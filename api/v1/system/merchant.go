package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type CreateMerchantFileInfo struct {
	FileID   string `json:"file_id" v:"required#文件ID不能为空" dc:"文件ID"`
	FileType string `json:"file_type" v:"required#文件类型不能为空" dc:"文件类型"`
}

type CreateMerchantReq struct {
	g.Meta `path:"/merchants" method:"post" tags:"展商管理" summary:"创建展商"`
	model.AuthorRequired
	CompanyID          string                    `json:"company_id" v:"required#公司ID不能为空" dc:"公司ID"`
	ExhibitionID       string                    `json:"exhibition_id" v:"required#展会ID不能为空" dc:"展会ID"`
	Name               string                    `json:"name" v:"required#展商名称不能为空" dc:"展商名称"`
	Description        string                    `json:"description" dc:"展商描述"`
	BoothNumber        string                    `json:"booth_number" dc:"展位号"`
	ContactPersonName  string                    `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string                    `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string                    `json:"contact_person_email" dc:"联系人邮箱"`
	Files              []*CreateMerchantFileInfo `json:"files" dc:"展商相关文件"`
}

type CreateMerchantRes struct {
	g.Meta
	ID string `json:"id" dc:"展商ID"`
}

type GetMerchantReq struct {
	g.Meta `path:"/merchants/{id}" method:"get" tags:"展商管理" summary:"获取展商详情"`
	ID     string `p:"id" v:"required#展商ID不能为空" dc:"展商ID"`
}

type GetMerchantRes struct {
	g.Meta
	MerchantInfo
	Files []*FileInfo `json:"files" dc:"展商相关文件"`
}

type ListMerchantsReq struct {
	g.Meta       `path:"/merchants" method:"get" tags:"展商管理" summary:"列表展商"`
	ExhibitionID string `p:"exhibition_id" dc:"展会ID"`
	Name         string `p:"name" dc:"展商名称(模糊搜索)"`
	Page         int    `p:"page" d:"1" dc:"页码"`
	Size         int    `p:"size" d:"10" dc:"每页数量"`
}

type ListMerchantsRes struct {
	g.Meta
	Merchants   []*model.Merchant `json:"merchants" dc:"展商列表"`
	Total       int               `json:"total" dc:"总数"`
	CurrentPage int               `json:"current_page" dc:"当前页"`
}

type GetPendingMerchantsReq struct {
	g.Meta `path:"/merchants/pending" method:"get" tags:"展商管理" summary:"获取待审核展商列表"`
	Page   int `p:"page" d:"1" dc:"页码"`
	Size   int `p:"size" d:"10" dc:"每页数量"`
}

type GetPendingMerchantsRes struct {
	g.Meta
	Merchants   []*model.Merchant `json:"merchants" dc:"展商列表"`
	Total       int               `json:"total" dc:"总数"`
	CurrentPage int               `json:"current_page" dc:"当前页"`
}

// 展商状态管理相关接口
type ApproveMerchantReq struct {
	g.Meta `path:"/merchants/{id}/approve" method:"post" tags:"展商管理" summary:"审核通过展商"`
	ID     string `p:"id" v:"required#展商ID不能为空" dc:"展商ID"`
}

type ApproveMerchantRes struct {
	g.Meta
}

type RejectMerchantReq struct {
	g.Meta `path:"/merchants/{id}/reject" method:"post" tags:"展商管理" summary:"审核拒绝展商"`
	ID     string `p:"id" v:"required#展商ID不能为空" dc:"展商ID"`
}

type RejectMerchantRes struct {
	g.Meta
}

type DisableMerchantReq struct {
	g.Meta `path:"/merchants/{id}/disable" method:"post" tags:"展商管理" summary:"禁用展商"`
	ID     string `p:"id" v:"required#展商ID不能为空" dc:"展商ID"`
}

type DisableMerchantRes struct {
	g.Meta
}

type EnableMerchantReq struct {
	g.Meta `path:"/merchants/{id}/enable" method:"post" tags:"展商管理" summary:"启用展商"`
	ID     string `p:"id" v:"required#展商ID不能为空" dc:"展商ID"`
}

type EnableMerchantRes struct {
	g.Meta
}

type MerchantInfo struct {
	ID                 string `json:"id" dc:"展商ID"`
	CompanyID          string `json:"company_id" dc:"公司ID"`
	ExhibitionID       string `json:"exhibition_id" dc:"展会ID"`
	Name               string `json:"name" dc:"展商名称"`
	Status             string `json:"status" dc:"状态"`
	Description        string `json:"description" dc:"展商描述"`
	BoothNumber        string `json:"booth_number" dc:"展位号"`
	ContactPersonName  string `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string `json:"contact_person_email" dc:"联系人邮箱"`
	CreateTime         string `json:"create_time" dc:"创建时间"`
	UpdateTime         string `json:"update_time" dc:"更新时间"`
}
