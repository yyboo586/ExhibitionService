package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPendingMerchantsReq struct {
	g.Meta `path:"/merchant/pending-list" method:"get" tags:"商户管理/状态管理" summary:"获取待审核列表"`
	model.PageReq
}

type GetPendingMerchantsRes struct {
	g.Meta  `mime:"application/json"`
	List    []*MerchantInfo `json:"list" dc:"列表"`
	PageRes *model.PageRes  `json:"page_res" dc:"分页响应"`
}

type ApproveMerchantReq struct {
	g.Meta `path:"/merchant/{id}/approve" method:"patch" tags:"商户管理/状态管理" summary:"审核通过"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type ApproveMerchantRes struct {
	g.Meta `mime:"application/json"`
}

type RejectMerchantReq struct {
	g.Meta `path:"/merchant/{id}/reject" method:"patch" tags:"商户管理/状态管理" summary:"审核拒绝"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type RejectMerchantRes struct {
	g.Meta `mime:"application/json"`
}

type DisableMerchantReq struct {
	g.Meta `path:"/merchant/{id}/disable" method:"patch" tags:"商户管理/状态管理" summary:"禁用"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type DisableMerchantRes struct {
	g.Meta `mime:"application/json"`
}

type EnableMerchantReq struct {
	g.Meta `path:"/merchant/{id}/enable" method:"patch" tags:"商户管理/状态管理" summary:"启用"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type EnableMerchantRes struct {
	g.Meta `mime:"application/json"`
}

type UnregisterMerchantReq struct {
	g.Meta `path:"/merchant/{id}/unregister" method:"patch" tags:"商户管理/状态管理" summary:"注销"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type UnregisterMerchantRes struct {
	g.Meta `mime:"application/json"`
}

type RecommitMerchantReq struct {
	g.Meta             `path:"/merchant/{id}/recommit" method:"patch" tags:"商户管理/状态管理" summary:"重新提交审核"`
	ID                 string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
	Name               string `json:"name" v:"required#商户名称不能为空" dc:"商户名称"`
	Website            string `json:"website" dc:"商户官网"`
	ContactPersonName  string `json:"contact_person_name" v:"required#联系人姓名不能为空" dc:"联系人姓名"`
	ContactPersonPhone string `json:"contact_person_phone" v:"required#联系人电话不能为空" dc:"联系人电话"`
	ContactPersonEmail string `json:"contact_person_email" v:"required#联系人邮箱不能为空" dc:"联系人邮箱"`
	Description        string `json:"description" dc:"商户描述"`

	Files       []*FileInfo  `json:"files" dc:"商户相关文件"`
	CompanyInfo *CompanyInfo `json:"company_info" v:"required#公司信息不能为空" dc:"公司信息"`
}

type RecommitMerchantRes struct {
	g.Meta `mime:"application/json"`
}
