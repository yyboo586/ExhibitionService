package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type CreateMerchantReq struct {
	g.Meta `path:"/merchant" method:"post" tags:"商户管理" summary:"创建"`
	model.AuthorRequired
	Name               string       `json:"name" v:"required#商户名称不能为空" dc:"商户名称"`
	Website            string       `json:"website" dc:"商户官网"`
	ContactPersonName  string       `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string       `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string       `json:"contact_person_email" dc:"联系人邮箱"`
	Description        string       `json:"description" dc:"商户描述"`
	Files              []*FileInfo  `json:"files" dc:"商户相关文件"`
	CompanyInfo        *CompanyInfo `json:"company_info" v:"required#公司信息不能为空" dc:"公司信息"`
}

type CreateMerchantRes struct {
	g.Meta
	ID string `json:"id" dc:"商户ID"`
}

type GetMerchantReq struct {
	g.Meta `path:"/merchant/{id}" method:"get" tags:"商户管理" summary:"详情"`
	ID     string `p:"id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type GetMerchantRes struct {
	g.Meta
	*MerchantInfo
}

type ListMerchantsReq struct {
	g.Meta `path:"/merchant" method:"get" tags:"商户管理" summary:"列表/搜索"`
	Name   string `p:"name" dc:"商户名称(模糊搜索)"`
	model.PageReq
}

type ListMerchantsRes struct {
	g.Meta  `mime:"application/json"`
	List    []*MerchantInfo `json:"list" dc:"列表"`
	PageRes *model.PageRes  `json:"page_res" dc:"分页响应"`
}

type MerchantInfo struct {
	ID                  string `json:"id" dc:"商户ID"`
	Name                string `json:"name" dc:"商户名称"`
	Status              string `json:"status" dc:"状态"`
	Website             string `json:"website" dc:"商户官网"`
	ContactPersonName   string `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone  string `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail  string `json:"contact_person_email" dc:"联系人邮箱"`
	Description         string `json:"description" dc:"商户描述"`
	CreateTime          string `json:"create_time" dc:"创建时间"`
	SubmitForReviewTime string `json:"submit_for_review_time" dc:"提交审核时间"`
	ApproveTime         string `json:"approve_time" dc:"审核通过时间"`
	UpdateTime          string `json:"update_time" dc:"更新时间"`

	Files []*FileInfo `json:"files" dc:"商户相关文件"`

	CompanyInfo *CompanyInfo `json:"company_info" dc:"公司信息"`
}
