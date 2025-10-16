package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

type FileModule int

const (
	_                         FileModule = iota
	FileModuleCompany                    // 公司
	FileModuleOrganizer                  // 主办方
	FileModuleExhibition                 // 展会
	FileModuleMerchant                   // 展商
	FileModuleServiceProvider            // 服务提供商
)

func GetFileModuleText(module FileModule) string {
	switch module {
	case FileModuleCompany:
		return "公司"
	case FileModuleOrganizer:
		return "主办方"
	case FileModuleExhibition:
		return "展会"
	case FileModuleMerchant:
		return "展商"
	case FileModuleServiceProvider:
		return "服务提供商"
	default:
		return "未知模块"
	}
}

type FileType int

const (
	FileTypeUnknown                 FileType = iota
	FileTypeCompanyLicense                   // Company License(公司营业执照)
	FileTypeCompanyLegalPersonPhoto          // Company Legal Person Photo(公司法人证件照)
	FileTypeLogo                             // Logo
	FileTypeBanner                           // Banner(横幅图)
	FileTypeProductImage                     // Product Image(展品图)
	FileTypePoster                           // Poster(海报)
	FileTypeDocument                         // Document(文档)
	FileTypeVideo                            // Video(视频)
)

func GetFileType(text string) FileType {
	switch text {
	case "Company License":
		return FileTypeCompanyLicense
	case "Company Legal Person Photo":
		return FileTypeCompanyLegalPersonPhoto
	case "Logo":
		return FileTypeLogo
	case "Banner":
		return FileTypeBanner
	case "Product Image":
		return FileTypeProductImage
	case "Poster":
		return FileTypePoster
	case "Document":
		return FileTypeDocument
	case "Video":
		return FileTypeVideo
	default:
		return FileTypeUnknown
	}
}

func GetFileTypeText(typ FileType) string {
	switch typ {
	case FileTypeCompanyLicense:
		return "Company License"
	case FileTypeCompanyLegalPersonPhoto:
		return "Company Legal Person Photo"
	case FileTypeLogo:
		return "Logo"
	case FileTypeBanner:
		return "Banner"
	case FileTypeProductImage:
		return "Product Image"
	case FileTypePoster:
		return "Poster"
	case FileTypeDocument:
		return "Document"
	case FileTypeVideo:
		return "Video"
	default:
		return "Unknown File Type"
	}
}

type FileStatus int

const (
	FileStatusInit          FileStatus = iota // Init
	FileStatusUploadSuccess                   // Upload Success
	FileStatusUploadFailed                    // Upload Failed
)

func GetFileStatusText(status FileStatus) string {
	switch status {
	case FileStatusInit:
		return "Init"
	case FileStatusUploadSuccess:
		return "Upload Success"
	case FileStatusUploadFailed:
		return "Upload Failed"
	default:
		return "Unknown File Status"
	}
}

type File struct {
	ID         string     `json:"id"`
	Module     FileModule `json:"module"`
	CustomID   string     `json:"custom_id"`
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
		Module:     FileModule(t.Module),
		CustomID:   t.CustomID,
		Type:       FileType(t.Type),
		FileID:     t.FileID,
		FileName:   t.FileName,
		FileLink:   t.FileLink,
		Status:     FileStatus(t.Status),
		CreateTime: time.Unix(t.CreateTime, 0),
		UpdateTime: time.Unix(t.UpdateTime, 0),
	}
}

func GetFileStatusFromBool(success bool) FileStatus {
	if success {
		return FileStatusUploadSuccess
	}
	return FileStatusUploadFailed
}
