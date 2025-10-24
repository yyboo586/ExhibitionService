package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/model"
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func GetExhibitionEventText(event model.ExhibitionEvent) string {
	switch event {
	case model.ExhibitionEventSubmitForReview:
		return "提交审核"
	case model.ExhibitionEventApprove:
		return "审核通过"
	case model.ExhibitionEventReject:
		return "审核驳回"
	case model.ExhibitionEventStartEnrolling:
		return "开始报名"
	case model.ExhibitionEventStartRunning:
		return "开始进行"
	case model.ExhibitionEventEnd:
		return "结束展会"
	case model.ExhibitionEventCancel:
		return "取消展会"
	default:
		return "未知事件"
	}
}

type ExhibitionAction func(ctx context.Context, exhibition *model.Exhibition, data interface{}) error

type ExhibitionTransition struct {
	State  model.ExhibitionStatus
	Action ExhibitionAction
}

var exhibitionTransitionMap = map[model.ExhibitionStatus]map[model.ExhibitionEvent]ExhibitionTransition{
	model.ExhibitionStatusPreparing: {
		model.ExhibitionEventSubmitForReview: {
			State:  model.ExhibitionStatusPending,
			Action: handleExhibitionSubmitForReview,
		},
		model.ExhibitionEventCancel: {
			State:  model.ExhibitionStatusCancelled,
			Action: handleExhibitionCancel,
		},
	},
	model.ExhibitionStatusPending: {
		model.ExhibitionEventApprove: {
			State:  model.ExhibitionStatusApproved,
			Action: handleExhibitionApprove,
		},
		model.ExhibitionEventReject: {
			State:  model.ExhibitionStatusPreparing,
			Action: handleExhibitionReject,
		},
	},
	model.ExhibitionStatusApproved: {
		model.ExhibitionEventStartEnrolling: {
			State:  model.ExhibitionStatusEnrolling,
			Action: handleExhibitionStartEnrolling,
		},
		model.ExhibitionEventCancel: {
			State:  model.ExhibitionStatusCancelled,
			Action: handleExhibitionCancel,
		},
	},
	model.ExhibitionStatusEnrolling: {
		model.ExhibitionEventStartRunning: {
			State:  model.ExhibitionStatusRunning,
			Action: handleExhibitionStartRunning,
		},
		model.ExhibitionEventCancel: {
			State:  model.ExhibitionStatusCancelled,
			Action: handleExhibitionCancel,
		},
	},
	model.ExhibitionStatusRunning: {
		model.ExhibitionEventEnd: {
			State:  model.ExhibitionStatusEnded,
			Action: handleExhibitionEnd,
		},
		model.ExhibitionEventCancel: {
			State:  model.ExhibitionStatusCancelled,
			Action: handleExhibitionCancel,
		},
	},
}

func (e *exhibition) HandleEvent(ctx context.Context, exhibitionID string, event model.ExhibitionEvent, data interface{}) (err error) {
	// 获取展会信息
	exhibition, err := e.GetExhibition(ctx, exhibitionID)
	if err != nil {
		return err
	}

	// 验证当前状态是否支持该事件
	transition, ok := exhibitionTransitionMap[exhibition.Status][event]
	if !ok {
		return gerror.Newf("exhibition current status: %s, not supported event: %s", model.GetExhibitionStatusText(exhibition.Status), GetExhibitionEventText(event))
	}

	// 执行状态操作
	err = transition.Action(ctx, exhibition, data)
	if err != nil {
		return gerror.Newf("handle event failed, exhibition id: %s, event: %s, err: %s", exhibitionID, GetExhibitionEventText(event), err.Error())
	}

	return nil
}

func handleExhibitionSubmitForReview(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusPending),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("submit for review failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("submit for review failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("submit for review failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionApprove(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusApproved),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("approve exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("approve exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("approve exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionReject(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusPreparing),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("reject exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("reject exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("reject exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionStartEnrolling(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusEnrolling),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("start enrolling failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("start enrolling failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("start enrolling failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionStartRunning(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusRunning),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("start running failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("start running failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("start running failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionEnd(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusEnded),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("end exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("end exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("end exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func handleExhibitionCancel(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusCancelled),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("cancel exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("cancel exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("cancel exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}
