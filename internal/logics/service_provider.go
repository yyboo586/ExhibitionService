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
	serviceProviderOnce   sync.Once
	serviceProviderDomain *serviceProvider
)

type serviceProvider struct {
	transitionMap map[model.ServiceProviderStatus]map[model.ServiceProviderEvent]ServiceProviderTransition

	fileDomain    interfaces.IFile
	companyDomain interfaces.ICompany
}

func NewServiceProvider(fileDomain interfaces.IFile, companyDomain interfaces.ICompany) interfaces.IServiceProvider {
	serviceProviderOnce.Do(func() {
		serviceProviderDomain = &serviceProvider{
			fileDomain:    fileDomain,
			companyDomain: companyDomain,
		}
		serviceProviderDomain.initTransitionMap()
	})

	return serviceProviderDomain
}

func (sp *serviceProvider) Create(ctx context.Context, tx gdb.TX, in *model.ServiceProvider) (id string, err error) {
	companyID, err := sp.companyDomain.Create(ctx, tx, in.CompanyInfo)
	if err != nil {
		return "", err
	}

	in.CompanyID = companyID
	err = sp.checkFileUploadSuccess(ctx, in.Files)
	if err != nil {
		return "", err
	}

	err = sp.checkFileComplete(ctx, in.Files)
	if err != nil {
		return "", err
	}

	id = uuid.New().String()
	data := map[string]any{
		dao.ServiceProvider.Columns().ID:                  id,
		dao.ServiceProvider.Columns().CompanyID:           in.CompanyID,
		dao.ServiceProvider.Columns().Name:                in.Name,
		dao.ServiceProvider.Columns().Description:         in.Description,
		dao.ServiceProvider.Columns().ContactPersonName:   in.ContactPersonName,
		dao.ServiceProvider.Columns().ContactPersonPhone:  in.ContactPersonPhone,
		dao.ServiceProvider.Columns().ContactPersonEmail:  in.ContactPersonEmail,
		dao.ServiceProvider.Columns().Website:             in.Website,
		dao.ServiceProvider.Columns().Status:              int(model.ServiceProviderStatusPending),
		dao.ServiceProvider.Columns().CreateTime:          time.Now().Unix(),
		dao.ServiceProvider.Columns().SubmitForReviewTime: time.Now().Unix(),
		dao.ServiceProvider.Columns().UpdateTime:          time.Now().Unix(),
	}

	_, err = dao.ServiceProvider.Ctx(ctx).TX(tx).Data(data).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry 'name'") {
			return "", gerror.Newf("create service provider failed, service provider name already exists, name: %s", in.Name)
		}
		return "", gerror.Newf("create service provider failed, err: %s", err.Error())
	}

	err = sp.createFileRelation(ctx, tx, id, in.Files)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (sp *serviceProvider) Get(ctx context.Context, id string) (out *model.ServiceProvider, err error) {
	var tsp entity.TServiceProvider
	err = dao.ServiceProvider.Ctx(ctx).Where(dao.ServiceProvider.Columns().ID, id).Scan(&tsp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get service provider failed, service provider not found, id: %s", id)
		}
		return nil, gerror.Newf("get service provider failed, err: %s", err.Error())
	}

	out = model.ConvertServiceProvider(&tsp)
	out.Files, err = sp.getFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	out.CompanyInfo, err = sp.companyDomain.Get(ctx, out.CompanyID)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (sp *serviceProvider) List(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.ServiceProvider, pageRes *model.PageRes, err error) {
	if pageReq.Page == 0 {
		pageReq.Page = 1
	}
	if pageReq.Size == 0 {
		pageReq.Size = 10
	}

	query := dao.ServiceProvider.Ctx(ctx)
	if name != "" {
		query = query.WhereLike(dao.ServiceProvider.Columns().Name, name+"%")
	}
	total, err := query.Count()
	if err != nil {
		return nil, nil, gerror.Newf("list service providers failed, query count err: %s", err.Error())
	}

	var tsp []*entity.TServiceProvider
	query = query.Page(pageReq.Page, pageReq.Size).OrderDesc(dao.ServiceProvider.Columns().CreateTime)
	err = query.Scan(&tsp)
	if err != nil {
		return nil, nil, gerror.Newf("list service providers failed, query scan err: %s", err.Error())
	}

	for _, r := range tsp {
		tmp := model.ConvertServiceProvider(r)
		tmp.CompanyInfo, err = sp.companyDomain.Get(ctx, tmp.CompanyID)
		if err != nil {
			return nil, nil, gerror.Newf("list service providers failed, get company info failed, err: %s", err.Error())
		}

		out = append(out, tmp)
	}

	pageRes = &model.PageRes{
		Total:       total,
		CurrentPage: pageReq.Page,
	}
	return out, pageRes, nil
}

func (sp *serviceProvider) IsAvailable(ctx context.Context, in *model.ServiceProvider) (err error) {
	if in.Status != model.ServiceProviderStatusApproved {
		return gerror.Newf("service provider is not available, status: %s", model.GetServiceProviderStatusText(in.Status))
	}

	return nil
}

// ---------------私有方法--------------------------------
// 检查文件是否上传成功
func (sp *serviceProvider) checkFileUploadSuccess(ctx context.Context, files []*model.File) (err error) {
	fileIDs := make([]string, len(files))
	for i, v := range files {
		fileIDs[i] = v.FileID
	}
	err = sp.fileDomain.CheckFileUploadSuccess(ctx, fileIDs)
	if err != nil {
		return err
	}
	return nil
}

func (sp *serviceProvider) checkFileComplete(ctx context.Context, files []*model.File) (err error) {
	return nil
}

// 建立 文件 - 服务提供商 关联
func (sp *serviceProvider) createFileRelation(ctx context.Context, tx gdb.TX, serviceProviderID string, files []*model.File) (err error) {
	for _, v := range files {
		err = sp.fileDomain.UpdateCustomInfo(ctx, tx, v.FileID, model.FileModuleServiceProvider, serviceProviderID, v.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

// 清除文件 自定义属性
func (sp *serviceProvider) clearFileRelation(ctx context.Context, tx gdb.TX, serviceProviderID string) (err error) {
	err = sp.fileDomain.ClearCustomInfo(ctx, tx, model.FileModuleServiceProvider, serviceProviderID)
	if err != nil {
		return err
	}
	return nil
}

func (sp *serviceProvider) getFiles(ctx context.Context, serviceProviderID string) (files []*model.File, err error) {
	files, err = sp.fileDomain.ListByModuleAndCustomID(ctx, model.FileModuleServiceProvider, serviceProviderID)
	if err != nil {
		return nil, err
	}

	return files, nil
}
