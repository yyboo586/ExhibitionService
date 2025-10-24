package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/model"
	"context"
)

// GetPendingMerchants 获取待审核商户列表
func (s *MerchantService) GetPendingMerchants(ctx context.Context, req *system.GetPendingMerchantsReq) (res *system.GetPendingMerchantsRes, err error) {
	pageReq := &model.PageReq{
		Page: req.Page,
		Size: req.Size,
	}

	merchants, pageRes, err := s.merchantDomain.GetPendingList(ctx, pageReq)
	if err != nil {
		return nil, err
	}

	res = &system.GetPendingMerchantsRes{
		PageRes: pageRes,
	}
	for _, v := range merchants {
		res.List = append(res.List, s.convertMerchant(v))
	}
	return res, nil
}

// ApproveMerchant 审核通过商户
func (s *MerchantService) ApproveMerchant(ctx context.Context, req *system.ApproveMerchantReq) (res *system.ApproveMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventApprove, nil)
	if err != nil {
		return nil, err
	}
	return &system.ApproveMerchantRes{}, nil
}

// RejectMerchant 审核拒绝商户
func (s *MerchantService) RejectMerchant(ctx context.Context, req *system.RejectMerchantReq) (res *system.RejectMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventReject, nil)
	if err != nil {
		return nil, err
	}
	return &system.RejectMerchantRes{}, nil
}

// DisableMerchant 禁用商户
func (s *MerchantService) DisableMerchant(ctx context.Context, req *system.DisableMerchantReq) (res *system.DisableMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventDisable, nil)
	if err != nil {
		return nil, err
	}
	return &system.DisableMerchantRes{}, nil
}

// EnableMerchant 启用商户
func (s *MerchantService) EnableMerchant(ctx context.Context, req *system.EnableMerchantReq) (res *system.EnableMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventEnable, nil)
	if err != nil {
		return nil, err
	}
	return &system.EnableMerchantRes{}, nil
}

// UnregisterMerchant 注销商户
func (s *MerchantService) UnregisterMerchant(ctx context.Context, req *system.UnregisterMerchantReq) (res *system.UnregisterMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventUnregister, nil)
	if err != nil {
		return nil, err
	}
	return &system.UnregisterMerchantRes{}, nil
}

func (s *MerchantService) Recommit(ctx context.Context, req *system.RecommitMerchantReq) (res *system.RecommitMerchantRes, err error) {
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
	merchantInfo := &model.Merchant{
		ID:                 req.ID,
		Name:               req.Name,
		Website:            req.Website,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonEmail: req.ContactPersonEmail,
		Description:        req.Description,
		CompanyInfo:        companyInfo,
	}
	for _, v := range req.Files {
		merchantInfo.Files = append(merchantInfo.Files, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}

	err = s.merchantDomain.HandleEvent(ctx, req.ID, model.MerchantEventReCommit, merchantInfo)
	if err != nil {
		return nil, err
	}

	return &system.RecommitMerchantRes{}, nil
}
