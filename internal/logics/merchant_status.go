package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type MerchantAction func(ctx context.Context, merchant *model.Merchant, data interface{}) error

type MerchantTransition struct {
	State  model.MerchantStatus
	Action MerchantAction
}

func (m *merchant) initTransitionMap() {
	m.transitionMap = map[model.MerchantStatus]map[model.MerchantEvent]MerchantTransition{
		model.MerchantStatusPending: {
			model.MerchantEventApprove: {
				State:  model.MerchantStatusApproved,
				Action: m.handleMerchantApprove,
			},
			model.MerchantEventReject: {
				State:  model.MerchantStatusRejected,
				Action: m.handleMerchantReject,
			},
		},
		model.MerchantStatusRejected: {
			model.MerchantEventReCommit: {
				State:  model.MerchantStatusPending,
				Action: m.handleMerchantReCommit,
			},
		},
		model.MerchantStatusApproved: {
			model.MerchantEventDisable: {
				State:  model.MerchantStatusDisabled,
				Action: m.handleMerchantDisable,
			},
			model.MerchantEventUnregister: {
				State:  model.MerchantStatusUnregistered,
				Action: m.handleMerchantUnregister,
			},
		},
		model.MerchantStatusDisabled: {
			model.MerchantEventEnable: {
				State:  model.MerchantStatusApproved,
				Action: m.handleMerchantEnable,
			},
		},
	}
}
func (m *merchant) GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().Status, int(model.MerchantStatusPending))
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("get pending merchants failed, query count err: %s", err.Error())
	}

	var tm []*entity.TMerchant
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Merchant.Columns().CreateTime)
	err = query.Scan(&tm)
	if err != nil {
		return nil, nil, gerror.Newf("get pending merchants failed, query scan err: %s", err.Error())
	}

	for _, r := range tm {
		tmp := model.ConvertMerchant(r)
		tmp.CompanyInfo, err = m.companyDomain.Get(ctx, tmp.CompanyID)
		if err != nil {
			return nil, nil, gerror.Newf("get pending merchants failed, get company info failed, err: %s", err.Error())
		}

		out = append(out, tmp)
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

// 状态流转
func (m *merchant) HandleEvent(ctx context.Context, merchantID string, event model.MerchantEvent, data interface{}) (err error) {
	// 获取展商信息
	merchant, err := m.Get(ctx, merchantID)
	if err != nil {
		return err
	}

	// 验证当前状态是否支持该事件
	transition, ok := m.transitionMap[merchant.Status][event]
	if !ok {
		return gerror.Newf("merchant current status: %s, not supported event: %s", model.GetMerchantStatusText(merchant.Status), model.GetMerchantEventText(event))
	}

	// 执行状态操作
	err = transition.Action(ctx, merchant, data)
	if err != nil {
		return gerror.Newf("handle event failed, merchant id: %s, event: %s, err: %s", merchantID, model.GetMerchantEventText(event), err.Error())
	}

	return nil
}

func (m *merchant) handleMerchantApprove(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:      int(model.MerchantStatusApproved),
		dao.Merchant.Columns().Version:     merchant.Version + 1,
		dao.Merchant.Columns().ApproveTime: time.Now().Unix(),
		dao.Merchant.Columns().UpdateTime:  time.Now().Unix(),
	}).
		Where(dao.Merchant.Columns().ID, merchant.ID).
		Where(dao.Merchant.Columns().Version, merchant.Version).
		Update()
	if err != nil {
		return gerror.Newf("approve merchant failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("approve merchant failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("approve merchant failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (m *merchant) handleMerchantReject(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:              int(model.MerchantStatusRejected),
		dao.Merchant.Columns().Version:             merchant.Version + 1,
		dao.Merchant.Columns().SubmitForReviewTime: 0,
		dao.Merchant.Columns().UpdateTime:          time.Now().Unix(),
	}).
		Where(dao.Merchant.Columns().ID, merchant.ID).
		Where(dao.Merchant.Columns().Version, merchant.Version).
		Update()
	if err != nil {
		return gerror.Newf("reject merchant failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("reject merchant failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("reject merchant failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (m *merchant) handleMerchantDisable(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:     int(model.MerchantStatusDisabled),
		dao.Merchant.Columns().Version:    merchant.Version + 1,
		dao.Merchant.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Merchant.Columns().ID, merchant.ID).
		Where(dao.Merchant.Columns().Version, merchant.Version).
		Update()
	if err != nil {
		return gerror.Newf("disable merchant failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("disable merchant failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("disable merchant failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (m *merchant) handleMerchantEnable(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:     int(model.MerchantStatusApproved),
		dao.Merchant.Columns().Version:    merchant.Version + 1,
		dao.Merchant.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Merchant.Columns().ID, merchant.ID).
		Where(dao.Merchant.Columns().Version, merchant.Version).
		Update()
	if err != nil {
		return gerror.Newf("enable merchant failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("enable merchant failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("enable merchant failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (m *merchant) handleMerchantReCommit(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	newMerchantInfo := data.(*model.Merchant)

	var dataUpdate map[string]any = make(map[string]any)
	if merchant.Name != newMerchantInfo.Name {
		dataUpdate[dao.Merchant.Columns().Name] = newMerchantInfo.Name
	}
	if merchant.Description != newMerchantInfo.Description {
		dataUpdate[dao.Merchant.Columns().Description] = newMerchantInfo.Description
	}
	if merchant.ContactPersonName != newMerchantInfo.ContactPersonName {
		dataUpdate[dao.Merchant.Columns().ContactPersonName] = newMerchantInfo.ContactPersonName
	}
	if merchant.ContactPersonPhone != newMerchantInfo.ContactPersonPhone {
		dataUpdate[dao.Merchant.Columns().ContactPersonPhone] = newMerchantInfo.ContactPersonPhone
	}
	if merchant.ContactPersonEmail != newMerchantInfo.ContactPersonEmail {
		dataUpdate[dao.Merchant.Columns().ContactPersonEmail] = newMerchantInfo.ContactPersonEmail
	}
	if merchant.Website != newMerchantInfo.Website {
		dataUpdate[dao.Merchant.Columns().Website] = newMerchantInfo.Website
	}
	dataUpdate[dao.Merchant.Columns().Status] = int(model.MerchantStatusPending)
	dataUpdate[dao.Merchant.Columns().SubmitForReviewTime] = time.Now().Unix()
	dataUpdate[dao.Merchant.Columns().Version] = merchant.Version + 1
	dataUpdate[dao.Merchant.Columns().UpdateTime] = time.Now().Unix()

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = m.companyDomain.Recommit(ctx, tx, model.CompanyTypeMerchant, newMerchantInfo.CompanyInfo)
		if err != nil {
			return gerror.Newf("recommit merchant failed, recommit company failed, err: %s", err.Error())
		}

		_, err = dao.Merchant.Ctx(ctx).TX(tx).Data(dataUpdate).Where(dao.Merchant.Columns().ID, newMerchantInfo.ID).Update()
		if err != nil {
			return gerror.Newf("recommit merchant failed, update merchant failed, err: %s", err.Error())
		}

		// 检查文件是否上传成功
		err = m.checkFileUploadSuccess(ctx, newMerchantInfo.Files)
		if err != nil {
			return gerror.Newf("recommit merchant failed, check file upload success failed, err: %s", err.Error())
		}

		// 检查商户文件是否完整
		err = m.checkFileComplete(ctx, newMerchantInfo.Files)
		if err != nil {
			return gerror.Newf("recommit merchant failed, check file complete failed, err: %s", err.Error())
		}

		// 清除旧文件关联
		err = m.clearFileRelation(ctx, tx, newMerchantInfo.ID)
		if err != nil {
			return gerror.Newf("recommit merchant failed, clear file relation failed, err: %s", err.Error())
		}

		err = m.createFileRelation(ctx, tx, newMerchantInfo.ID, newMerchantInfo.Files)
		if err != nil {
			return gerror.Newf("recommit merchant failed, create file relation failed, err: %s", err.Error())
		}

		return nil
	})
}

func (m *merchant) handleMerchantUnregister(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:     int(model.MerchantStatusUnregistered),
		dao.Merchant.Columns().Version:    merchant.Version + 1,
		dao.Merchant.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Merchant.Columns().ID, merchant.ID).
		Where(dao.Merchant.Columns().Version, merchant.Version).
		Update()
	if err != nil {
		return gerror.Newf("unregister merchant failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("unregister merchant failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("unregister merchant failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}
