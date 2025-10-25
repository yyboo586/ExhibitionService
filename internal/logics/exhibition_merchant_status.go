package logics

import (
	"ExhibitionService/internal/dao"
	"context"
	"time"

	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ExMerchantTransition 状态转换函数类型
type ExMerchantAction func(ctx context.Context, exhibitionID string, merchantID string, data interface{}) error

type ExMerchantTransition struct {
	State  model.ExMerchantStatus
	Action ExMerchantAction
}

// initExMerchantTransitionMap 初始化状态转换映射
func (e *exhibition) initExMerchantTransitionMap() {
	e.exMerchantTransitionMap = map[model.ExMerchantStatus]map[model.ExMerchantEvent]ExMerchantTransition{
		model.ExMerchantStatusPending: {
			model.ExMerchantEventApprove: {
				State:  model.ExMerchantStatusApproved,
				Action: e.approveExMerchant,
			},
			model.ExMerchantEventReject: {
				State:  model.ExMerchantStatusRejected,
				Action: e.rejectExMerchant,
			},
			model.ExMerchantEventWithdraw: {
				State:  model.ExMerchantStatusWithdrawn,
				Action: e.withdrawExMerchant,
			},
		},
		model.ExMerchantStatusApproved: {
			model.ExMerchantEventWithdraw: {
				State:  model.ExMerchantStatusWithdrawn,
				Action: e.withdrawExMerchant,
			},
		},
		model.ExMerchantStatusRejected: {
			model.ExMerchantEventReApply: {
				State:  model.ExMerchantStatusPending,
				Action: e.reapplyForExhibition,
			},
		},
	}
}

// HandleEvent 处理展会与商户关联状态事件
func (e *exhibition) HandleExMerchantEvent(ctx context.Context, exhibitionID string, merchantID string, event model.ExMerchantEvent, data interface{}) (err error) {
	// 获取当前申请状态
	application, err := e.GetExMerchantApplication(ctx, exhibitionID, merchantID)
	if err != nil {
		return err
	}

	// 检查状态转换是否合法
	action, exists := e.exMerchantTransitionMap[application.Status][event]
	if !exists {
		return gerror.Newf("invalid state transition from %s to %s",
			model.GetExMerchantStatusText(application.Status),
			model.GetExMerchantEventText(event))
	}

	// 执行状态转换
	err = action.Action(ctx, exhibitionID, merchantID, data)
	if err != nil {
		return gerror.Newf("handle ex merchant event failed, exhibition id: %s, merchant id: %s, event: %s, err: %s", exhibitionID, merchantID, model.GetExMerchantEventText(event), err.Error())
	}

	g.Log().Infof(ctx, "handle ex merchant event success, exhibition id: %s, merchant id: %s, event: %s", exhibitionID, merchantID, model.GetExMerchantEventText(event))
	return nil
}

// approveExMerchant 审核通过
func (e *exhibition) approveExMerchant(ctx context.Context, exhibitionID string, merchantID string, data interface{}) error {
	_, err := dao.ExhibitionMerchant.Ctx(ctx).
		Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID).
		Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID).
		Data(map[string]any{
			dao.ExhibitionMerchant.Columns().Status:      int(model.ExMerchantStatusApproved),
			dao.ExhibitionMerchant.Columns().ApproveTime: time.Now().Unix(),
			dao.ExhibitionMerchant.Columns().UpdateTime:  time.Now().Unix(),
		}).Update()
	if err != nil {
		return gerror.Newf("approve ex merchant failed, err: %s", err.Error())
	}
	return nil
}

// rejectExMerchant 审核驳回
func (e *exhibition) rejectExMerchant(ctx context.Context, exhibitionID string, merchantID string, data interface{}) error {
	_, err := dao.ExhibitionMerchant.Ctx(ctx).
		Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID).
		Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID).
		Data(map[string]any{
			dao.ExhibitionMerchant.Columns().Status:     int(model.ExMerchantStatusRejected),
			dao.ExhibitionMerchant.Columns().UpdateTime: time.Now().Unix(),
		}).Update()
	if err != nil {
		return gerror.Newf("reject application failed, err: %s", err.Error())
	}
	return nil
}

// withdrawExMerchant 退出
func (e *exhibition) withdrawExMerchant(ctx context.Context, exhibitionID string, merchantID string, data interface{}) error {
	_, err := dao.ExhibitionMerchant.Ctx(ctx).
		Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID).
		Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID).
		Data(map[string]any{
			dao.ExhibitionMerchant.Columns().Status:     int(model.ExMerchantStatusWithdrawn),
			dao.ExhibitionMerchant.Columns().UpdateTime: time.Now().Unix(),
		}).Update()
	if err != nil {
		return gerror.Newf("withdraw application failed, err: %s", err.Error())
	}
	return nil
}

// reapplyForExhibition 重新申请
func (e *exhibition) reapplyForExhibition(ctx context.Context, exhibitionID string, merchantID string, data interface{}) error {
	_, err := dao.ExhibitionMerchant.Ctx(ctx).
		Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID).
		Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID).
		Data(map[string]any{
			dao.ExhibitionMerchant.Columns().Status:              int(model.ExMerchantStatusPending),
			dao.ExhibitionMerchant.Columns().SubmitForReviewTime: time.Now().Unix(),
			dao.ExhibitionMerchant.Columns().UpdateTime:          time.Now().Unix(),
		}).Update()
	if err != nil {
		return gerror.Newf("reapply for exhibition failed, err: %s", err.Error())
	}
	return nil
}
