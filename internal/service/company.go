package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type companyService struct {
	companyDomain interfaces.ICompany
	fileDomain    interfaces.IFile
}

func NewCompanyService(companyDomain interfaces.ICompany, fileDomain interfaces.IFile) *companyService {
	return &companyService{
		companyDomain: companyDomain,
		fileDomain:    fileDomain,
	}
}

func (s *companyService) CreateCompany(ctx context.Context, req *system.CreateCompanyReq) (res *system.CreateCompanyRes, err error) {
	file, err := s.fileDomain.GetFile(ctx, req.LegalPersonPhotoFileID)
	if err != nil {
		return nil, err
	}

	in := model.Company{
		Name:        req.Name,
		Country:     req.Country,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		Description: req.Description,

		BusinessLicense:       req.BusinessLicense,
		SocialCreditCode:      req.SocialCreditCode,
		LegalPersonName:       req.LegalPersonName,
		LegalPersonCardNumber: req.LegalPersonCardNumber,
		LegalPersonPhotoUrl:   file.FileLink,
		LegalPersonPhone:      req.LegalPersonPhone,
	}

	var companyID string
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		companyID, err = s.companyDomain.CreateCompany(ctx, tx, &in)
		if err != nil {
			return err
		}

		err = s.fileDomain.UpdateFileCompanyInfo(ctx, tx, req.LegalPersonPhotoFileID, companyID, model.FileTypeLegalPersonPhoto)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.CreateCompanyRes{
		ID: companyID,
	}, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, req *system.DeleteCompanyReq) (res *system.DeleteCompanyRes, err error) {
	if err := s.companyDomain.DeleteCompany(ctx, req.ID); err != nil {
		return nil, err
	}
	return &system.DeleteCompanyRes{}, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, req *system.UpdateCompanyReq) (*system.UpdateCompanyRes, error) {
	in := &model.Company{
		ID:                    req.ID,
		Name:                  req.Name,
		Country:               req.Country,
		Phone:                 req.Phone,
		Email:                 req.Email,
		Address:               req.Address,
		Description:           req.Description,
		BusinessLicense:       req.BusinessLicense,
		SocialCreditCode:      req.SocialCreditCode,
		LegalPersonName:       req.LegalPersonName,
		LegalPersonCardNumber: req.LegalPersonCardNumber,
		LegalPersonPhotoUrl:   req.LegalPersonPhotoUrl,
		LegalPersonPhone:      req.LegalPersonPhone,
	}
	if err := s.companyDomain.UpdateCompany(ctx, in); err != nil {
		return nil, err
	}
	return &system.UpdateCompanyRes{}, nil
}

func (s *companyService) GetCompany(ctx context.Context, req *system.GetCompanyReq) (*system.GetCompanyRes, error) {
	c, err := s.companyDomain.GetCompany(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &system.GetCompanyRes{Data: convertCompany(c)}, nil
}

func (s *companyService) ListCompanies(ctx context.Context, req *system.ListCompanyReq) (*system.ListCompanyRes, error) {
	list, pageRes, err := s.companyDomain.ListCompanies(ctx, req.Name, &req.PageReq)
	if err != nil {
		return nil, err
	}

	companies := make([]*system.Company, 0, len(list))
	for _, c := range list {
		companies = append(companies, convertCompany(c))
	}
	return &system.ListCompanyRes{List: companies, PageRes: pageRes}, nil
}

func (s *companyService) ApproveCompany(ctx context.Context, req *system.ApproveCompanyReq) (*system.ApproveCompanyRes, error) {
	if err := s.companyDomain.ApproveCompany(ctx, req.ID); err != nil {
		return nil, err
	}
	return &system.ApproveCompanyRes{}, nil
}

func (s *companyService) RejectCompany(ctx context.Context, req *system.RejectCompanyReq) (*system.RejectCompanyRes, error) {
	if err := s.companyDomain.RejectCompany(ctx, req.ID); err != nil {
		return nil, err
	}
	return &system.RejectCompanyRes{}, nil
}

func (s *companyService) BanCompany(ctx context.Context, req *system.BanCompanyReq) (*system.BanCompanyRes, error) {
	if err := s.companyDomain.BanCompany(ctx, req.ID); err != nil {
		return nil, err
	}
	return &system.BanCompanyRes{}, nil
}

func (s *companyService) UnbanCompany(ctx context.Context, req *system.UnbanCompanyReq) (*system.UnbanCompanyRes, error) {
	if err := s.companyDomain.UnbanCompany(ctx, req.ID); err != nil {
		return nil, err
	}
	return &system.UnbanCompanyRes{}, nil
}

func (s *companyService) ListApplications(ctx context.Context, req *system.ListApplicationsReq) (*system.ListApplicationsRes, error) {
	list, pageRes, err := s.companyDomain.ListApplications(ctx, &req.PageReq)
	if err != nil {
		return nil, err
	}

	companies := make([]*system.Company, 0, len(list))
	for _, c := range list {
		companies = append(companies, convertCompany(c))
	}
	return &system.ListApplicationsRes{List: companies, PageRes: pageRes}, nil
}

func convertCompany(in *model.Company) *system.Company {
	return &system.Company{
		ID:      in.ID,
		Name:    in.Name,
		Country: in.Country,
		Status:  model.GetCompanyStatusText(in.Status),
		Phone:   in.Phone,
		Email:   in.Email,
		Address: in.Address,

		BusinessLicense:       in.BusinessLicense,
		SocialCreditCode:      in.SocialCreditCode,
		LegalPersonName:       in.LegalPersonName,
		LegalPersonCardNumber: in.LegalPersonCardNumber,
		LegalPersonPhotoUrl:   in.LegalPersonPhotoUrl,
		LegalPersonPhone:      in.LegalPersonPhone,
		ApplyTime:             in.ApplyTime.Format(time.DateTime),
		ApproveTime:           in.ApproveTime.Format(time.DateTime),
		CreateTime:            in.CreateTime.Format(time.DateTime),
		UpdateTime:            in.UpdateTime.Format(time.DateTime),
	}
}
