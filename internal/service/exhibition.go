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

/*
// SubmitExhibitionForReview 提交展会审核
func (s *ExhibitionService) SubmitExhibitionForReview(ctx context.Context, req *system.SubmitExhibitionForReviewReq) (res *system.SubmitExhibitionForReviewRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, interfaces.ExhibitionEventSubmitForReview, nil)
	if err != nil {
		return nil, err
	}
	return &system.SubmitExhibitionForReviewRes{}, nil
}

// GetSubmitForReviewList 获取待审核列表
func (s *ExhibitionService) GetSubmitForReviewList(ctx context.Context, req *system.GetSubmitForReviewListReq) (res *system.GetSubmitForReviewListRes, err error) {
	exhibitions, pageRes, err := s.exhibitionDomain.GetPendingList(ctx, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.GetSubmitForReviewListRes{PageRes: pageRes}
	for _, v := range exhibitions {
		var (
			exhibitionOrganizers []*model.ExhibitionOrganizer
			organizers           []*model.Organizer
			fileInfos            []*model.File
		)

		// 获取主办方信息
		exhibitionOrganizers, err = s.organizerDomain.GetExhibitionOrganizers(ctx, v.ID)
		if err != nil {
			return nil, err
		}

		for _, eo := range exhibitionOrganizers {
			tmp, err := s.organizerDomain.GetOrganizer(ctx, eo.OrganizerID)
			if err != nil {
				return nil, err
			}
			organizers = append(organizers, tmp)
		}

		// 获取展会文件
		fileInfos, err = s.fileDomain.ListFilesByModuleAndCustomID(ctx, model.FileModuleExhibition, v.ID)
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, &system.ExhibitionUnit{
			ExhibitionInfo: s.convertExhibition(v),
			Files:          s.convertFiles(fileInfos),
			Organizers:     s.convertOrganizers(organizers),
		})
	}
	return res, nil
}

// ApproveExhibition 审核通过展会
func (s *ExhibitionService) ApproveExhibition(ctx context.Context, req *system.ApproveExhibitionReq) (res *system.ApproveExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, interfaces.ExhibitionEventApprove, nil)
	if err != nil {
		return nil, err
	}
	return &system.ApproveExhibitionRes{}, nil
}

// RejectExhibition 审核驳回展会
func (s *ExhibitionService) RejectExhibition(ctx context.Context, req *system.RejectExhibitionReq) (res *system.RejectExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, interfaces.ExhibitionEventReject, nil)
	if err != nil {
		return nil, err
	}
	return &system.RejectExhibitionRes{}, nil
}

// CancelExhibition 取消展会
func (s *ExhibitionService) CancelExhibition(ctx context.Context, req *system.CancelExhibitionReq) (res *system.CancelExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, interfaces.ExhibitionEventCancel, nil)
	if err != nil {
		return nil, err
	}
	return &system.CancelExhibitionRes{}, nil
}

*/

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
