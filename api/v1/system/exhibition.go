package system

import (
	"ExhibitionService/internal/model"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type OrganizerInfo struct {
	CompanyID          string `json:"company_id" v:"required#公司ID不能为空" dc:"公司ID"`
	Name               string `json:"name" v:"required#主办方名称不能为空" dc:"主办方名称"`
	RoleType           string `json:"role_type" v:"required#主办方角色类型不能为空" dc:"主办方角色类型(主办方、联合主办方、协办方)"`
	Description        string `json:"description" dc:"主办方描述"`
	ContactPersonName  string `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string `json:"contact_person_email" dc:"联系人邮箱"`
	Website            string `json:"website" dc:"主办方官网"`
}

type CreateExhibitionFileInfo struct {
	FileID   string `json:"file_id" v:"required#文件ID不能为空" dc:"文件ID"`
	FileType string `json:"file_type" v:"required#文件类型不能为空" dc:"文件类型(Company License、Company Legal Person Photo、Logo、Banner、Product Image、Poster、Document、Video)"`
}

type CreateExhibitionReq struct {
	g.Meta `path:"/exhibitions" method:"post" tags:"展会管理" summary:"创建展会"`
	model.AuthorRequired
	ServiceProviderID string                      `json:"service_provider_id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
	Title             string                      `json:"title" v:"required#展会名称不能为空" dc:"展会名称"`
	Industry          string                      `json:"industry" v:"required#所属行业不能为空" dc:"所属行业"`
	Tags              string                      `json:"tags" v:"required#展会标签不能为空" dc:"展会标签"`
	Website           string                      `json:"website" dc:"展会官网"`
	Venue             string                      `json:"venue" v:"required#展会地点不能为空" dc:"展会地点"`
	VenueAddress      string                      `json:"venue_address" v:"required#展会详细地址不能为空" dc:"展会详细地址"`
	Country           string                      `json:"country" v:"required#国家不能为空" dc:"国家"`
	City              string                      `json:"city" v:"required#城市不能为空" dc:"城市"`
	Description       string                      `json:"description" v:"required#展会描述不能为空" dc:"展会描述"`
	RegistrationStart time.Time                   `json:"registration_start" v:"required#报名开始时间不能为空" dc:"报名开始时间"`
	RegistrationEnd   time.Time                   `json:"registration_end" v:"required#报名结束时间不能为空" dc:"报名结束时间"`
	StartTime         time.Time                   `json:"start_time" v:"required#展会开始时间不能为空" dc:"展会开始时间"`
	EndTime           time.Time                   `json:"end_time" v:"required#展会结束时间不能为空" dc:"展会结束时间"`
	Files             []*CreateExhibitionFileInfo `json:"files" dc:"展会相关文件(Logo、Banner、Product Image、Poster、Document、Video)"`
}

type CreateExhibitionRes struct {
	g.Meta
	ID string `json:"id" dc:"展会ID"`
}

type GetExhibitionReq struct {
	g.Meta `path:"/exhibitions/{id}" method:"get" tags:"展会管理" summary:"获取展会详情"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type GetExhibitionRes struct {
	g.Meta
	*ExhibitionUnit
}

type ExhibitionUnit struct {
	ExhibitionInfo
	Files      []*FileInfo      `json:"files" dc:"展会相关文件(Logo、Banner、Product Image、Poster、Document、Video)"`
	Organizers []*OrganizerInfo `json:"organizers" dc:"主办方信息"`
}

type ListExhibitionsReq struct {
	g.Meta `path:"/exhibitions" method:"get" tags:"展会管理" summary:"列表展会"`
	Name   string `p:"name" dc:"展会名称(模糊搜索)"`
	model.PageReq
}

type ListExhibitionsRes struct {
	g.Meta
	List    []*ExhibitionUnit `json:"list" dc:"展会列表"`
	PageRes *model.PageRes    `json:"page_res" dc:"分页响应"`
}

// 展会状态管理相关接口
type SubmitExhibitionForReviewReq struct {
	g.Meta `path:"/exhibitions/{id}/submit-review" method:"patch" tags:"展会管理" summary:"提交展会审核"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type SubmitExhibitionForReviewRes struct {
	g.Meta
}

type GetSubmitForReviewListReq struct {
	g.Meta `path:"/exhibitions/pending-list" method:"get" tags:"展会管理" summary:"获取待审核列表"`
	model.PageReq
}

type GetSubmitForReviewListRes struct {
	g.Meta
	List    []*ExhibitionUnit `json:"list" dc:"展会列表"`
	PageRes *model.PageRes    `json:"page_res" dc:"分页响应"`
}

type ApproveExhibitionReq struct {
	g.Meta `path:"/exhibitions/{id}/approve" method:"patch" tags:"展会管理" summary:"审核通过"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type ApproveExhibitionRes struct {
	g.Meta
}

type RejectExhibitionReq struct {
	g.Meta `path:"/exhibitions/{id}/reject" method:"patch" tags:"展会管理" summary:"审核驳回"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type RejectExhibitionRes struct {
	g.Meta
}

type CancelExhibitionReq struct {
	g.Meta `path:"/exhibitions/{id}/cancel" method:"patch" tags:"展会管理" summary:"取消展会"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type CancelExhibitionRes struct {
	g.Meta
}

type ExhibitionInfo struct {
	ID                string `json:"id" dc:"展会ID"`
	ServiceProviderID string `json:"service_provider_id" dc:"服务提供商ID"`
	Title             string `json:"title" dc:"展会名称"`
	Status            string `json:"status"`
	Industry          string `json:"industry" dc:"所属行业"`
	Tags              string `json:"tags" dc:"展会标签"`
	Website           string `json:"website" dc:"展会官网"`
	Venue             string `json:"venue" dc:"展会地点"`
	VenueAddress      string `json:"venue_address" dc:"展会详细地址"`
	Country           string `json:"country" dc:"国家"`
	City              string `json:"city" dc:"城市"`
	Description       string `json:"description" dc:"展会描述"`
	RegistrationStart string `json:"registration_start" dc:"报名开始时间"`
	RegistrationEnd   string `json:"registration_end" dc:"报名结束时间"`
	StartTime         string `json:"start_time" dc:"展会开始时间"`
	EndTime           string `json:"end_time" dc:"展会结束时间"`
	CreateTime        string `json:"create_time" dc:"创建时间"`
	UpdateTime        string `json:"update_time" dc:"更新时间"`
}
