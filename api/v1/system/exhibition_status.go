package system

import (
	"ExhibitionService/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 展会状态管理相关接口
type SubmitExhibitionForReviewReq struct {
	g.Meta `path:"/exhibitions/{id}/submit-review" method:"patch" tags:"展会管理" summary:"提交展会审核"`
	ID     string `p:"id" v:"required#展会ID不能为空" dc:"展会ID"`
}

type SubmitExhibitionForReviewRes struct {
	g.Meta
}

type GetPendingListReq struct {
	g.Meta `path:"/exhibitions/pending-list" method:"get" tags:"展会管理" summary:"获取待审核列表"`
	model.PageReq
}

type GetPendingListRes struct {
	g.Meta
	List    []*ExhibitionInfo `json:"list" dc:"展会列表"`
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
