package interfaces

import (
	"ExhibitionService/internal/model"
	"context"
)

type IFileEngine interface {
	// 获取文件上传连接
	PreUpload(ctx context.Context, in *model.PreUploadReq) (out *model.PreUploadRes, err error)
	// 获取文件下载链接
	PreDownload(ctx context.Context, fileID string) (out *model.PreDownloadRes, err error)
	// 删除文件
	Delete(ctx context.Context, fileID string) error
	// 上报文件上传结果
	ReportUploadResult(ctx context.Context, fileID string, success bool) error
}
