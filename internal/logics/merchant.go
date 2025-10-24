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

type merchant struct {
	fileDomain    interfaces.IFile
	companyDomain interfaces.ICompany

	transitionMap map[model.MerchantStatus]map[model.MerchantEvent]MerchantTransition
}

func NewMerchant(fileDomain interfaces.IFile, companyDomain interfaces.ICompany) interfaces.IMerchant {
	merchantOnce.Do(func() {
		merchantDomain = &merchant{
			fileDomain:    fileDomain,
			companyDomain: companyDomain,
		}
		merchantDomain.initTransitionMap()
	})
	return merchantDomain
}

func (m *merchant) Create(ctx context.Context, tx gdb.TX, in *model.Merchant) (id string, err error) {
	companyID, err := m.companyDomain.Create(ctx, tx, in.CompanyInfo)
	if err != nil {
		return "", err
	}

	in.CompanyID = companyID
	err = m.checkFileUploadSuccess(ctx, in.Files)
	if err != nil {
		return "", err
	}

	err = m.checkFileComplete(ctx, in.Files)
	if err != nil {
		return "", err
	}

	id = uuid.New().String()
	data := map[string]any{
		dao.Merchant.Columns().ID:                  id,
		dao.Merchant.Columns().CompanyID:           in.CompanyID,
		dao.Merchant.Columns().Name:                in.Name,
		dao.Merchant.Columns().Description:         in.Description,
		dao.Merchant.Columns().ContactPersonName:   in.ContactPersonName,
		dao.Merchant.Columns().ContactPersonPhone:  in.ContactPersonPhone,
		dao.Merchant.Columns().ContactPersonEmail:  in.ContactPersonEmail,
		dao.Merchant.Columns().Website:             in.Website,
		dao.Merchant.Columns().Status:              int(model.MerchantStatusPending),
		dao.Merchant.Columns().CreateTime:          time.Now().Unix(),
		dao.Merchant.Columns().SubmitForReviewTime: time.Now().Unix(),
		dao.Merchant.Columns().UpdateTime:          time.Now().Unix(),
	}

	_, err = dao.Merchant.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return "", gerror.Newf("create merchant failed, merchant name already exists, name: %s", in.Name)
		}
		return "", gerror.Newf("create merchant failed, err: %s", err.Error())
	}

	err = m.createFileRelation(ctx, tx, id, in.Files)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *merchant) Get(ctx context.Context, id string) (out *model.Merchant, err error) {
	var tm entity.TMerchant
	err = dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().ID, id).Scan(&tm)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get merchant failed, merchant not found, id: %s", id)
		}
		return nil, gerror.Newf("get merchant failed, err: %s", err.Error())
	}

	out = model.ConvertMerchant(&tm)
	out.Files, err = m.getFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	out.CompanyInfo, err = m.companyDomain.Get(ctx, out.CompanyID)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (m *merchant) List(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.Merchant.Ctx(ctx)
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
		tmp := model.ConvertMerchant(r)
		tmp.Files, err = m.getFiles(ctx, tmp.ID)
		if err != nil {
			return nil, nil, gerror.Newf("list merchants failed, get files failed, err: %s", err.Error())
		}
		tmp.CompanyInfo, err = m.companyDomain.Get(ctx, tmp.CompanyID)
		if err != nil {
			return nil, nil, gerror.Newf("list merchants failed, get company info failed, err: %s", err.Error())
		}

		out = append(out, tmp)
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

// ---------------私有方法--------------------------------
// 检查文件是否上传成功
func (m *merchant) checkFileUploadSuccess(ctx context.Context, files []*model.File) (err error) {
	for _, v := range files {
		fileInfo, err := m.fileDomain.Get(ctx, v.FileID)
		if err != nil {
			return err
		}
		err = m.fileDomain.IsUploadSuccess(ctx, fileInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *merchant) checkFileComplete(ctx context.Context, files []*model.File) (err error) {
	return nil
}

// 建立 文件 - 商户 关联
func (m *merchant) createFileRelation(ctx context.Context, tx gdb.TX, merchantID string, files []*model.File) (err error) {
	for _, v := range files {
		err = m.fileDomain.UpdateCustomInfo(ctx, tx, v.FileID, model.FileModuleMerchant, merchantID, v.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

// 清除文件 自定义属性
func (m *merchant) clearFileRelation(ctx context.Context, tx gdb.TX, merchantID string) (err error) {
	err = m.fileDomain.ClearCustomInfo(ctx, tx, model.FileModuleMerchant, merchantID)
	if err != nil {
		return err
	}
	return nil
}

func (m *merchant) getFiles(ctx context.Context, merchantID string) (files []*model.File, err error) {
	files, err = m.fileDomain.ListByModuleAndCustomID(ctx, model.FileModuleMerchant, merchantID)
	if err != nil {
		return nil, err
	}

	return files, nil
}
