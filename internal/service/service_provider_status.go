package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/model"
	"context"
)

// GetPendingServiceProviders 获取待审核服务提供商列表
func (s *ServiceProviderService) GetPendingServiceProviders(ctx context.Context, req *system.GetPendingSPReq) (res *system.GetPendingSPRes, err error) {
	serviceProviders, pageRes, err := s.serviceProviderDomain.GetPendingList(ctx, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.GetPendingSPRes{
		PageRes: pageRes,
	}
	for _, v := range serviceProviders {
		res.List = append(res.List, s.convertServiceProvider(v))
	}
	return res, nil
}

// ApproveServiceProvider 审核通过服务提供商
func (s *ServiceProviderService) ApproveServiceProvider(ctx context.Context, req *system.ApproveServiceProviderReq) (res *system.ApproveServiceProviderRes, err error) {
	err = s.serviceProviderDomain.HandleEvent(ctx, req.ID, model.ServiceProviderEventApprove, nil)
	if err != nil {
		return nil, err
	}
	return &system.ApproveServiceProviderRes{}, nil
}

// RejectServiceProvider 审核拒绝服务提供商
func (s *ServiceProviderService) RejectServiceProvider(ctx context.Context, req *system.RejectServiceProviderReq) (res *system.RejectServiceProviderRes, err error) {
	err = s.serviceProviderDomain.HandleEvent(ctx, req.ID, model.ServiceProviderEventReject, nil)
	if err != nil {
		return nil, err
	}
	return &system.RejectServiceProviderRes{}, nil
}

// DisableServiceProvider 禁用服务提供商
func (s *ServiceProviderService) DisableServiceProvider(ctx context.Context, req *system.DisableServiceProviderReq) (res *system.DisableServiceProviderRes, err error) {
	err = s.serviceProviderDomain.HandleEvent(ctx, req.ID, model.ServiceProviderEventDisable, nil)
	if err != nil {
		return nil, err
	}
	return &system.DisableServiceProviderRes{}, nil
}

// EnableServiceProvider 启用服务提供商
func (s *ServiceProviderService) EnableServiceProvider(ctx context.Context, req *system.EnableServiceProviderReq) (res *system.EnableServiceProviderRes, err error) {
	err = s.serviceProviderDomain.HandleEvent(ctx, req.ID, model.ServiceProviderEventEnable, nil)
	if err != nil {
		return nil, err
	}
	return &system.EnableServiceProviderRes{}, nil
}

func (s *ServiceProviderService) RecommitServiceProvider(ctx context.Context, req *system.RecommitServiceProviderReq) (res *system.RecommitServiceProviderRes, err error) {
	companyInfo := &model.Company{
		ID:          req.CompanyInfo.ID,
		Name:        req.CompanyInfo.Name,
		Country:     req.CompanyInfo.Country,
		City:        req.CompanyInfo.City,
		Address:     req.CompanyInfo.Address,
		Email:       req.CompanyInfo.Email,
		Description: req.CompanyInfo.Description,
	}
	for _, v := range req.CompanyInfo.Files {
		companyInfo.Files = append(companyInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}
	serviceProviderInfo := &model.ServiceProvider{
		ID:                 req.ID,
		Name:               req.Name,
		Website:            req.Website,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonEmail: req.ContactPersonEmail,
		Description:        req.Description,
	}
	for _, v := range req.Files {
		serviceProviderInfo.Files = append(serviceProviderInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}

	data := map[string]interface{}{
		"company_info":          companyInfo,
		"service_provider_info": serviceProviderInfo,
	}
	err = s.serviceProviderDomain.HandleEvent(ctx, req.ID, model.ServiceProviderEventReCommit, data)
	if err != nil {
		return nil, err
	}
	return &system.RecommitServiceProviderRes{}, nil
}
