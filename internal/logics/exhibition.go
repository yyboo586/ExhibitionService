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

type exhibition struct {
	fileDomain interfaces.IFile
}

func NewExhibition(fileDomain interfaces.IFile) interfaces.IExhibition {
	exhibitionOnce.Do(func() {
		exhibitionDomain = &exhibition{
			fileDomain: fileDomain,
		}
	})
	return exhibitionDomain
}

func (e *exhibition) Create(ctx context.Context, tx gdb.TX, in *model.Exhibition) (id string, err error) {
	// 检查文件是否上传成功
	fileIDs := make([]string, len(in.Files))
	for i, v := range in.Files {
		fileIDs[i] = v.FileID
	}
	err = e.fileDomain.CheckFileUploadSuccess(ctx, fileIDs)
	if err != nil {
		return "", gerror.Newf("create exhibition failed, file upload failed, err: %s", err.Error())
	}

	// 时间逻辑校验
	err = validateExhibitionTime(in.RegistrationStart, in.RegistrationEnd, in.StartTime, in.EndTime)
	if err != nil {
		return "", gerror.Newf("create exhibition failed, %s", err.Error())
	}

	id = uuid.New().String()
	// 创建展会
	data := map[string]any{
		dao.Exhibition.Columns().ID:           id,
		dao.Exhibition.Columns().Title:        in.Title,
		dao.Exhibition.Columns().Website:      in.Website,
		dao.Exhibition.Columns().Status:       int(model.ExhibitionStatusPreparing),
		dao.Exhibition.Columns().Industry:     in.Industry,
		dao.Exhibition.Columns().Tags:         in.Tags,
		dao.Exhibition.Columns().Country:      in.Country,
		dao.Exhibition.Columns().City:         in.City,
		dao.Exhibition.Columns().Venue:        in.Venue,
		dao.Exhibition.Columns().VenueAddress: in.VenueAddress,
		dao.Exhibition.Columns().Description:  in.Description,

		dao.Exhibition.Columns().RegistrationStart: in.RegistrationStart.Unix(),
		dao.Exhibition.Columns().RegistrationEnd:   in.RegistrationEnd.Unix(),
		dao.Exhibition.Columns().StartTime:         in.StartTime.Unix(),
		dao.Exhibition.Columns().EndTime:           in.EndTime.Unix(),

		dao.Exhibition.Columns().CreateTime: time.Now().Unix(),
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}

	_, err = dao.Exhibition.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return "", gerror.Newf("create exhibition failed, exhibition name already exists, name: %s", in.Title)
		}
		return "", gerror.Newf("create exhibition failed, err: %s", err.Error())
	}

	err = e.createFileRelation(ctx, tx, id, in.Files)
	if err != nil {
		return "", gerror.Newf("create exhibition failed, create file relation failed, err: %s", err.Error())
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

	out = model.ConvertExhibition(&te)
	out.Files, err = e.getFiles(ctx, id)
	if err != nil {
		return nil, err
	}
	out.Organizers, err = e.getOrganizers(ctx, id)
	if err != nil {
		return nil, gerror.Newf("get exhibition failed, get organizers failed, err: %s", err.Error())
	}

	return out, nil
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
		exhibition := model.ConvertExhibition(r)
		exhibition.Files, err = e.getFiles(ctx, exhibition.ID)
		if err != nil {
			return nil, nil, gerror.Newf("list exhibitions failed, get files failed, err: %s", err.Error())
		}

		exhibition.Organizers, err = e.getOrganizers(ctx, exhibition.ID)
		if err != nil {
			return nil, nil, gerror.Newf("list exhibitions failed, get organizers failed, err: %s", err.Error())
		}
		out = append(out, exhibition)
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
// 获取展会主办方
func (e *exhibition) getOrganizers(ctx context.Context, exhibitionID string) (organizers []*model.ExOrganizer, err error) {
	var teo []*entity.TExOrganizer
	err = dao.ExOrganizer.Ctx(ctx).Where(dao.ExOrganizer.Columns().ExhibitionID, exhibitionID).Scan(&teo)
	if err != nil {
		return nil, gerror.Newf("get organizers failed, err: %s", err.Error())
	}
	for _, v := range teo {
		organizers = append(organizers, model.ConvertExOrganizer(v))
	}
	return organizers, nil
}

// 获取展会文件
func (e *exhibition) getFiles(ctx context.Context, exhibitionID string) (files []*model.File, err error) {
	files, err = e.fileDomain.ListByModuleAndCustomID(ctx, model.FileModuleExhibition, exhibitionID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// 建立 文件 - 展会 关联
func (e *exhibition) createFileRelation(ctx context.Context, tx gdb.TX, exhibitionID string, files []*model.File) (err error) {
	for _, v := range files {
		err = e.fileDomain.UpdateCustomInfo(ctx, tx, v.FileID, model.FileModuleExhibition, exhibitionID, v.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

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

	return nil
}
