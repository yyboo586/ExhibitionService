package interfaces

import (
	"ExhibitionService/internal/model"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

type HTTPClient interface {
	GET(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error)
	POST(ctx context.Context, url string, header map[string]interface{}, body interface{}) (status int, respBody []byte, err error)
	DELETE(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error)
}

type IFile interface {
	// 创建文件
	Create(ctx context.Context, fileID string, fileName string, fileLink string) (err error)
	// 更新文件状态
	UpdateFileStatus(ctx context.Context, fileID string, status model.FileStatus) (err error)
	// 更新文件 公司信息
	UpdateFileCompanyInfo(ctx context.Context, tx gdb.TX, fileID string, companyID string, typ model.FileType) (err error)
	// 获取文件
	GetFile(ctx context.Context, fileID string) (out *model.File, err error)
}

type ICompany interface {
	// 创建展会公司
	CreateCompany(ctx context.Context, tx gdb.TX, in *model.Company) (id string, err error)
	// 删除公司
	DeleteCompany(ctx context.Context, id string) (err error)
	// 更新公司信息
	UpdateCompany(ctx context.Context, in *model.Company) (err error)
	// 更新法人照片URL
	UpdateLegalPersonPhotoURL(ctx context.Context, tx gdb.TX, id string, url string) (err error)
	// 获取公司详情
	GetCompany(ctx context.Context, id string) (out *model.Company, err error)
	// 列表/搜索公司
	ListCompanies(ctx context.Context, name string, pageReq *model.PageReq) (out []*model.Company, pageRes *model.PageRes, err error)

	// 状态流转
	// 审核通过
	ApproveCompany(ctx context.Context, id string) (err error)
	// 审核拒绝
	RejectCompany(ctx context.Context, id string) (err error)
	// 禁用公司
	BanCompany(ctx context.Context, id string) (err error)
	// 解禁公司
	UnbanCompany(ctx context.Context, id string) (err error)
	// 获取待审核申请列表
	ListApplications(ctx context.Context, pageReq *model.PageReq) (out []*model.Company, pageRes *model.PageRes, err error)
}
