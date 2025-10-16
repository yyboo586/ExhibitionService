package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/google/uuid"
)

var (
	merchantOnce   sync.Once
	merchantDomain *merchant
)

type merchant struct{}

func NewMerchant() interfaces.IMerchant {
	merchantOnce.Do(func() {
		merchantDomain = &merchant{}
	})
	return merchantDomain
}

func (m *merchant) Create(ctx context.Context, tx gdb.TX, in *model.Merchant) (id string, err error) {
	id = uuid.New().String()

	// 创建展商
	data := map[string]any{
		dao.Merchant.Columns().ID:                 id,
		dao.Merchant.Columns().CompanyID:          in.CompanyID,
		dao.Merchant.Columns().ExhibitionID:       in.ExhibitionID,
		dao.Merchant.Columns().Name:               in.Name,
		dao.Merchant.Columns().Description:        in.Description,
		dao.Merchant.Columns().BoothNumber:        in.BoothNumber,
		dao.Merchant.Columns().ContactPersonName:  in.ContactPersonName,
		dao.Merchant.Columns().ContactPersonPhone: in.ContactPersonPhone,
		dao.Merchant.Columns().ContactPersonEmail: in.ContactPersonEmail,
		dao.Merchant.Columns().Status:             int(model.MerchantStatusPending),
		dao.Merchant.Columns().Version:            1,
		dao.Merchant.Columns().CreateTime:         time.Now().Unix(),
		dao.Merchant.Columns().UpdateTime:         time.Now().Unix(),
	}

	_, err = dao.Merchant.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return "", gerror.Newf("create merchant failed, merchant already exists for this exhibition")
		}
		return "", gerror.Newf("create merchant failed, err: %s", err.Error())
	}

	return id, nil
}

func (m *merchant) GetMerchant(ctx context.Context, id string) (out *model.Merchant, err error) {
	var tm entity.TMerchant
	err = dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().ID, id).Scan(&tm)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get merchant failed, merchant not found, id: %s", id)
		}
		return nil, gerror.Newf("get merchant failed, err: %s", err.Error())
	}

	return model.ConvertMerchant(&tm), nil
}

func (m *merchant) UpdateMerchant(ctx context.Context, in *model.Merchant) (err error) {
	data := map[string]any{}
	if in.Name != "" {
		data[dao.Merchant.Columns().Name] = in.Name
	}
	if in.Description != "" {
		data[dao.Merchant.Columns().Description] = in.Description
	}
	if in.BoothNumber != "" {
		data[dao.Merchant.Columns().BoothNumber] = in.BoothNumber
	}
	if in.ContactPersonName != "" {
		data[dao.Merchant.Columns().ContactPersonName] = in.ContactPersonName
	}
	if in.ContactPersonPhone != "" {
		data[dao.Merchant.Columns().ContactPersonPhone] = in.ContactPersonPhone
	}
	if in.ContactPersonEmail != "" {
		data[dao.Merchant.Columns().ContactPersonEmail] = in.ContactPersonEmail
	}

	if len(data) == 0 {
		return nil
	}

	data[dao.Merchant.Columns().UpdateTime] = time.Now().Unix()
	_, err = dao.Merchant.Ctx(ctx).Data(data).Where(dao.Merchant.Columns().ID, in.ID).Update()
	if err != nil {
		return gerror.Newf("update merchant failed, err: %s", err.Error())
	}

	return nil
}

func (m *merchant) DeleteMerchant(ctx context.Context, id string) (err error) {
	_, err = dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().ID, id).Delete()
	if err != nil {
		return gerror.Newf("delete merchant failed, err: %s", err.Error())
	}
	return nil
}

func (m *merchant) ListMerchants(ctx context.Context, exhibitionID string, name string, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Merchant.Ctx(ctx)
	if exhibitionID != "" {
		query = query.Where(dao.Merchant.Columns().ExhibitionID, exhibitionID)
	}
	if name != "" {
		query = query.WhereLike(dao.Merchant.Columns().Name, name+"%")
	}
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("list merchants failed, query count err: %s", err.Error())
	}

	var tm []*entity.TMerchant
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Merchant.Columns().CreateTime)
	err = query.Scan(&tm)
	if err != nil {
		return nil, nil, gerror.Newf("list merchants failed, query scan err: %s", err.Error())
	}

	for _, r := range tm {
		out = append(out, model.ConvertMerchant(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

func (m *merchant) GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().Status, int(model.MerchantStatusPending))
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("get pending merchants failed, query count err: %s", err.Error())
	}

	var tm []*entity.TMerchant
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Merchant.Columns().CreateTime)
	err = query.Scan(&tm)
	if err != nil {
		return nil, nil, gerror.Newf("get pending merchants failed, query scan err: %s", err.Error())
	}

	for _, r := range tm {
		out = append(out, model.ConvertMerchant(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

// 状态流转
func (m *merchant) HandleEvent(ctx context.Context, merchantID string, event interfaces.MerchantEvent, data interface{}) (err error) {
	// 获取展商信息
	merchant, err := m.GetMerchant(ctx, merchantID)
	if err != nil {
		return err
	}

	// 验证当前状态是否支持该事件
	transition, ok := merchantTransitionMap[merchant.Status][event]
	if !ok {
		return gerror.Newf("merchant current status: %s, not supported event: %s", model.GetMerchantStatusText(merchant.Status), GetMerchantEventText(event))
	}

	// 执行状态操作
	err = transition.Action(ctx, merchant, data)
	if err != nil {
		return gerror.Newf("handle event failed, merchant id: %s, event: %s, err: %s", merchantID, GetMerchantEventText(event), err.Error())
	}

	return nil
}
