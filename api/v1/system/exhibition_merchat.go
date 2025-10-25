package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 商户申请参加展会相关接口
type ApplyForExhibitionReq struct {
	g.Meta `path:"/exhibitions/{exhibition_id}/apply" method:"post" tags:"展会管理" summary:"商户申请参加展会"`
	model.AuthorRequired
	ExhibitionID string `p:"exhibition_id" v:"required#展会ID不能为空" dc:"展会ID"`
	MerchantID   string `json:"merchant_id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type ApplyForExhibitionRes struct {
	g.Meta
	Message string `json:"message" dc:"申请结果消息"`
}

type GetMerchantApplicationReq struct {
	g.Meta       `path:"/exhibitions/{exhibition_id}/applications/{merchant_id}" method:"get" tags:"展会管理" summary:"获取商户申请状态"`
	ExhibitionID string `p:"exhibition_id" v:"required#展会ID不能为空" dc:"展会ID"`
	MerchantID   string `p:"merchant_id" v:"required#商户ID不能为空" dc:"商户ID"`
}

type GetMerchantApplicationRes struct {
	g.Meta
	*ExhibitionMerchantInfo
}

type ListExhibitionApplicationsReq struct {
	g.Meta       `path:"/exhibitions/{exhibition_id}/applications" method:"get" tags:"展会管理" summary:"获取展会申请列表"`
	ExhibitionID string `p:"exhibition_id" v:"required#展会ID不能为空" dc:"展会ID"`
	model.PageReq
}

type ListExhibitionApplicationsRes struct {
	g.Meta
	List    []*ExhibitionMerchantInfo `json:"list" dc:"申请列表"`
	PageRes *model.PageRes            `json:"page_res" dc:"分页响应"`
}

type ListMerchantApplicationsReq struct {
	g.Meta     `path:"/merchants/{merchant_id}/applications" method:"get" tags:"商户管理" summary:"获取商户申请列表"`
	MerchantID string `p:"merchant_id" v:"required#商户ID不能为空" dc:"商户ID"`
	model.PageReq
}

type ListMerchantApplicationsRes struct {
	g.Meta
	List    []*ExhibitionMerchantInfo `json:"list" dc:"申请列表"`
	PageRes *model.PageRes            `json:"page_res" dc:"分页响应"`
}

type ExhibitionMerchantInfo struct {
	ID                  int64  `json:"id" dc:"申请ID"`
	ExhibitionID        string `json:"exhibition_id" dc:"展会ID"`
	MerchantID          string `json:"merchant_id" dc:"商户ID"`
	Status              string `json:"status" dc:"申请状态"`
	CreateTime          string `json:"create_time" dc:"申请时间"`
	SubmitForReviewTime string `json:"submit_for_review_time" dc:"提交审核时间"`
	ApproveTime         string `json:"approve_time" dc:"审核通过时间"`
	UpdateTime          string `json:"update_time" dc:"更新时间"`

	// 关联信息
	Exhibition *ExhibitionInfo `json:"exhibition,omitempty" dc:"展会信息"`
	Merchant   *MerchantInfo   `json:"merchant,omitempty" dc:"商户信息"`
}
