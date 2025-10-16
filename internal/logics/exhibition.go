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
	exhibitionOnce   sync.Once
	exhibitionDomain *exhibition
)

type exhibition struct{}

func NewExhibition() interfaces.IExhibition {
	exhibitionOnce.Do(func() {
		exhibitionDomain = &exhibition{}
	})
	return exhibitionDomain
}

func (e *exhibition) Create(ctx context.Context, tx gdb.TX, in *model.Exhibition) (id string, err error) {
	// 时间逻辑校验
	err = validateExhibitionTime(in.RegistrationStart, in.RegistrationEnd, in.StartTime, in.EndTime)
	if err != nil {
		return "", gerror.Newf("create exhibition failed, %s", err.Error())
	}

	id = uuid.New().String()
	// 创建展会
	data := map[string]any{
		dao.Exhibition.Columns().ID:                id,
		dao.Exhibition.Columns().ServiceProviderID: in.ServiceProviderID,
		dao.Exhibition.Columns().Title:             in.Title,
		dao.Exhibition.Columns().Status:            int(model.ExhibitionStatusPreparing),
		dao.Exhibition.Columns().Industry:          in.Industry,
		dao.Exhibition.Columns().Tags:              in.Tags,
		dao.Exhibition.Columns().Website:           in.Website,
		dao.Exhibition.Columns().Venue:             in.Venue,
		dao.Exhibition.Columns().VenueAddress:      in.VenueAddress,
		dao.Exhibition.Columns().Country:           in.Country,
		dao.Exhibition.Columns().City:              in.City,
		dao.Exhibition.Columns().Description:       in.Description,
		dao.Exhibition.Columns().RegistrationStart: in.RegistrationStart.Unix(),
		dao.Exhibition.Columns().RegistrationEnd:   in.RegistrationEnd.Unix(),
		dao.Exhibition.Columns().StartTime:         in.StartTime.Unix(),
		dao.Exhibition.Columns().EndTime:           in.EndTime.Unix(),
		dao.Exhibition.Columns().Version:           1,
		dao.Exhibition.Columns().CreateTime:        time.Now().Unix(),
		dao.Exhibition.Columns().UpdateTime:        time.Now().Unix(),
	}

	_, err = dao.Exhibition.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return "", gerror.Newf("create exhibition failed, exhibition name already exists, name: %s", in.Title)
		}
		return "", gerror.Newf("create exhibition failed, err: %s", err.Error())
	}

	return id, nil
}

func (e *exhibition) GetExhibition(ctx context.Context, id string) (out *model.Exhibition, err error) {
	var te entity.TExhibition
	err = dao.Exhibition.Ctx(ctx).Where(dao.Exhibition.Columns().ID, id).Scan(&te)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get exhibition failed, exhibition not found, id: %s", id)
		}
		return nil, gerror.Newf("get exhibition failed, err: %s", err.Error())
	}

	return model.ConvertExhibition(&te), nil
}

func (e *exhibition) UpdateExhibition(ctx context.Context, in *model.Exhibition) (err error) {
	data := map[string]any{}
	if in.Title != "" {
		data[dao.Exhibition.Columns().Title] = in.Title
	}
	if in.Industry != "" {
		data[dao.Exhibition.Columns().Industry] = in.Industry
	}
	if in.Tags != "" {
		data[dao.Exhibition.Columns().Tags] = in.Tags
	}
	if in.Website != "" {
		data[dao.Exhibition.Columns().Website] = in.Website
	}
	if in.Venue != "" {
		data[dao.Exhibition.Columns().Venue] = in.Venue
	}
	if in.VenueAddress != "" {
		data[dao.Exhibition.Columns().VenueAddress] = in.VenueAddress
	}
	if in.Country != "" {
		data[dao.Exhibition.Columns().Country] = in.Country
	}
	if in.City != "" {
		data[dao.Exhibition.Columns().City] = in.City
	}
	if in.Description != "" {
		data[dao.Exhibition.Columns().Description] = in.Description
	}
	if !in.RegistrationStart.IsZero() {
		data[dao.Exhibition.Columns().RegistrationStart] = in.RegistrationStart.Unix()
	}
	if !in.RegistrationEnd.IsZero() {
		data[dao.Exhibition.Columns().RegistrationEnd] = in.RegistrationEnd.Unix()
	}
	if !in.StartTime.IsZero() {
		data[dao.Exhibition.Columns().StartTime] = in.StartTime.Unix()
	}
	if !in.EndTime.IsZero() {
		data[dao.Exhibition.Columns().EndTime] = in.EndTime.Unix()
	}

	if len(data) == 0 {
		return nil
	}

	data[dao.Exhibition.Columns().UpdateTime] = time.Now().Unix()
	_, err = dao.Exhibition.Ctx(ctx).Data(data).Where(dao.Exhibition.Columns().ID, in.ID).Update()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return gerror.Newf("update exhibition failed, exhibition name already exists, name: %s", in.Title)
		}
		return gerror.Newf("update exhibition failed, err: %s", err.Error())
	}

	return nil
}

func (e *exhibition) DeleteExhibition(ctx context.Context, id string) (err error) {
	_, err = dao.Exhibition.Ctx(ctx).Where(dao.Exhibition.Columns().ID, id).Delete()
	if err != nil {
		return gerror.Newf("delete exhibition failed, err: %s", err.Error())
	}
	return nil
}

func (e *exhibition) ListExhibitions(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Exhibition.Ctx(ctx)
	if name != "" {
		query = query.WhereLike(dao.Exhibition.Columns().Title, name+"%")
	}
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("list exhibitions failed, query count err: %s", err.Error())
	}

	var te []*entity.TExhibition
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Exhibition.Columns().CreateTime)
	err = query.Scan(&te)
	if err != nil {
		return nil, nil, gerror.Newf("list exhibitions failed, query scan err: %s", err.Error())
	}

	for _, r := range te {
		out = append(out, model.ConvertExhibition(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

func (e *exhibition) GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Exhibition.Ctx(ctx).Where(dao.Exhibition.Columns().Status, int(model.ExhibitionStatusPending))
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("get submit for review list failed, query count err: %s", err.Error())
	}

	var te []*entity.TExhibition
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Exhibition.Columns().CreateTime)
	err = query.Scan(&te)
	if err != nil {
		return nil, nil, gerror.Newf("get submit for review list failed, query scan err: %s", err.Error())
	}

	for _, r := range te {
		out = append(out, model.ConvertExhibition(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

// 私有方法
/*
报名开始时间 > 报名结束时间	"registration start time must be before registration end time"
报名结束时间 >= 展会开始时间	"registration end time must be strictly before exhibition start time"
展会开始时间 >= 展会结束时间	"exhibition start time must be strictly before exhibition end time"
报名开始时间 < 当前时间	"registration start time must be after current time"
报名开始时间 < (当前时间+30分钟)	"registration start time must be at least 30 minutes from now"
展会开始时间 > (报名结束时间+7天)	"exhibition must start within 7 days after registration ends"
展会结束时间 > (展会开始时间+30天)	"exhibition duration cannot exceed 30 days"
*/
func validateExhibitionTime(registrationStart time.Time, registrationEnd time.Time, startTime time.Time, endTime time.Time) error {
	now := time.Now()

	// 1. 报名开始时间不能大于报名结束时间
	if registrationStart.After(registrationEnd) {
		return gerror.New("registration start time must be before registration end time")
	}

	// 2. 报名结束时间不能大于等于展会开始时间
	if !registrationEnd.Before(startTime) {
		return gerror.New("registration end time must be strictly before exhibition start time")
	}

	// 3. 展会开始时间不能大于等于展会结束时间
	if !startTime.Before(endTime) {
		return gerror.New("exhibition start time must be strictly before exhibition end time")
	}

	// 4. 报名开始时间不能小于当前时间
	if registrationStart.Before(now) {
		return gerror.New("registration start time must be after current time")
	}

	// 5. 报名开始时间必须至少在当前时间30分钟后
	minRegistrationStart := now.Add(30 * time.Minute)
	if registrationStart.Before(minRegistrationStart) {
		return gerror.New("registration start time must be at least 30 minutes from now")
	}

	// 6. 展会开始时间必须在报名结束时间7天内
	maxStartTime := registrationEnd.Add(7 * 24 * time.Hour)
	if startTime.After(maxStartTime) {
		return gerror.New("exhibition must start within 7 days after registration ends")
	}

	// 7. 展会持续时间不能超过30天
	maxEndTime := startTime.Add(30 * 24 * time.Hour)
	if endTime.After(maxEndTime) {
		return gerror.New("exhibition duration cannot exceed 30 days")
	}

	return nil
}
