package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

var (
	companyOnce   sync.Once
	companyDomain *company
)

type company struct {
}

func NewCompany() interfaces.ICompany {
	companyOnce.Do(func() {
		companyDomain = &company{}
	})
	return companyDomain
}

func (c *company) CreateCompany(ctx context.Context, tx gdb.TX, in *model.Company) (id string, err error) {
	id = uuid.New().String()
	data := map[string]any{
		dao.Company.Columns().ID:          id,
		dao.Company.Columns().Name:        in.Name,
		dao.Company.Columns().Country:     in.Country,
		dao.Company.Columns().Status:      model.CompanyStatusPending,
		dao.Company.Columns().Phone:       in.Phone,
		dao.Company.Columns().Email:       in.Email,
		dao.Company.Columns().Address:     in.Address,
		dao.Company.Columns().Description: in.Description,

		dao.Company.Columns().BusinessLicense:       in.BusinessLicense,
		dao.Company.Columns().SocialCreditCode:      in.SocialCreditCode,
		dao.Company.Columns().LegalPersonName:       in.LegalPersonName,
		dao.Company.Columns().LegalPersonCardNumber: in.LegalPersonCardNumber,
		dao.Company.Columns().LegalPersonPhotoUrl:   in.LegalPersonPhotoUrl,
		dao.Company.Columns().LegalPersonPhone:      in.LegalPersonPhone,

		dao.Company.Columns().ApplyTime:  time.Now().Unix(),
		dao.Company.Columns().CreateTime: time.Now().Unix(),
		dao.Company.Columns().UpdateTime: time.Now().Unix(),
	}
	_, err = dao.Company.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return "", gerror.New("公司名称已存在")
		}
		return
	}

	return
}

func (c *company) UpdateCompany(ctx context.Context, in *model.Company) (err error) {
	data := map[string]any{}
	if in.Name != "" {
		data[dao.Company.Columns().Name] = in.Name
	}
	if in.Country != "" {
		data[dao.Company.Columns().Country] = in.Country
	}
	if in.Phone != "" {
		data[dao.Company.Columns().Phone] = in.Phone
	}
	if in.Email != "" {
		data[dao.Company.Columns().Email] = in.Email
	}
	if in.Address != "" {
		data[dao.Company.Columns().Address] = in.Address
	}
	if in.Description != "" {
		data[dao.Company.Columns().Description] = in.Description
	}
	if in.BusinessLicense != "" {
		data[dao.Company.Columns().BusinessLicense] = in.BusinessLicense
	}
	if in.SocialCreditCode != "" {
		data[dao.Company.Columns().SocialCreditCode] = in.SocialCreditCode
	}
	if in.LegalPersonName != "" {
		data[dao.Company.Columns().LegalPersonName] = in.LegalPersonName
	}
	if in.LegalPersonCardNumber != "" {
		data[dao.Company.Columns().LegalPersonCardNumber] = in.LegalPersonCardNumber
	}
	if in.LegalPersonPhotoUrl != "" {
		data[dao.Company.Columns().LegalPersonPhotoUrl] = in.LegalPersonPhotoUrl
	}
	if in.LegalPersonPhone != "" {
		data[dao.Company.Columns().LegalPersonPhone] = in.LegalPersonPhone
	}
	if len(data) == 0 {
		return nil
	}

	data[dao.Company.Columns().UpdateTime] = time.Now().Unix()
	_, err = dao.Company.Ctx(ctx).Data(data).Where(dao.Company.Columns().ID, in.ID).Update()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return fmt.Errorf("公司名称已存在")
		}
		return
	}

	return
}

func (c *company) UpdateLegalPersonPhotoURL(ctx context.Context, tx gdb.TX, id string, url string) (err error) {
	_, err = dao.Company.Ctx(ctx).TX(tx).Data(g.Map{
		dao.Company.Columns().LegalPersonPhotoUrl: url,
		dao.Company.Columns().UpdateTime:          time.Now().Unix(),
	}).Where(dao.Company.Columns().ID, id).Update()
	return
}

func (c *company) GetCompany(ctx context.Context, id string) (out *model.Company, err error) {
	var tc entity.TCompany
	err = dao.Company.Ctx(ctx).Where(dao.Company.Columns().ID, id).Scan(&tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.New("公司不存在")
		}
		return
	}

	return model.ConvertCompany(&tc), nil
}

func (c *company) ListCompanies(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Company, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Company.Ctx(ctx)
	if name != "" {
		query = query.WhereLike(dao.Company.Columns().Name, name+"%")
	}
	total, err := query.Count()
	if err != nil {
		return
	}

	var tc []*entity.TCompany
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Company.Columns().CreateTime)
	err = query.Scan(&tc)
	if err != nil {
		return
	}

	for _, r := range tc {
		tmp := model.ConvertCompany(r)
		out = append(out, tmp)
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return
}

func (c *company) DeleteCompany(ctx context.Context, id string) (err error) {
	_, err = dao.Company.Ctx(ctx).Where(dao.Company.Columns().ID, id).Delete()
	return
}

func (c *company) ApproveCompany(ctx context.Context, id string) (err error) {
	_, err = dao.Company.Ctx(ctx).Data(g.Map{
		dao.Company.Columns().Status:      int(model.CompanyStatusApproved),
		dao.Company.Columns().ApproveTime: time.Now().Unix(),
		dao.Company.Columns().UpdateTime:  time.Now().Unix(),
	}).Where(dao.Company.Columns().ID, id).Update()
	return
}

func (c *company) RejectCompany(ctx context.Context, id string) (err error) {
	_, err = dao.Company.Ctx(ctx).Data(g.Map{
		dao.Company.Columns().Status:     int(model.CompanyStatusUnregistered),
		dao.Company.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.Company.Columns().ID, id).Update()
	return
}

func (c *company) BanCompany(ctx context.Context, id string) (err error) {
	_, err = dao.Company.Ctx(ctx).Data(g.Map{
		dao.Company.Columns().Status:     int(model.CompanyStatusDisabled),
		dao.Company.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.Company.Columns().ID, id).Update()
	return
}

func (c *company) UnbanCompany(ctx context.Context, id string) (err error) {
	_, err = dao.Company.Ctx(ctx).Data(g.Map{
		dao.Company.Columns().Status:     int(model.CompanyStatusApproved),
		dao.Company.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.Company.Columns().ID, id).Update()
	return
}

func (c *company) ListApplications(ctx context.Context, pageReq *model.PageReq) (out []*model.Company, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Company.Ctx(ctx).Where(dao.Company.Columns().Status, int(model.CompanyStatusPending))
	total, err := query.Count()
	if err != nil {
		return
	}

	var tc []*entity.TCompany
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.Company.Columns().CreateTime)
	err = query.Scan(&tc)
	if err != nil {
		return
	}

	for _, r := range tc {
		out = append(out, model.ConvertCompany(r))
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return
}
