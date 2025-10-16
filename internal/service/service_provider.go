package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type ServiceProviderService struct {
	serviceProviderDomain interfaces.IServiceProvider
}

func NewServiceProviderService(serviceProviderDomain interfaces.IServiceProvider) *ServiceProviderService {
	return &ServiceProviderService{
		serviceProviderDomain: serviceProviderDomain,
	}
}

// CreateServiceProvider 创建服务提供商
func (s *ServiceProviderService) CreateServiceProvider(ctx context.Context, req *system.CreateServiceProviderReq) (res *system.CreateServiceProviderRes, err error) {
	var (
		serviceProviderID   string
		companyInfo         *model.Company
		serviceProviderInfo *model.ServiceProvider
	)

	companyInfo = &model.Company{
		Name:        req.CompanyInfo.Name,
		Country:     req.CompanyInfo.Country,
		City:        req.CompanyInfo.City,
		Address:     req.CompanyInfo.Address,
		Email:       req.CompanyInfo.Email,
		Description: req.CompanyInfo.Description,

		SocialCreditCode:      req.CompanyInfo.SocialCreditCode,
		LegalPersonName:       req.CompanyInfo.LegalPersonName,
		LegalPersonCardNumber: req.CompanyInfo.LegalPersonCardNumber,
	}
	for _, v := range req.CompanyInfo.Files {
		companyInfo.Files = append(companyInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}
	// 创建服务提供商
	serviceProviderInfo = &model.ServiceProvider{
		Name:               req.Name,
		Website:            req.Website,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonEmail: req.ContactPersonEmail,
		Description:        req.Description,

		CompanyInfo: companyInfo,
	}
	for _, v := range req.Files {
		serviceProviderInfo.Files = append(serviceProviderInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		serviceProviderID, err = s.serviceProviderDomain.Create(ctx, tx, serviceProviderInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.CreateServiceProviderRes{ID: serviceProviderID}, nil
}

// GetServiceProvider 获取服务提供商详情
func (s *ServiceProviderService) GetServiceProvider(ctx context.Context, req *system.GetServiceProviderReq) (res *system.GetServiceProviderRes, err error) {
	var (
		serviceProviderInfo *model.ServiceProvider
	)

	serviceProviderInfo, err = s.serviceProviderDomain.GetServiceProvider(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res = &system.GetServiceProviderRes{
		ServiceProviderInfo: s.convertServiceProvider(serviceProviderInfo),
	}
	return res, nil
}

/*
// ListServiceProviders 列表服务提供商
func (s *ServiceProviderService) ListServiceProviders(ctx context.Context, req *model.ListServiceProvidersReq) (res *model.ListServiceProvidersRes, err error) {
	serviceProviders, pageRes, err := s.serviceProviderDomain.ListServiceProviders(ctx, req.Name, req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &model.ListServiceProvidersRes{
		ServiceProviders: serviceProviders,
		PageRes:          pageRes,
	}

	return res, nil
}

*/

func (s *ServiceProviderService) convertServiceProvider(in *model.ServiceProvider) (out *system.ServiceProviderInfo) {
	out = &system.ServiceProviderInfo{
		ID:                 in.ID,
		CompanyID:          in.CompanyID,
		Name:               in.Name,
		Status:             model.GetServiceProviderStatusText(in.Status),
		Description:        in.Description,
		ContactPersonName:  in.ContactPersonName,
		ContactPersonPhone: in.ContactPersonPhone,
		ContactPersonEmail: in.ContactPersonEmail,
		Website:            in.Website,

		CreateTime:          model.FormatTime(in.CreateTime),
		SubmitForReviewTime: model.FormatTime(in.SubmitForReviewTime),
		ApproveTime:         model.FormatTime(in.ApproveTime),
		UpdateTime:          model.FormatTime(in.UpdateTime),

		Files:       s.convertFiles(in.Files),
		CompanyInfo: s.convertCompanyInfo(in.CompanyInfo),
	}
	return out
}

func (s *ServiceProviderService) convertCompanyInfo(in *model.Company) *system.CompanyInfo {
	return &system.CompanyInfo{
		ID:          in.ID,
		Name:        in.Name,
		Country:     in.Country,
		City:        in.City,
		Address:     in.Address,
		Email:       in.Email,
		Description: in.Description,

		BusinessLicense:       in.BusinessLicense,
		SocialCreditCode:      in.SocialCreditCode,
		LegalPersonName:       in.LegalPersonName,
		LegalPersonCardNumber: in.LegalPersonCardNumber,
		LegalPersonPhoto:      in.LegalPersonPhoto,

		Files: s.convertFiles(in.Files),
	}
}

func (s *ServiceProviderService) convertFiles(in []*model.File) (out []*system.FileInfo) {
	for _, v := range in {
		out = append(out, &system.FileInfo{
			FileID:     v.FileID,
			FileType:   model.GetFileTypeText(v.Type),
			FileName:   v.FileName,
			FileLink:   v.FileLink,
			CreateTime: model.FormatTime(v.CreateTime),
			UpdateTime: model.FormatTime(v.UpdateTime),
		})
	}
	return out
}
