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
}

func NewMerchantService(merchantDomain interfaces.IMerchant, fileDomain interfaces.IFile) *MerchantService {
	return &MerchantService{
		merchantDomain: merchantDomain,
		fileDomain:     fileDomain,
	}
}

// CreateMerchant 创建展商
func (s *MerchantService) CreateMerchant(ctx context.Context, req *system.CreateMerchantReq) (res *system.CreateMerchantRes, err error) {
	var merchantID string

	// 创建展商
	merchantInfo := &model.Merchant{
		CompanyID:          req.CompanyID,
		ExhibitionID:       req.ExhibitionID,
		Name:               req.Name,
		Description:        req.Description,
		BoothNumber:        req.BoothNumber,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonEmail: req.ContactPersonEmail,
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建展商
		merchantID, err = s.merchantDomain.Create(ctx, tx, merchantInfo)
		if err != nil {
			return err
		}

		// 更新展商文件
		if len(req.Files) > 0 {
			for _, v := range req.Files {
				err := s.fileDomain.UpdateFileCustomInfo(ctx, tx, v.FileID, model.FileModuleMerchant, merchantID, model.GetFileType(v.FileType))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.CreateMerchantRes{ID: merchantID}, nil
}

// GetMerchant 获取展商详情
func (s *MerchantService) GetMerchant(ctx context.Context, req *system.GetMerchantReq) (res *system.GetMerchantRes, err error) {
	var (
		merchant  *model.Merchant
		fileInfos []*model.File
	)

	merchant, err = s.merchantDomain.GetMerchant(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 获取展商文件
	fileInfos, err = s.fileDomain.ListFilesByModuleAndCustomID(ctx, model.FileModuleMerchant, merchant.ID)
	if err != nil {
		return nil, err
	}

	res = &system.GetMerchantRes{
		MerchantInfo: s.convertMerchant(merchant),
		Files:        s.convertFiles(fileInfos),
	}

	return res, nil
}

// ListMerchants 列表展商
func (s *MerchantService) ListMerchants(ctx context.Context, req *model.ListMerchantsReq) (res *model.ListMerchantsRes, err error) {
	merchants, pageRes, err := s.merchantDomain.ListMerchants(ctx, req.ExhibitionID, req.Name, req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &model.ListMerchantsRes{
		Merchants: merchants,
		PageRes:   pageRes,
	}

	return res, nil
}

// GetPendingMerchants 获取待审核展商列表
func (s *MerchantService) GetPendingMerchants(ctx context.Context, req *model.PageReq) (res *model.ListMerchantsRes, err error) {
	merchants, pageRes, err := s.merchantDomain.GetPendingList(ctx, req)
	if err != nil {
		return nil, err
	}

	res = &model.ListMerchantsRes{
		Merchants: merchants,
		PageRes:   pageRes,
	}

	return res, nil
}

// ApproveMerchant 审核通过展商
func (s *MerchantService) ApproveMerchant(ctx context.Context, req *system.ApproveMerchantReq) (res *system.ApproveMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, interfaces.MerchantEventApprove, nil)
	if err != nil {
		return nil, err
	}
	return &system.ApproveMerchantRes{}, nil
}

// RejectMerchant 审核拒绝展商
func (s *MerchantService) RejectMerchant(ctx context.Context, req *system.RejectMerchantReq) (res *system.RejectMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, interfaces.MerchantEventReject, nil)
	if err != nil {
		return nil, err
	}
	return &system.RejectMerchantRes{}, nil
}

// DisableMerchant 禁用展商
func (s *MerchantService) DisableMerchant(ctx context.Context, req *system.DisableMerchantReq) (res *system.DisableMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, interfaces.MerchantEventDisable, nil)
	if err != nil {
		return nil, err
	}
	return &system.DisableMerchantRes{}, nil
}

// EnableMerchant 启用展商
func (s *MerchantService) EnableMerchant(ctx context.Context, req *system.EnableMerchantReq) (res *system.EnableMerchantRes, err error) {
	err = s.merchantDomain.HandleEvent(ctx, req.ID, interfaces.MerchantEventEnable, nil)
	if err != nil {
		return nil, err
	}
	return &system.EnableMerchantRes{}, nil
}

func (s *MerchantService) convertMerchant(in *model.Merchant) system.MerchantInfo {
	return system.MerchantInfo{
		ID:                 in.ID,
		CompanyID:          in.CompanyID,
		ExhibitionID:       in.ExhibitionID,
		Name:               in.Name,
		Status:             model.GetMerchantStatusText(in.Status),
		Description:        in.Description,
		BoothNumber:        in.BoothNumber,
		ContactPersonName:  in.ContactPersonName,
		ContactPersonPhone: in.ContactPersonPhone,
		ContactPersonEmail: in.ContactPersonEmail,
		CreateTime:         model.FormatTime(in.CreateTime),
		UpdateTime:         model.FormatTime(in.UpdateTime),
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
