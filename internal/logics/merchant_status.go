package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func GetMerchantEventText(event interfaces.MerchantEvent) string {
	switch event {
	case interfaces.MerchantEventApprove:
		return "审核通过"
	case interfaces.MerchantEventReject:
		return "审核拒绝"
	case interfaces.MerchantEventDisable:
		return "禁用展商"
	case interfaces.MerchantEventEnable:
		return "启用展商"
	default:
		return "未知事件"
	}
}

type MerchantAction func(ctx context.Context, merchant *model.Merchant, data interface{}) error

type MerchantTransition struct {
	State  model.MerchantStatus
	Action MerchantAction
}

var merchantTransitionMap = map[model.MerchantStatus]map[interfaces.MerchantEvent]MerchantTransition{
	model.MerchantStatusPending: {
		interfaces.MerchantEventApprove: {
			State:  model.MerchantStatusApproved,
			Action: handleMerchantApprove,
		},
		interfaces.MerchantEventReject: {
			State:  model.MerchantStatusDisabled,
			Action: handleMerchantReject,
		},
	},
	model.MerchantStatusApproved: {
		interfaces.MerchantEventDisable: {
			State:  model.MerchantStatusDisabled,
			Action: handleMerchantDisable,
		},
	},
	model.MerchantStatusDisabled: {
		interfaces.MerchantEventEnable: {
			State:  model.MerchantStatusApproved,
			Action: handleMerchantEnable,
		},
	},
}

func handleMerchantApprove(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:     int(model.MerchantStatusApproved),
		dao.Merchant.Columns().Version:    merchant.Version + 1,
		dao.Merchant.Columns().UpdateTime: time.Now().Unix(),
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

func handleMerchantReject(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
	result, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Status:     int(model.MerchantStatusDisabled),
		dao.Merchant.Columns().Version:    merchant.Version + 1,
		dao.Merchant.Columns().UpdateTime: time.Now().Unix(),
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

func handleMerchantDisable(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
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

func handleMerchantEnable(ctx context.Context, merchant *model.Merchant, data interface{}) (err error) {
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
