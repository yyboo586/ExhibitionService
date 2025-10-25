package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/model"
	"context"
)

// SubmitExhibitionForReview 提交展会审核
func (s *ExhibitionService) SubmitExhibitionForReview(ctx context.Context, req *system.SubmitExhibitionForReviewReq) (res *system.SubmitExhibitionForReviewRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, model.ExhibitionEventSubmitForReview, nil)
	if err != nil {
		return nil, err
	}
	return &system.SubmitExhibitionForReviewRes{}, nil
}

// GetSubmitForReviewList 获取待审核列表
func (s *ExhibitionService) GetSubmitForReviewList(ctx context.Context, req *system.GetPendingListReq) (res *system.GetPendingListRes, err error) {
	exhibitions, pageRes, err := s.exhibitionDomain.GetPendingList(ctx, &req.PageReq)
	if err != nil {
		return nil, err
	}

	res = &system.GetPendingListRes{PageRes: pageRes}
	for _, v := range exhibitions {
		res.List = append(res.List, s.convertExhibition(v))
	}
	return res, nil
}

// ApproveExhibition 审核通过展会
func (s *ExhibitionService) ApproveExhibition(ctx context.Context, req *system.ApproveExhibitionReq) (res *system.ApproveExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, model.ExhibitionEventApprove, nil)
	if err != nil {
		return nil, err
	}
	return &system.ApproveExhibitionRes{}, nil
}

// RejectExhibition 审核驳回展会
func (s *ExhibitionService) RejectExhibition(ctx context.Context, req *system.RejectExhibitionReq) (res *system.RejectExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, model.ExhibitionEventReject, nil)
	if err != nil {
		return nil, err
	}
	return &system.RejectExhibitionRes{}, nil
}

// CancelExhibition 取消展会
func (s *ExhibitionService) CancelExhibition(ctx context.Context, req *system.CancelExhibitionReq) (res *system.CancelExhibitionRes, err error) {
	err = s.exhibitionDomain.HandleEvent(ctx, req.ID, model.ExhibitionEventCancel, nil)
	if err != nil {
		return nil, err
	}
	return &system.CancelExhibitionRes{}, nil
}
