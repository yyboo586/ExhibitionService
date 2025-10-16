package service

import (
	"ExhibitionService/internal/interfaces"
)

type ExhibitionService struct {
	exhibitionDomain interfaces.IExhibition
	fileDomain       interfaces.IFile
}

func NewExhibitionService(exhibitionDomain interfaces.IExhibition, fileDomain interfaces.IFile) *ExhibitionService {
	return &ExhibitionService{
		exhibitionDomain: exhibitionDomain,
		fileDomain:       fileDomain,
	}
}

/*
// CreateExhibition 创建展会
func (s *ExhibitionService) CreateExhibition(ctx context.Context, req *system.CreateExhibitionReq) (res *system.CreateExhibitionRes, err error) {
	var exhibitionID string

	// 创建展会
	exhibitionInfo := &model.Exhibition{
		ServiceProviderID: req.ServiceProviderID,
		Title:             req.Title,
		Industry:          req.Industry,
		Tags:              req.Tags,
		Website:           req.Website,
		Venue:             req.Venue,
		VenueAddress:      req.VenueAddress,
		Country:           req.Country,
		City:              req.City,
		Description:       req.Description,
		RegistrationStart: req.RegistrationStart,
		RegistrationEnd:   req.RegistrationEnd,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建展会
		exhibitionID, err = s.exhibitionDomain.Create(ctx, tx, exhibitionInfo)
		if err != nil {
			return err
		}

		// 更新展会文件
		if len(req.Files) > 0 {
			for _, v := range req.Files {
				err := s.fileDomain.UpdateFileCustomInfo(ctx, tx, v.FileID, model.FileModuleExhibition, exhibitionID, model.GetFileType(v.FileType))
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

	return &system.CreateExhibitionRes{ID: exhibitionID}, nil
}

// GetExhibition 获取展会详情
func (s *ExhibitionService) GetExhibition(ctx context.Context, req *system.GetExhibitionReq) (res *system.GetExhibitionRes, err error) {
	var (
		exhibition           *model.Exhibition
		exhibitionOrganizers []*model.ExhibitionOrganizer
		organizers           []*model.Organizer
		fileInfos            []*model.File
	)

	exhibition, err = s.exhibitionDomain.GetExhibition(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 获取主办方信息
	exhibitionOrganizers, err = s.organizerDomain.GetExhibitionOrganizers(ctx, req.ID)
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
	fileInfos, err = s.fileDomain.ListFilesByModuleAndCustomID(ctx, model.FileModuleExhibition, exhibition.ID)
	if err != nil {
		return nil, err
	}

	res = &system.GetExhibitionRes{
		ExhibitionUnit: &system.ExhibitionUnit{
			ExhibitionInfo: s.convertExhibition(exhibition),
			Files:          s.convertFiles(fileInfos),
			Organizers:     s.convertOrganizers(organizers),
		},
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

func (s *ExhibitionService) convertExhibition(in *model.Exhibition) system.ExhibitionInfo {
	return system.ExhibitionInfo{
		ID:                in.ID,
		ServiceProviderID: in.ServiceProviderID,
		Title:             in.Title,
		Status:            model.GetExhibitionStatusText(in.Status),
		Industry:          in.Industry,
		Tags:              in.Tags,
		Website:           in.Website,
		Venue:             in.Venue,
		VenueAddress:      in.VenueAddress,
		Country:           in.Country,
		City:              in.City,
		Description:       in.Description,
		RegistrationStart: model.FormatTime(in.RegistrationStart),
		RegistrationEnd:   model.FormatTime(in.RegistrationEnd),
		StartTime:         model.FormatTime(in.StartTime),
		EndTime:           model.FormatTime(in.EndTime),
		CreateTime:        model.FormatTime(in.CreateTime),
		UpdateTime:        model.FormatTime(in.UpdateTime),
	}
}

func (s *ExhibitionService) convertOrganizers(in []*model.Organizer) (out []*system.OrganizerInfo) {
	for _, v := range in {
		out = append(out, &system.OrganizerInfo{
			CompanyID:          v.CompanyID,
			Name:               v.Name,
			RoleType:           model.GetOrganizerRoleText(v.RoleType),
			Description:        v.Description,
			ContactPersonName:  v.ContactPersonName,
			ContactPersonPhone: v.ContactPersonPhone,
			ContactPersonEmail: v.ContactPersonEmail,
			Website:            v.Website,
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
*/
