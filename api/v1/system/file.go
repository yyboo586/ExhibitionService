package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type UploadFileReq struct {
	g.Meta `path:"/files" method:"post" tags:"文件管理" summary:"获取文件上传链接"`
	model.AuthorRequired
	FileName    string `json:"file_name" v:"required#文件名不能为空" dc:"文件名"`
	ContentType string `json:"content_type" v:"required#文件类型不能为空" dc:"文件类型"`
	FileSize    int64  `json:"file_size" v:"required#文件大小不能为空" dc:"文件大小"`
}

type UploadFileRes struct {
	g.Meta
	FileID       string `json:"file_id" dc:"文件ID"`
	OriginalName string `json:"original_name" dc:"文件名称"`
	UploadURL    string `json:"upload_url" dc:"上传URL"`
}

type ReportFileUploadResultReq struct {
	g.Meta `path:"/files/{file_id}/report-upload-result" method:"patch" tags:"文件管理" summary:"上报文件上传结果"`
	model.AuthorRequired
	FileID  string `p:"file_id" v:"required#文件ID不能为空"`
	Success bool   `json:"success" v:"required#上传结果不能为空" dc:"上传结果"`
}

type ReportFileUploadResultRes struct {
	g.Meta
}

type FileInfo struct {
	FileID     string `json:"file_id" dc:"文件ID"`
	FileType   string `json:"file_type" dc:"文件类型"`
	FileName   string `json:"file_name" dc:"文件名称"`
	FileLink   string `json:"file_link" dc:"文件链接"`
	CreateTime string `json:"create_time" dc:"创建时间"`
	UpdateTime string `json:"update_time" dc:"更新时间"`
}
