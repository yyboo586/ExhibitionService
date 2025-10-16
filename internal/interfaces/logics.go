package interfaces

import (
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// 展商事件类型
type MerchantEvent uint8

const (
	MerchantEventApprove MerchantEvent = iota // 审核通过
	MerchantEventReject                       // 审核拒绝
	MerchantEventDisable                      // 禁用展商
	MerchantEventEnable                       // 启用展商
)

type IFile interface {
	// 创建文件
	Create(ctx context.Context, fileID string, fileName string, fileLink string) (err error)
	// 更新文件状态
	UpdateFileStatus(ctx context.Context, fileID string, status model.FileStatus) (err error)
	// 更新文件 自定义属性
	UpdateFileCustomInfo(ctx context.Context, tx gdb.TX, fileID string, module model.FileModule, customID string, typ model.FileType) (err error)
	// 获取文件
	GetFile(ctx context.Context, fileID string) (out *model.File, err error)
	// 按模块与自定义ID获取文件列表
	ListFilesByModuleAndCustomID(ctx context.Context, module model.FileModule, customID string) (out []*model.File, err error)
	// 清除文件 自定义属性
	ClearFileCustomInfo(ctx context.Context, tx gdb.TX, module model.FileModule, customID string) (err error)
	// 检查文件是否上传成功
	IsUploadSuccess(ctx context.Context, fileInfo *model.File) (err error)
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
	GetServiceProvider(ctx context.Context, id string) (out *model.ServiceProvider, err error)
	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.ServiceProvider, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理服务提供商状态事件
	HandleEvent(ctx context.Context, serviceProviderID string, event model.ServiceProviderEvent, data interface{}) (err error)
}

// 展会事件类型
type ExhibitionEvent uint8

const (
	ExhibitionEventSubmitForReview ExhibitionEvent = iota // 提交审核
	ExhibitionEventApprove                                // 审核通过
	ExhibitionEventReject                                 // 审核驳回
	ExhibitionEventStartEnrolling                         // 开始报名(由定时任务触发)
	ExhibitionEventStartRunning                           // 开始进行(由定时任务触发)
	ExhibitionEventEnd                                    // 结束展会(由定时任务触发)
	ExhibitionEventCancel                                 // 取消展会
)

// 展会业务逻辑接口
type IExhibition interface {
	// 创建展会
	Create(ctx context.Context, tx gdb.TX, in *model.Exhibition) (id string, err error)
	// 获取展会详情
	GetExhibition(ctx context.Context, id string) (out *model.Exhibition, err error)
	// 更新展会信息
	UpdateExhibition(ctx context.Context, in *model.Exhibition) (err error)
	// 删除展会
	DeleteExhibition(ctx context.Context, id string) (err error)
	// 列表展会
	ListExhibitions(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error)

	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Exhibition, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理展会状态事件
	HandleEvent(ctx context.Context, exhibitionID string, event ExhibitionEvent, data interface{}) (err error)
}

// 展商业务逻辑接口
type IMerchant interface {
	// 创建展商
	Create(ctx context.Context, tx gdb.TX, in *model.Merchant) (id string, err error)
	// 获取展商详情
	GetMerchant(ctx context.Context, id string) (out *model.Merchant, err error)
	// 更新展商信息
	UpdateMerchant(ctx context.Context, in *model.Merchant) (err error)
	// 删除展商
	DeleteMerchant(ctx context.Context, id string) (err error)
	// 列表展商
	ListMerchants(ctx context.Context, exhibitionID string, name string, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error)
	// 获取待审核列表
	GetPendingList(ctx context.Context, pageReq *model.PageReq) (out []*model.Merchant, pageRes *model.PageRes, err error)

	// 状态流转
	// 处理展商状态事件
	HandleEvent(ctx context.Context, merchantID string, event MerchantEvent, data interface{}) (err error)
}
