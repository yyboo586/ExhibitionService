package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type ExhibitionService struct {
	exhibitionDomain      interfaces.IExhibition
	serviceProviderDomain interfaces.IServiceProvider
}

func NewExhibitionService(exhibitionDomain interfaces.IExhibition, serviceProviderDomain interfaces.IServiceProvider) *ExhibitionService {
	return &ExhibitionService{
		exhibitionDomain:      exhibitionDomain,
		serviceProviderDomain: serviceProviderDomain,
	}
}

// CreateExhibition 创建展会
func (s *ExhibitionService) CreateExhibition(ctx context.Context, req *system.CreateExhibitionReq) (res *system.CreateExhibitionRes, err error) {
	var (
		exhibitionID string
		exOrganizers []*model.ExOrganizer
		fileInfos    []*model.File
	)
	for _, v := range req.Organizers {
		serviceProvider, err := s.serviceProviderDomain.Get(ctx, v.ServiceProviderID)
		if err != nil {
			return nil, err
		}
		err = s.serviceProviderDomain.IsAvailable(ctx, serviceProvider)
		if err != nil {
			return nil, err
		}

		exOrganizers = append(exOrganizers, &model.ExOrganizer{
			ServiceProviderID: v.ServiceProviderID,
			RoleType:          model.GetOrganizerRole(v.RoleType),
		})
	}
	for _, v := range req.Files {
		fileInfos = append(fileInfos, &model.File{
			FileID: v.FileID,
			Type:   model.GetFileType(v.FileType),
		})
	}

	// 创建展会
	exhibitionInfo := &model.Exhibition{
		Organizers:   exOrganizers,
		Title:        req.Title,
		Website:      req.Website,
		Industry:     req.Industry,
		Tags:         req.Tags,
		Country:      req.Country,
		City:         req.City,
		Venue:        req.Venue,
		VenueAddress: req.VenueAddress,
		Description:  req.Description,

		RegistrationStart: req.RegistrationStart,
		RegistrationEnd:   req.RegistrationEnd,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,

		Files: fileInfos,
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建展会
		exhibitionID, err = s.exhibitionDomain.Create(ctx, tx, exhibitionInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.CreateExhibitionRes{ID: exhibitionID}, nil
}

// GetExhibition 获取展会详情
func (s *ExhibitionService) GetExhibition(ctx context.Context, req *system.GetExhibitionReq) (res *system.GetExhibitionRes, err error) {
	var (
		exhibition *model.Exhibition
	)

	exhibition, err = s.exhibitionDomain.GetExhibition(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res = &system.GetExhibitionRes{
		ExhibitionInfo: s.convertExhibition(exhibition),
	}
	return res, nil
}

// ListExhibitions 列表展会
func (s *ExhibitionService) ListExhibitions(ctx context.Context, req *system.ListExhibitionsReq) (res *system.ListExhibitionsRes, err error) {
	exhibitions, pageRes, err := s.exhibitionDomain.ListExhibitions(ctx, req.Name, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.ListExhibitionsRes{PageRes: pageRes}
	for _, v := range exhibitions {
		res.List = append(res.List, s.convertExhibition(v))
	}
	return res, nil
}

func (s *ExhibitionService) convertExhibition(in *model.Exhibition) *system.ExhibitionInfo {
	return &system.ExhibitionInfo{
		ID:           in.ID,
		Title:        in.Title,
		Website:      in.Website,
		Status:       model.GetExhibitionStatusText(in.Status),
		Industry:     in.Industry,
		Tags:         in.Tags,
		Country:      in.Country,
		City:         in.City,
		Venue:        in.Venue,
		VenueAddress: in.VenueAddress,
		Description:  in.Description,

		RegistrationStart: model.FormatTime(in.RegistrationStart),
		RegistrationEnd:   model.FormatTime(in.RegistrationEnd),
		StartTime:         model.FormatTime(in.StartTime),
		EndTime:           model.FormatTime(in.EndTime),

		CreateTime:          model.FormatTime(in.CreateTime),
		SubmitForReviewTime: model.FormatTime(in.SubmitForReviewTime),
		ApproveTime:         model.FormatTime(in.ApproveTime),
		UpdateTime:          model.FormatTime(in.UpdateTime),

		Files:      s.convertFiles(in.Files),
		Organizers: s.convertOrganizers(in.Organizers),
	}
}

func (s *ExhibitionService) convertOrganizers(in []*model.ExOrganizer) (out []*system.OrganizerInfo) {
	for _, v := range in {
		out = append(out, &system.OrganizerInfo{
			ServiceProviderID: v.ServiceProviderID,
			RoleType:          model.GetOrganizerRoleText(v.RoleType),
		})
	}
	return out
}

func (s *ExhibitionService) convertFiles(in []*model.File) (out []*system.FileInfo) {
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

// ApplyForExhibition 商户申请参加展会
func (s *ExhibitionService) ApplyForExhibition(ctx context.Context, req *system.ApplyForExhibitionReq) (res *system.ApplyForExhibitionRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = s.exhibitionDomain.CreateApplyForExhibition(ctx, tx, req.ExhibitionID, req.MerchantID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &system.ApplyForExhibitionRes{Message: "申请提交成功，等待审核"}, nil
}

// GetMerchantApplication 获取商户申请状态
func (s *ExhibitionService) GetMerchantApplication(ctx context.Context, req *system.GetMerchantApplicationReq) (res *system.GetMerchantApplicationRes, err error) {
	application, err := s.exhibitionDomain.GetExMerchantApplication(ctx, req.ExhibitionID, req.MerchantID)
	if err != nil {
		return nil, err
	}

	res = &system.GetMerchantApplicationRes{
		ExhibitionMerchantInfo: s.convertExhibitionMerchant(application),
	}
	return res, nil
}

// ListExhibitionApplications 获取展会申请列表
func (s *ExhibitionService) ListExhibitionApplications(ctx context.Context, req *system.ListExhibitionApplicationsReq) (res *system.ListExhibitionApplicationsRes, err error) {
	applications, pageRes, err := s.exhibitionDomain.ListExhibitionApplications(ctx, req.ExhibitionID, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.ListExhibitionApplicationsRes{PageRes: pageRes}
	for _, v := range applications {
		res.List = append(res.List, s.convertExhibitionMerchant(v))
	}
	return res, nil
}

// ListMerchantApplications 获取商户申请列表
func (s *ExhibitionService) ListMerchantApplications(ctx context.Context, req *system.ListMerchantApplicationsReq) (res *system.ListMerchantApplicationsRes, err error) {
	applications, pageRes, err := s.exhibitionDomain.ListMerchantApplications(ctx, req.MerchantID, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.ListMerchantApplicationsRes{PageRes: pageRes}
	for _, v := range applications {
		res.List = append(res.List, s.convertExhibitionMerchant(v))
	}
	return res, nil
}

func (s *ExhibitionService) convertExhibitionMerchant(in *model.ExhibitionMerchant) *system.ExhibitionMerchantInfo {
	return &system.ExhibitionMerchantInfo{
		ID:                  in.ID,
		ExhibitionID:        in.ExhibitionID,
		MerchantID:          in.MerchantID,
		Status:              model.GetExMerchantStatusText(in.Status),
		CreateTime:          model.FormatTime(in.CreateTime),
		SubmitForReviewTime: model.FormatTime(in.SubmitForReviewTime),
		ApproveTime:         model.FormatTime(in.ApproveTime),
		UpdateTime:          model.FormatTime(in.UpdateTime),
	}
}
