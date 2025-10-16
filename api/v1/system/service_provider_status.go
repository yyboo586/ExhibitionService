package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPendingSPReq struct {
	g.Meta `path:"/service-provider/pending-list" method:"get" tags:"展会公司" summary:"获取待审核列表"`
	model.PageReq
}

type GetPendingSPRes struct {
	g.Meta  `mime:"application/json"`
	List    []*ServiceProviderInfo `json:"list" dc:"列表"`
	PageRes *model.PageRes         `json:"page_res" dc:"分页响应"`
}

type ApproveServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}/approve" method:"patch" tags:"展会公司" summary:"审核通过"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type ApproveServiceProviderRes struct {
	g.Meta `mime:"application/json"`
}

type RejectServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}/reject" method:"patch" tags:"展会公司" summary:"审核拒绝"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type RejectServiceProviderRes struct {
	g.Meta
}

type DisableServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}/disable" method:"patch" tags:"展会公司" summary:"禁用"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type DisableServiceProviderRes struct {
	g.Meta `mime:"application/json"`
}

type EnableServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}/enable" method:"patch" tags:"展会公司" summary:"启用"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type EnableServiceProviderRes struct {
	g.Meta `mime:"application/json"`
}

type UnregisterServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}/unregister" method:"patch" tags:"展会公司" summary:"注销"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type UnregisterServiceProviderRes struct {
	g.Meta `mime:"application/json"`
}

type RecommitServiceProviderReq struct {
	g.Meta             `path:"/service-provider/{id}/recommit" method:"patch" tags:"展会公司" summary:"重新提交审核"`
	ID                 string       `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
	CompanyInfo        *CompanyInfo `json:"company_info" v:"required#公司信息不能为空" dc:"公司信息"`
	Name               string       `json:"name" v:"required#服务提供商名称不能为空" dc:"服务提供商名称"`
	Website            string       `json:"website" dc:"官网"`
	ContactPersonName  string       `json:"contact_person_name" v:"required#联系人姓名不能为空" dc:"联系人姓名"`
	ContactPersonPhone string       `json:"contact_person_phone" v:"required#联系人电话不能为空" dc:"联系人电话"`
	ContactPersonEmail string       `json:"contact_person_email" v:"required#联系人邮箱不能为空" dc:"联系人邮箱"`
	Description        string       `json:"description" dc:"服务提供商描述"`
	Files              []*FileInfo  `json:"files" dc:"服务提供商相关文件"`
}

type RecommitServiceProviderRes struct {
	g.Meta `mime:"application/json"`
}
