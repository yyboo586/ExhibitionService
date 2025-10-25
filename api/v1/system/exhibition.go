package system

import (
	"ExhibitionService/internal/model"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type OrganizerInfo struct {
	ServiceProviderID string `json:"service_provider_id" v:"required#服务提供商ID不能为空" dc:"服务提供商ID"`
	RoleType          string `json:"role_type" v:"required#主办方角色类型不能为空" dc:"主办方角色类型(主办方、联合主办方、协办方)"`
}

type CreateExhibitionReq struct {
	g.Meta `path:"/exhibitions" method:"post" tags:"展会管理" summary:"创建展会"`
	model.AuthorRequired
	Organizers   []*OrganizerInfo `json:"organizers" v:"required#主办方不能为空" dc:"主办方"`
	Title        string           `json:"title" v:"required#展会名称不能为空" dc:"展会名称"`
	Website      string           `json:"website" dc:"展会官网"`
	Industry     string           `json:"industry" v:"required#所属行业不能为空" dc:"所属行业"`
	Tags         string           `json:"tags" v:"required#展会标签不能为空" dc:"展会标签"`
	Country      string           `json:"country" v:"required#国家不能为空" dc:"国家"`
	City         string           `json:"city" v:"required#城市不能为空" dc:"城市"`
	Venue        string           `json:"venue" v:"required#展会地点不能为空" dc:"展会地点"`
	VenueAddress string           `json:"venue_address" v:"required#展会详细地址不能为空" dc:"展会详细地址"`
	Description  string           `json:"description" v:"required#展会描述不能为空" dc:"展会描述"`

	RegistrationStart time.Time `json:"registration_start" v:"required#报名开始时间不能为空" dc:"报名开始时间"`
	RegistrationEnd   time.Time `json:"registration_end" v:"required#报名结束时间不能为空" dc:"报名结束时间"`
	StartTime         time.Time `json:"start_time" v:"required#展会开始时间不能为空" dc:"展会开始时间"`
	EndTime           time.Time `json:"end_time" v:"required#展会结束时间不能为空" dc:"展会结束时间"`

	Files []*FileInfo `json:"files" dc:"展会相关文件(Logo、Banner、Product Image、Poster、Document、Video)"`
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
	*ExhibitionInfo
}

type ListExhibitionsReq struct {
	g.Meta `path:"/exhibitions" method:"get" tags:"展会管理" summary:"列表展会"`
	Name   string `p:"name" dc:"展会名称(模糊搜索)"`
	model.PageReq
}

type ListExhibitionsRes struct {
	g.Meta
	List    []*ExhibitionInfo `json:"list" dc:"展会列表"`
	PageRes *model.PageRes    `json:"page_res" dc:"分页响应"`
}

type ExhibitionInfo struct {
	ID           string `json:"id" dc:"展会ID"`
	Title        string `json:"title" dc:"展会名称"`
	Website      string `json:"website" dc:"展会官网"`
	Status       string `json:"status"`
	Industry     string `json:"industry" dc:"所属行业"`
	Tags         string `json:"tags" dc:"展会标签"`
	Country      string `json:"country" dc:"国家"`
	City         string `json:"city" dc:"城市"`
	Venue        string `json:"venue" dc:"展会地点"`
	VenueAddress string `json:"venue_address" dc:"展会详细地址"`
	Description  string `json:"description" dc:"展会描述"`

	RegistrationStart string `json:"registration_start" dc:"报名开始时间"`
	RegistrationEnd   string `json:"registration_end" dc:"报名结束时间"`
	StartTime         string `json:"start_time" dc:"展会开始时间"`
	EndTime           string `json:"end_time" dc:"展会结束时间"`

	CreateTime          string `json:"create_time" dc:"创建时间"`
	SubmitForReviewTime string `json:"submit_for_review_time" dc:"提交审核时间"`
	ApproveTime         string `json:"approve_time" dc:"审核通过时间"`
	UpdateTime          string `json:"update_time" dc:"更新时间"`

	Files      []*FileInfo      `json:"files" dc:"展会相关文件(Logo、Banner、Product Image、Poster、Document、Video)"`
	Organizers []*OrganizerInfo `json:"organizers" dc:"主办方信息"`
}
