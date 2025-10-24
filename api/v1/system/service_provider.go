package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

/*
- ServiceProvider模块不提供删除接口。
*/

type CreateServiceProviderReq struct {
	g.Meta `path:"/service-provider" method:"post" tags:"展会公司" summary:"创建"`
	model.AuthorRequired
	CompanyInfo        *CompanyInfo `json:"company_info" v:"required#公司信息不能为空" dc:"公司信息"`
	Name               string       `json:"name" v:"required#服务提供商名称不能为空" dc:"服务提供商名称"`
	Website            string       `json:"website" dc:"官网"`
	ContactPersonName  string       `json:"contact_person_name" v:"required#联系人姓名不能为空" dc:"联系人姓名"`
	ContactPersonPhone string       `json:"contact_person_phone" v:"required#联系人电话不能为空" dc:"联系人电话"`
	ContactPersonEmail string       `json:"contact_person_email" v:"required#联系人邮箱不能为空" dc:"联系人邮箱"`
	Description        string       `json:"description" dc:"服务提供商描述"`
	Files              []*FileInfo  `json:"files" dc:"服务提供商相关文件"`
}

type CreateServiceProviderRes struct {
	g.Meta `mime:"application/json"`
	ID     string `json:"id" dc:"服务提供商ID"`
}

type GetServiceProviderReq struct {
	g.Meta `path:"/service-provider/{id}" method:"get" tags:"展会公司" summary:"详情"`
	ID     string `p:"id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
}

type GetServiceProviderRes struct {
	g.Meta               `mime:"application/json"`
	*ServiceProviderInfo `json:"service_provider_info" dc:"服务提供商信息"`
}

type ListServiceProvidersReq struct {
	g.Meta `path:"/service-provider" method:"get" tags:"展会公司" summary:"列表/搜索"`
	Name   string `p:"name" dc:"服务提供商名称(模糊搜索)"`
	model.PageReq
}

type ListServiceProvidersRes struct {
	g.Meta  `mime:"application/json"`
	List    []*ServiceProviderInfo `json:"list" dc:"列表"`
	PageRes *model.PageRes         `json:"page_res" dc:"分页响应"`
}

type ServiceProviderInfo struct {
	ID                 string `json:"id" dc:"服务提供商ID"`
	CompanyID          string `json:"company_id" dc:"公司ID"`
	Name               string `json:"name" dc:"服务提供商名称"`
	Status             string `json:"status" dc:"状态"`
	Website            string `json:"website" dc:"官网"`
	ContactPersonName  string `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string `json:"contact_person_email" dc:"联系人邮箱"`
	Description        string `json:"description" dc:"服务提供商描述"`

	CreateTime          string `json:"create_time" dc:"创建时间"`
	SubmitForReviewTime string `json:"submit_for_review_time" dc:"提交审核时间"`
	ApproveTime         string `json:"approve_time" dc:"审核通过时间"`
	UpdateTime          string `json:"update_time" dc:"更新时间"`

	Files []*FileInfo `json:"files" dc:"服务提供商相关文件"`

	CompanyInfo *CompanyInfo `json:"company_info" dc:"公司信息"`
}
