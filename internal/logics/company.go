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
	companyOnce   sync.Once
	companyDomain *company
)

type company struct {
	fileDomain interfaces.IFile
}

func NewCompany(fileDomain interfaces.IFile) interfaces.ICompany {
	companyOnce.Do(func() {
		companyDomain = &company{
			fileDomain: fileDomain,
		}
	})
	return companyDomain
}

func (c *company) Create(ctx context.Context, tx gdb.TX, in *model.Company) (id string, err error) {
	// 检查文件是否上传成功
	err = c.checkFileUploadSuccess(ctx, in.Files)
	if err != nil {
		return "", err
	}

	// 检查公司资质文件是否完整
	err = c.checkFileComplete(ctx, in.Files)
	if err != nil {
		return "", err
	}

	id = uuid.New().String()
	data := map[string]any{
		dao.Company.Columns().ID:          id,
		dao.Company.Columns().Name:        in.Name,
		dao.Company.Columns().Type:        int(in.Type),
		dao.Company.Columns().Country:     in.Country,
		dao.Company.Columns().City:        in.City,
		dao.Company.Columns().Address:     in.Address,
		dao.Company.Columns().Email:       in.Email,
		dao.Company.Columns().Description: in.Description,

		dao.Company.Columns().SocialCreditCode:      in.SocialCreditCode,
		dao.Company.Columns().LegalPersonName:       in.LegalPersonName,
		dao.Company.Columns().LegalPersonCardNumber: in.LegalPersonCardNumber,

		dao.Company.Columns().CreateTime: time.Now().Unix(),
		dao.Company.Columns().UpdateTime: time.Now().Unix(),
	}
	_, err = dao.Company.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry ") {
			return "", gerror.Newf("create company failed, company social credit code already exists, social credit code: %s", in.SocialCreditCode)
		}
		return "", gerror.Newf("create company failed, err: %s", err.Error())
	}

	err = c.createFileRelation(ctx, tx, id, in.Files)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *company) Get(ctx context.Context, id string) (out *model.Company, err error) {
	var tc entity.TCompany
	err = dao.Company.Ctx(ctx).Where(dao.Company.Columns().ID, id).Scan(&tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get company failed, company not found, id: %s", id)
		}
		return nil, gerror.Newf("get company failed, err: %s", err.Error())
	}

	out = model.ConvertCompany(&tc)
	out.Files, err = c.getFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	out.BusinessLicense = c.getFileLink(ctx, out.Files, model.FileTypeCompanyLicense)
	out.LegalPersonPhoto = c.getFileLink(ctx, out.Files, model.FileTypeCompanyLegalPersonPhoto)

	return out, nil
}

func (c *company) Recommit(ctx context.Context, tx gdb.TX, typ model.CompanyType, in *model.Company) (err error) {
	switch typ {
	case model.CompanyTypeServiceProvider:
		return c.recommitServiceProvider(ctx, tx, in)
	case model.CompanyTypeMerchant:
		return c.recommitMerchant(ctx, tx, in)
	default:
		return gerror.Newf("recommit company failed, company type not supported, type: %s", model.GetCompanyTypeText(typ))
	}
}

// ---------------私有方法--------------------------------

func (c *company) recommitServiceProvider(ctx context.Context, tx gdb.TX, in *model.Company) (err error) {
	oldCompany, err := c.Get(ctx, in.ID)
	if err != nil {
		return err
	}
	var dataUpdate map[string]any = make(map[string]any)
	if oldCompany.Name != in.Name {
		dataUpdate[dao.Company.Columns().Name] = in.Name
	}
	if oldCompany.Country != in.Country {
		dataUpdate[dao.Company.Columns().Country] = in.Country
	}
	if oldCompany.City != in.City {
		dataUpdate[dao.Company.Columns().City] = in.City
	}
	if oldCompany.Address != in.Address {
		dataUpdate[dao.Company.Columns().Address] = in.Address
	}
	if oldCompany.Email != in.Email {
		dataUpdate[dao.Company.Columns().Email] = in.Email
	}
	if oldCompany.Description != in.Description {
		dataUpdate[dao.Company.Columns().Description] = in.Description
	}
	if oldCompany.SocialCreditCode != in.SocialCreditCode {
		dataUpdate[dao.Company.Columns().SocialCreditCode] = in.SocialCreditCode
	}
	if oldCompany.LegalPersonName != in.LegalPersonName {
		dataUpdate[dao.Company.Columns().LegalPersonName] = in.LegalPersonName
	}
	if oldCompany.LegalPersonCardNumber != in.LegalPersonCardNumber {
		dataUpdate[dao.Company.Columns().LegalPersonCardNumber] = in.LegalPersonCardNumber
	}

	// 如果数据有更新，则更新数据
	if len(dataUpdate) > 0 {
		dataUpdate[dao.Company.Columns().UpdateTime] = time.Now().Unix()
		_, err = dao.Company.Ctx(ctx).Data(dataUpdate).Where(dao.Company.Columns().ID, in.ID).Update()
		if err != nil {
			return gerror.Newf("recommit service provider failed, err: %s", err.Error())
		}
	}

	// 检查文件是否上传成功
	err = c.checkFileUploadSuccess(ctx, in.Files)
	if err != nil {
		return err
	}

	// 检查公司资质文件是否完整
	err = c.checkFileComplete(ctx, in.Files)
	if err != nil {
		return err
	}

	// 清除旧文件关联
	err = c.clearFileRelation(ctx, tx, in.ID)
	if err != nil {
		return err
	}

	// 创建新文件关联
	err = c.createFileRelation(ctx, tx, in.ID, in.Files)
	if err != nil {
		return err
	}

	return nil
}

func (c *company) recommitMerchant(ctx context.Context, tx gdb.TX, in *model.Company) (err error) {
	return nil
}

// 检查文件是否上传成功
func (c *company) checkFileUploadSuccess(ctx context.Context, files []*model.File) (err error) {
	for _, v := range files {
		fileInfo, err := c.fileDomain.GetFile(ctx, v.FileID)
		if err != nil {
			return err
		}
		err = c.fileDomain.IsUploadSuccess(ctx, fileInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// 检查公司资质文件是否完整
func (c *company) checkFileComplete(ctx context.Context, files []*model.File) (err error) {
	fileMap := make(map[model.FileType]*model.File)
	for _, v := range files {
		fileMap[v.Type] = v
	}

	requiresTypes := []model.FileType{
		model.FileTypeCompanyLicense,
		model.FileTypeCompanyLegalPersonPhoto,
	}

	for _, v := range requiresTypes {
		_, ok := fileMap[v]
		if !ok {
			return gerror.Newf("check company file not complete, file type: %s", model.GetFileTypeText(v))
		}
	}
	return nil
}

// 建立 文件 - 公司 关联
func (c *company) createFileRelation(ctx context.Context, tx gdb.TX, companyID string, files []*model.File) (err error) {
	for _, v := range files {
		err = c.fileDomain.UpdateFileCustomInfo(ctx, tx, v.FileID, model.FileModuleCompany, companyID, v.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

// 清除文件 自定义属性
func (c *company) clearFileRelation(ctx context.Context, tx gdb.TX, companyID string) (err error) {
	err = c.fileDomain.ClearFileCustomInfo(ctx, tx, model.FileModuleCompany, companyID)
	if err != nil {
		return err
	}
	return nil
}

func (c *company) getFiles(ctx context.Context, companyID string) (files []*model.File, err error) {
	files, err = c.fileDomain.ListFilesByModuleAndCustomID(ctx, model.FileModuleCompany, companyID)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (c *company) getFileLink(ctx context.Context, files []*model.File, fileType model.FileType) (link string) {
	for _, v := range files {
		if v.Type == fileType {
			return v.FileLink
		}
	}
	return ""
}
