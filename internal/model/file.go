package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

type FileType int

const (
	FileTypeUnknown          FileType = iota
	FileTypeLegalPersonPhoto          // 法人证件照
)

func GetFileType(text string) FileType {
	switch text {
	case "法人证件照":
		return FileTypeLegalPersonPhoto
	default:
		return FileTypeUnknown
	}
}

type FileStatus int

const (
	FileStatusInit          FileStatus = iota // 初始状态
	FileStatusUploadSuccess                   // 上传成功
	FileStatusUploadFailed                    // 上传失败
)

func GetFileStatusFromBool(success bool) FileStatus {
	if success {
		return FileStatusUploadSuccess
	}
	return FileStatusUploadFailed
}

func GetFileStatusText(status FileStatus) string {
	switch status {
	case FileStatusInit:
		return "初始状态"
	case FileStatusUploadSuccess:
		return "上传成功"
	case FileStatusUploadFailed:
		return "上传失败"
	default:
		return "未知状态"
	}
}

type File struct {
	ID         string     `json:"id"`
	CompanyID  string     `json:"company_id"`
	Type       FileType   `json:"type"`
	FileID     string     `json:"file_id"`
	FileName   string     `json:"file_name"`
	FileLink   string     `json:"file_link"`
	Status     FileStatus `json:"status"`
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time  `json:"update_time"`
}

func ConvertFile(t *entity.TFile) *File {
	return &File{
		ID:         t.ID,
		CompanyID:  t.CompanyID,
		Type:       FileType(t.Type),
		FileID:     t.FileID,
		FileName:   t.FileName,
		FileLink:   t.FileLink,
		Status:     FileStatus(t.Status),
		CreateTime: time.Unix(t.CreateTime, 0),
		UpdateTime: time.Unix(t.UpdateTime, 0),
	}
}
