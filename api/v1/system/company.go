package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type CreateCompanyReq struct {
	g.Meta `path:"/company" method:"post" tags:"公司管理" summary:"创建公司"`
	model.AuthorRequired
	Name        string `p:"name" v:"required#公司名称不能为空" dc:"公司名称"`
	Country     string `p:"country" v:"required#国家不能为空" dc:"国家"`
	Phone       string `p:"phone" v:"required#手机号不能为空" dc:"手机号"`
	Email       string `p:"email" v:"required#邮箱不能为空" dc:"邮箱"`
	Address     string `p:"address" v:"required#地址不能为空" dc:"地址"`
	Description string `p:"description" v:"required#公司描述不能为空" dc:"公司描述"`

	BusinessLicense        string `p:"business_license" v:"required#营业执照不能为空" dc:"营业执照"`
	SocialCreditCode       string `p:"social_credit_code" v:"required#统一社会信用代码不能为空" dc:"统一社会信用代码"`
	LegalPersonName        string `p:"legal_person_name" v:"required#法人姓名不能为空" dc:"法人姓名"`
	LegalPersonCardNumber  string `p:"legal_person_card_number" v:"required#法人证件号不能为空" dc:"法人证件号"`
	LegalPersonPhotoFileID string `p:"legal_person_photo_file_id" v:"required#法人证件照不能为空" dc:"法人证件照"`
	LegalPersonPhone       string `p:"legal_person_phone" v:"required#法人手机号不能为空" dc:"法人手机号"`
}

type CreateCompanyRes struct {
	g.Meta
	ID string `json:"id" dc:"公司ID"`
}

type DeleteCompanyReq struct {
	g.Meta `path:"/company/{id}" method:"delete" tags:"公司管理" summary:"删除公司"`
	model.AuthorRequired
	ID string `p:"id" v:"required#公司ID不能为空"`
}

type DeleteCompanyRes struct {
	g.Meta
}

type UpdateCompanyReq struct {
	g.Meta `path:"/company/{id}" method:"patch" tags:"公司管理" summary:"更新公司"`
	model.AuthorRequired
	ID          string `json:"id" v:"required#公司ID不能为空"`
	Name        string `json:"name" dc:"公司名称"`
	Country     string `json:"country" dc:"国家"`
	Phone       string `json:"phone" dc:"手机号"`
	Email       string `json:"email" dc:"邮箱"`
	Address     string `json:"address" dc:"地址"`
	Description string `json:"description" dc:"公司描述"`

	BusinessLicense       string `json:"business_license" dc:"营业执照"`
	SocialCreditCode      string `json:"social_credit_code" dc:"统一社会信用代码"`
	LegalPersonName       string `json:"legal_person_name" dc:"法人姓名"`
	LegalPersonCardNumber string `json:"legal_person_card_number" dc:"法人证件号"`
	LegalPersonPhotoUrl   string `json:"legal_person_photo_url" dc:"法人证件照"`
	LegalPersonPhone      string `json:"legal_person_phone" dc:"法人手机号"`
}

type UpdateCompanyRes struct {
	g.Meta
}

type GetCompanyReq struct {
	g.Meta `path:"/company/{id}" method:"get" tags:"公司管理" summary:"公司详情"`
	model.AuthorRequired
	ID string `p:"id" v:"required#公司ID不能为空"`
}

type GetCompanyRes struct {
	g.Meta
	Data *Company `json:"data" dc:"公司详情"`
}

type ListCompanyReq struct {
	g.Meta `path:"/company" method:"get" tags:"公司管理" summary:"公司列表/搜索"`
	model.AuthorRequired
	Name string `p:"name" dc:"公司名称(模糊搜索)" v:"max:20#公司名称不能超过20个字符"`
	model.PageReq
}

type ListCompanyRes struct {
	g.Meta
	List    []*Company     `json:"list" dc:"公司列表"`
	PageRes *model.PageRes `json:"page_res" dc:"分页响应"`
}

type ApproveCompanyReq struct {
	g.Meta `path:"/company/{id}/approve" method:"patch" summary:"审核通过" tags:"公司管理"`
	model.AuthorRequired
	ID string `p:"id" v:"required#公司ID不能为空"`
}

type ApproveCompanyRes struct {
	g.Meta
}

type RejectCompanyReq struct {
	g.Meta `path:"/company/{id}/reject" method:"patch" summary:"审核拒绝" tags:"公司管理"`
	model.AuthorRequired
	ID     string `p:"id" v:"required#公司ID不能为空"`
	Reason string `p:"reason" v:"required#审核拒绝原因不能为空" dc:"审核拒绝原因"`
}

type RejectCompanyRes struct {
	g.Meta
}

type BanCompanyReq struct {
	g.Meta `path:"/company/{id}/ban" method:"patch" summary:"禁用公司" tags:"公司管理"`
	model.AuthorRequired
	ID string `p:"id" v:"required#公司ID不能为空"`
}

type BanCompanyRes struct {
	g.Meta
}

type UnbanCompanyReq struct {
	g.Meta `path:"/company/{id}/unban" method:"patch" summary:"解禁公司" tags:"公司管理"`
	model.AuthorRequired
	ID string `p:"id" v:"required#公司ID不能为空"`
}

type UnbanCompanyRes struct {
	g.Meta
}

type ListApplicationsReq struct {
	g.Meta `path:"/company/applications" method:"get" tags:"公司管理" summary:"公司申请列表(待审核)"`
	model.AuthorRequired
	model.PageReq
}

type ListApplicationsRes struct {
	g.Meta
	List    []*Company     `json:"list"`
	PageRes *model.PageRes `json:"page_res"`
}

type Company struct {
	ID                    string `json:"id" dc:"公司ID"`
	Name                  string `json:"name" dc:"公司名称"`
	Country               string `json:"country" dc:"国家"`
	Status                string `json:"status" dc:"公司状态"`
	Phone                 string `json:"phone" dc:"手机号"`
	Email                 string `json:"email" dc:"邮箱"`
	Address               string `json:"address" dc:"地址"`
	Description           string `json:"description" dc:"公司描述"`
	BusinessLicense       string `json:"business_license" dc:"营业执照"`
	SocialCreditCode      string `json:"social_credit_code" dc:"统一社会信用代码"`
	LegalPersonName       string `json:"legal_person_name" dc:"法人姓名"`
	LegalPersonCardNumber string `json:"legal_person_card_number" dc:"法人证件号"`
	LegalPersonPhotoUrl   string `json:"legal_person_photo_url" dc:"法人证件照"`
	LegalPersonPhone      string `json:"legal_person_phone" dc:"法人手机号"`
	ApplyTime             string `json:"apply_time" dc:"申请时间"`
	ApproveTime           string `json:"approve_time" dc:"入驻时间"`
	CreateTime            string `json:"create_time" dc:"创建时间"`
	UpdateTime            string `json:"update_time" dc:"更新时间"`
}
