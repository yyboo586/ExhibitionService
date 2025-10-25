package interfaces

import (
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	asyncTask "github.com/yyboo586/common/AsyncTask"
)

type IFile interface {
	// 创建文件
	Create(ctx context.Context, fileID string, fileName string, fileLink string) (err error)
	// 更新文件状态
	UpdateStatus(ctx context.Context, fileID string, status model.FileStatus) (err error)
	// 更新文件 自定义属性
	UpdateCustomInfo(ctx context.Context, tx gdb.TX, fileID string, module model.FileModule, customID string, typ model.FileType) (err error)
	// 获取文件
	Get(ctx context.Context, fileID string) (out *model.File, err error)
	// 按模块与自定义ID获取文件列表
	ListByModuleAndCustomID(ctx context.Context, module model.FileModule, customID string) (out []*model.File, err error)
	// 清除文件 自定义属性
	ClearCustomInfo(ctx context.Context, tx gdb.TX, module model.FileModule, customID string) (err error)
	// 检查文件是否上传成功
	IsUploadSuccess(ctx context.Context, fileInfo *model.File) (err error)
	// 检查文件是否上传成功
	CheckFileUploadSuccess(ctx context.Context, fileIDs []string) (err error)
}

type ICompany interface {
	// 创建展会公司
	Create(ctx context.Context, tx gdb.TX, in *model.Company) (id string, err error)
	// 获取公司详情
	Get(ctx context.Context, id string) (out *model.Company, err error)
	// 重新提交审核接口
	Recommit(ctx context.Context, tx gdb.TX, typ model.CompanyType, in *model.Company) (err error)
}

// 服务提供商业务逻辑接口
type IServiceProvider interface {
	// 创建服务提供商
	Create(ctx context.Context, tx gdb.TX, in *model.ServiceProvider) (id string, err error)
	// 获取服务提供商详情
	Get(ctx context.Context, id string) (out *model.ServiceProvider, err error)
	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.ServiceProvider, pageRes *model.PageRes, err error)
	// 列表/搜索
	List(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.ServiceProvider, pageRes *model.PageRes, err error)

	// 检查服务提供商是否可用
	IsAvailable(ctx context.Context, in *model.ServiceProvider) (err error)

	// 状态流转
	// 处理服务提供商状态事件
	HandleEvent(ctx context.Context, serviceProviderID string, event model.ServiceProviderEvent, data interface{}) (err error)
}

// 展商业务逻辑接口
type IMerchant interface {
	// 创建展商
	Create(ctx context.Context, tx gdb.TX, in *model.Merchant) (id string, err error)
	// 获取展商详情
	Get(ctx context.Context, id string) (out *model.Merchant, err error)
	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error)
	// 列表/搜索
	List(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理展商状态事件
	HandleEvent(ctx context.Context, merchantID string, event model.MerchantEvent, data interface{}) (err error)
}

// 展会业务逻辑接口
type IExhibition interface {
	// 创建展会
	Create(ctx context.Context, tx gdb.TX, in *model.Exhibition) (id string, err error)
	// 获取展会详情
	GetExhibition(ctx context.Context, id string) (out *model.Exhibition, err error)
	// 列表展会
	ListExhibitions(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error)

	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理展会状态事件
	HandleEvent(ctx context.Context, exhibitionID string, event model.ExhibitionEvent, data interface{}) (err error)

	// 异步任务
	// 处理展会自动开始报名任务
	HandleTaskAutoStartEnrolling(ctx context.Context, task *asyncTask.Task) (err error)
	// 处理展会自动结束报名任务
	HandleTaskAutoEndEnrolling(ctx context.Context, task *asyncTask.Task) (err error)
	// 处理展会自动开始进行任务
	HandleTaskAutoStartRunning(ctx context.Context, task *asyncTask.Task) (err error)
	// 处理展会自动结束任务
	HandleTaskAutoEnd(ctx context.Context, task *asyncTask.Task) (err error)

	// 商户申请参加展会
	CreateApplyForExhibition(ctx context.Context, tx gdb.TX, exhibitionID string, merchantID string) (err error)
	// 获取商户在展会中的申请状态
	GetExMerchantApplication(ctx context.Context, exhibitionID string, merchantID string) (out *model.ExhibitionMerchant, err error)
	// 获取展会的商户申请列表
	ListExhibitionApplications(ctx context.Context, exhibitionID string, pageReq *model.PageReq) (out []*model.ExhibitionMerchant, pageRes *model.PageRes, err error)
	// 获取商户的展会申请列表
	ListMerchantApplications(ctx context.Context, merchantID string, pageReq *model.PageReq) (out []*model.ExhibitionMerchant, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理展会与商户关联状态事件
	HandleExMerchantEvent(ctx context.Context, exhibitionID string, merchantID string, event model.ExMerchantEvent, data interface{}) (err error)
}
