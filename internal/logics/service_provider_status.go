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

type ServiceProviderAction func(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) error

type ServiceProviderTransition struct {
	State  model.ServiceProviderStatus
	Action ServiceProviderAction
}

func (sp *serviceProvider) initTransitionMap() {
	sp.transitionMap = map[model.ServiceProviderStatus]map[model.ServiceProviderEvent]ServiceProviderTransition{
		model.ServiceProviderStatusPending: {
			model.ServiceProviderEventApprove: {
				State:  model.ServiceProviderStatusApproved,
				Action: sp.handleServiceProviderApprove,
			},
			model.ServiceProviderEventReject: {
				State:  model.ServiceProviderStatusDisabled,
				Action: sp.handleServiceProviderReject,
			},
		},
		model.ServiceProviderStatusRejected: {
			model.ServiceProviderEventReCommit: {
				State:  model.ServiceProviderStatusPending,
				Action: sp.handleServiceProviderReCommit,
			},
		},
		model.ServiceProviderStatusApproved: {
			model.ServiceProviderEventDisable: {
				State:  model.ServiceProviderStatusDisabled,
				Action: sp.handleServiceProviderDisable,
			},
			model.ServiceProviderEventUnregister: {
				State:  model.ServiceProviderStatusUnregistered,
				Action: sp.handleServiceProviderUnregister,
			},
		},
		model.ServiceProviderStatusDisabled: {
			model.ServiceProviderEventEnable: {
				State:  model.ServiceProviderStatusApproved,
				Action: sp.handleServiceProviderEnable,
			},
		},
	}
}

func (sp *serviceProvider) GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.ServiceProvider, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.ServiceProvider.Ctx(ctx).Where(dao.ServiceProvider.Columns().Status, int(model.ServiceProviderStatusPending))
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("get pending service providers failed, query count err: %s", err.Error())
	}

	var tsp []*entity.TServiceProvider
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.ServiceProvider.Columns().CreateTime)
	err = query.Scan(&tsp)
	if err != nil {
		return nil, nil, gerror.Newf("get pending service providers failed, query scan err: %s", err.Error())
	}

	for _, r := range tsp {
		tmp := model.ConvertServiceProvider(r)
		tmp.CompanyInfo, err = sp.companyDomain.Get(ctx, tmp.CompanyID)
		if err != nil {
			return nil, nil, gerror.Newf("get pending service providers failed, get company info failed, err: %s", err.Error())
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
func (sp *serviceProvider) HandleEvent(ctx context.Context, serviceProviderID string, event model.ServiceProviderEvent, data interface{}) (err error) {
	// 获取服务提供商信息
	serviceProvider, err := sp.GetServiceProvider(ctx, serviceProviderID)
	if err != nil {
		return err
	}

	// 验证当前状态是否支持该事件
	transition, ok := sp.transitionMap[serviceProvider.Status][event]
	if !ok {
		return gerror.Newf("service provider current status: %s, not supported event: %s", model.GetServiceProviderStatusText(serviceProvider.Status), model.GetServiceProviderEventText(event))
	}

	// 执行状态操作
	err = transition.Action(ctx, serviceProvider, data)
	if err != nil {
		return gerror.Newf("handle event failed, service provider id: %s, event: %s, err: %s", serviceProviderID, model.GetServiceProviderEventText(event), err.Error())
	}

	return nil
}

func (sp *serviceProvider) handleServiceProviderApprove(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	result, err := dao.ServiceProvider.Ctx(ctx).Data(g.Map{
		dao.ServiceProvider.Columns().Status:      int(model.ServiceProviderStatusApproved),
		dao.ServiceProvider.Columns().Version:     serviceProvider.Version + 1,
		dao.ServiceProvider.Columns().ApproveTime: time.Now().Unix(),
		dao.ServiceProvider.Columns().UpdateTime:  time.Now().Unix(),
	}).
		Where(dao.ServiceProvider.Columns().ID, serviceProvider.ID).
		Where(dao.ServiceProvider.Columns().Version, serviceProvider.Version).
		Update()
	if err != nil {
		return gerror.Newf("approve service provider failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("approve service provider failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("approve service provider failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (sp *serviceProvider) handleServiceProviderReject(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	result, err := dao.ServiceProvider.Ctx(ctx).Data(g.Map{
		dao.ServiceProvider.Columns().Status:              int(model.ServiceProviderStatusDisabled),
		dao.ServiceProvider.Columns().Version:             serviceProvider.Version + 1,
		dao.ServiceProvider.Columns().SubmitForReviewTime: 0,
		dao.ServiceProvider.Columns().UpdateTime:          time.Now().Unix(),
	}).
		Where(dao.ServiceProvider.Columns().ID, serviceProvider.ID).
		Where(dao.ServiceProvider.Columns().Version, serviceProvider.Version).
		Update()
	if err != nil {
		return gerror.Newf("reject service provider failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("reject service provider failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("reject service provider failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (sp *serviceProvider) handleServiceProviderReCommit(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	newCompanyInfo := data.(map[string]interface{})["company_info"].(*model.Company)
	newSPInfo := data.(map[string]interface{})["service_provider_info"].(*model.ServiceProvider)

	var dataUpdate map[string]any = make(map[string]any)
	if serviceProvider.Name != newSPInfo.Name {
		dataUpdate[dao.ServiceProvider.Columns().Name] = newSPInfo.Name
	}
	if serviceProvider.Description != newSPInfo.Description {
		dataUpdate[dao.ServiceProvider.Columns().Description] = newSPInfo.Description
	}
	if serviceProvider.ContactPersonName != newSPInfo.ContactPersonName {
		dataUpdate[dao.ServiceProvider.Columns().ContactPersonName] = newSPInfo.ContactPersonName
	}
	if serviceProvider.ContactPersonPhone != newSPInfo.ContactPersonPhone {
		dataUpdate[dao.ServiceProvider.Columns().ContactPersonPhone] = newSPInfo.ContactPersonPhone
	}
	if serviceProvider.ContactPersonEmail != newSPInfo.ContactPersonEmail {
		dataUpdate[dao.ServiceProvider.Columns().ContactPersonEmail] = newSPInfo.ContactPersonEmail
	}
	if serviceProvider.Website != newSPInfo.Website {
		dataUpdate[dao.ServiceProvider.Columns().Website] = newSPInfo.Website
	}
	dataUpdate[dao.ServiceProvider.Columns().Status] = int(model.ServiceProviderStatusPending)
	dataUpdate[dao.ServiceProvider.Columns().SubmitForReviewTime] = time.Now().Unix()
	dataUpdate[dao.ServiceProvider.Columns().Version] = serviceProvider.Version + 1
	dataUpdate[dao.ServiceProvider.Columns().UpdateTime] = time.Now().Unix()

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = sp.companyDomain.Recommit(ctx, tx, model.CompanyTypeServiceProvider, newCompanyInfo)
		if err != nil {
			return gerror.Newf("recommit service provider failed, recommit company failed, err: %s", err.Error())
		}

		_, err = dao.ServiceProvider.Ctx(ctx).TX(tx).Data(dataUpdate).Where(dao.ServiceProvider.Columns().ID, newSPInfo.ID).Update()
		if err != nil {
			return gerror.Newf("recommit service provider failed, update service provider failed, err: %s", err.Error())
		}

		// 检查文件是否上传成功
		err = sp.checkFileUploadSuccess(ctx, newSPInfo.Files)
		if err != nil {
			return gerror.Newf("recommit service provider failed, check file upload success failed, err: %s", err.Error())
		}

		// 检查服务提供商文件是否完整
		err = sp.checkFileComplete(ctx, newSPInfo.Files)
		if err != nil {
			return gerror.Newf("recommit service provider failed, check file complete failed, err: %s", err.Error())
		}

		// 清除旧文件关联
		err = sp.clearFileRelation(ctx, tx, newSPInfo.ID)
		if err != nil {
			return gerror.Newf("recommit service provider failed, clear file relation failed, err: %s", err.Error())
		}

		// 创建新文件关联
		err = sp.createFileRelation(ctx, tx, newSPInfo.ID, newSPInfo.Files)
		if err != nil {
			return gerror.Newf("recommit service provider failed, create file relation failed, err: %s", err.Error())
		}

		return nil
	})
}

func (sp *serviceProvider) handleServiceProviderDisable(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	result, err := dao.ServiceProvider.Ctx(ctx).Data(g.Map{
		dao.ServiceProvider.Columns().Status:     int(model.ServiceProviderStatusDisabled),
		dao.ServiceProvider.Columns().Version:    serviceProvider.Version + 1,
		dao.ServiceProvider.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.ServiceProvider.Columns().ID, serviceProvider.ID).
		Where(dao.ServiceProvider.Columns().Version, serviceProvider.Version).
		Update()
	if err != nil {
		return gerror.Newf("disable service provider failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("disable service provider failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("disable service provider failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (sp *serviceProvider) handleServiceProviderEnable(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	result, err := dao.ServiceProvider.Ctx(ctx).Data(g.Map{
		dao.ServiceProvider.Columns().Status:     int(model.ServiceProviderStatusApproved),
		dao.ServiceProvider.Columns().Version:    serviceProvider.Version + 1,
		dao.ServiceProvider.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.ServiceProvider.Columns().ID, serviceProvider.ID).
		Where(dao.ServiceProvider.Columns().Version, serviceProvider.Version).
		Update()
	if err != nil {
		return gerror.Newf("enable service provider failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("enable service provider failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("enable service provider failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (sp *serviceProvider) handleServiceProviderUnregister(ctx context.Context, serviceProvider *model.ServiceProvider, data interface{}) (err error) {
	result, err := dao.ServiceProvider.Ctx(ctx).Data(g.Map{
		dao.ServiceProvider.Columns().Status:     int(model.ServiceProviderStatusUnregistered),
		dao.ServiceProvider.Columns().Version:    serviceProvider.Version + 1,
		dao.ServiceProvider.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.ServiceProvider.Columns().ID, serviceProvider.ID).
		Where(dao.ServiceProvider.Columns().Version, serviceProvider.Version).
		Update()
	if err != nil {
		return gerror.Newf("unregister service provider failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("unregister service provider failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("unregister service provider failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}
