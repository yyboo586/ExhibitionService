package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"database/sql"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ApplyForExhibition 商户申请参加展会
func (e *exhibition) CreateApplyForExhibition(ctx context.Context, tx gdb.TX, exhibitionID string, merchantID string) (err error) {
	// 检查展会是否存在且状态允许申请
	exhibition, err := e.GetExhibition(ctx, exhibitionID)
	if err != nil {
		return gerror.Newf("apply for exhibition failed, get exhibition failed, err: %s", err.Error())
	}

	// 检查展会状态是否允许申请
	if exhibition.Status != model.ExhibitionStatusEnrolling {
		return gerror.Newf("apply for exhibition failed, exhibition status is not enrolling, current status: %s", model.GetExhibitionStatusText(exhibition.Status))
	}

	// 检查商户是否存在且状态正常
	merchant, err := e.merchantDomain.Get(ctx, merchantID)
	if err != nil {
		return gerror.Newf("apply for exhibition failed, get merchant failed, err: %s", err.Error())
	}

	// 检查商户状态是否允许申请
	if merchant.Status != model.MerchantStatusApproved {
		return gerror.Newf("apply for exhibition failed, merchant status is not approved, current status: %s", model.GetMerchantStatusText(merchant.Status))
	}

	// 检查是否已经申请过
	existingApplication, err := e.GetExMerchantApplication(ctx, exhibitionID, merchantID)
	if err != nil && !gerror.HasCode(err, gcode.CodeNotFound) {
		return gerror.Newf("apply for exhibition failed, check existing application failed, err: %s", err.Error())
	}

	if existingApplication != nil {
		// 如果已经申请过，检查状态
		switch existingApplication.Status {
		case model.ExMerchantStatusPending:
			return gerror.New("apply for exhibition failed, application is already pending")
		case model.ExMerchantStatusApproved:
			return gerror.New("apply for exhibition failed, application is already approved")
		case model.ExMerchantStatusRejected:
			// 如果之前被拒绝，可以重新申请
			break
		case model.ExMerchantStatusWithdrawn:
			// 如果之前退出，可以重新申请
			break
		}
	}

	// 创建申请记录
	data := map[string]any{
		dao.ExhibitionMerchant.Columns().ExhibitionID:        exhibitionID,
		dao.ExhibitionMerchant.Columns().MerchantID:          merchantID,
		dao.ExhibitionMerchant.Columns().Status:              int(model.ExMerchantStatusPending),
		dao.ExhibitionMerchant.Columns().CreateTime:          time.Now().Unix(),
		dao.ExhibitionMerchant.Columns().SubmitForReviewTime: time.Now().Unix(),
		dao.ExhibitionMerchant.Columns().UpdateTime:          time.Now().Unix(),
	}

	_, err = dao.ExhibitionMerchant.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		return gerror.Newf("apply for exhibition failed, insert application failed, err: %s", err.Error())
	}

	return nil
}

// GetMerchantApplication 获取商户在展会中的申请状态
func (e *exhibition) GetExMerchantApplication(ctx context.Context, exhibitionID string, merchantID string) (out *model.ExhibitionMerchant, err error) {
	var tem entity.TExhibitionMerchant
	err = dao.ExhibitionMerchant.Ctx(ctx).
		Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID).
		Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID).
		Scan(&tem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.NewCode(gcode.CodeNotFound, "merchant application not found")
		}
		return nil, gerror.Newf("get merchant application failed, err: %s", err.Error())
	}

	out = model.ConvertExhibitionMerchant(&tem)
	return out, nil
}

// ListExhibitionApplications 获取展会的商户申请列表
func (e *exhibition) ListExhibitionApplications(ctx context.Context, exhibitionID string, pageReq *model.PageReq) (out []*model.ExhibitionMerchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.ExhibitionMerchant.Ctx(ctx).Where(dao.ExhibitionMerchant.Columns().ExhibitionID, exhibitionID)
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("list exhibition applications failed, query count err: %s", err.Error())
	}

	var tem []*entity.TExhibitionMerchant
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.ExhibitionMerchant.Columns().CreateTime)
	err = query.Scan(&tem)
	if err != nil {
		return nil, nil, gerror.Newf("list exhibition applications failed, query scan err: %s", err.Error())
	}

	for _, r := range tem {
		out = append(out, model.ConvertExhibitionMerchant(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

// ListMerchantApplications 获取商户的展会申请列表
func (e *exhibition) ListMerchantApplications(ctx context.Context, merchantID string, pageReq *model.PageReq) (out []*model.ExhibitionMerchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.ExhibitionMerchant.Ctx(ctx).Where(dao.ExhibitionMerchant.Columns().MerchantID, merchantID)
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("list merchant applications failed, query count err: %s", err.Error())
	}

	var tem []*entity.TExhibitionMerchant
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.ExhibitionMerchant.Columns().CreateTime)
	err = query.Scan(&tem)
	if err != nil {
		return nil, nil, gerror.Newf("list merchant applications failed, query scan err: %s", err.Error())
	}

	for _, r := range tem {
		out = append(out, model.ConvertExhibitionMerchant(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}
