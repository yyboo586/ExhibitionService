package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type MerchantService struct {
	merchantDomain interfaces.IMerchant
	fileDomain     interfaces.IFile
	companyDomain  interfaces.ICompany
}

func NewMerchantService(merchantDomain interfaces.IMerchant, fileDomain interfaces.IFile, companyDomain interfaces.ICompany) *MerchantService {
	return &MerchantService{
		merchantDomain: merchantDomain,
		fileDomain:     fileDomain,
		companyDomain:  companyDomain,
	}
}

// CreateMerchant 创建商户
func (s *MerchantService) CreateMerchant(ctx context.Context, req *system.CreateMerchantReq) (res *system.CreateMerchantRes, err error) {
	var (
		merchantID   string
		companyInfo  *model.Company
		merchantInfo *model.Merchant
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
	// 创建商户
	merchantInfo = &model.Merchant{
		Name:               req.Name,
		Website:            req.Website,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonEmail: req.ContactPersonEmail,
		Description:        req.Description,

		CompanyInfo: companyInfo,
	}
	for _, v := range req.Files {
		merchantInfo.Files = append(merchantInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建商户
		merchantID, err = s.merchantDomain.Create(ctx, tx, merchantInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.CreateMerchantRes{ID: merchantID}, nil
}

// GetMerchant 获取商户详情
func (s *MerchantService) GetMerchant(ctx context.Context, req *system.GetMerchantReq) (res *system.GetMerchantRes, err error) {
	merchant, err := s.merchantDomain.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res = &system.GetMerchantRes{
		MerchantInfo: s.convertMerchant(merchant),
	}

	return res, nil
}

func (s *MerchantService) ListMerchants(ctx context.Context, req *system.ListMerchantsReq) (res *system.ListMerchantsRes, err error) {
	merchants, pageRes, err := s.merchantDomain.List(ctx, req.Name, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.ListMerchantsRes{
		PageRes: pageRes,
	}
	for _, v := range merchants {
		res.List = append(res.List, s.convertMerchant(v))
	}

	return res, nil
}

func (s *MerchantService) convertMerchant(in *model.Merchant) (out *system.MerchantInfo) {
	out = &system.MerchantInfo{
		ID:                  in.ID,
		Name:                in.Name,
		Status:              model.GetMerchantStatusText(in.Status),
		Website:             in.Website,
		ContactPersonName:   in.ContactPersonName,
		ContactPersonPhone:  in.ContactPersonPhone,
		ContactPersonEmail:  in.ContactPersonEmail,
		Description:         in.Description,
		CreateTime:          model.FormatTime(in.CreateTime),
		SubmitForReviewTime: model.FormatTime(in.SubmitForReviewTime),
		ApproveTime:         model.FormatTime(in.ApproveTime),
		UpdateTime:          model.FormatTime(in.UpdateTime),

		Files:       s.convertFiles(in.Files),
		CompanyInfo: s.convertCompanyInfo(in.CompanyInfo),
	}

	return out
}

func (s *MerchantService) convertCompanyInfo(in *model.Company) *system.CompanyInfo {
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

func (s *MerchantService) convertFiles(in []*model.File) (out []*system.FileInfo) {
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
