package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

// 展会事件类型
type ExhibitionEvent uint8

const (
	ExhibitionEventSubmitForReview ExhibitionEvent = iota // 提交审核
	ExhibitionEventApprove                                // 审核通过
	ExhibitionEventReject                                 // 审核驳回
	ExhibitionEventStartEnrolling                         // 开始报名(由定时任务触发)
	ExhibitionEventEndEnrolling                           // 结束报名(由定时任务触发)
	ExhibitionEventStartRunning                           // 开始进行(由定时任务触发)
	ExhibitionEventEnd                                    // 结束展会(由定时任务触发)
	ExhibitionEventCancel                                 // 取消展会
)

func GetExhibitionEventText(event ExhibitionEvent) string {
	switch event {
	case ExhibitionEventSubmitForReview:
		return "提交审核"
	case ExhibitionEventApprove:
		return "审核通过"
	case ExhibitionEventReject:
		return "审核驳回"
	case ExhibitionEventStartEnrolling:
		return "开始报名"
	case ExhibitionEventEndEnrolling:
		return "结束报名"
	case ExhibitionEventStartRunning:
		return "开始进行"
	case ExhibitionEventEnd:
		return "结束展会"
	case ExhibitionEventCancel:
		return "取消展会"
	default:
		return "未知事件"
	}
}

type ExhibitionStatus int

const (
	ExhibitionStatusPreparing      ExhibitionStatus = iota // 筹备中
	ExhibitionStatusPending                                // 待审核
	ExhibitionStatusApproved                               // 已批准
	ExhibitionStatusEnrolling                              // 报名中
	ExhibitionStatusEnrollingEnded                         // 报名结束
	ExhibitionStatusRunning                                // 进行中
	ExhibitionStatusEnded                                  // 已结束
	ExhibitionStatusCancelled                              // 已取消
)

func GetExhibitionStatusText(status ExhibitionStatus) string {
	switch status {
	case ExhibitionStatusPreparing:
		return "筹备中"
	case ExhibitionStatusPending:
		return "待审核"
	case ExhibitionStatusApproved:
		return "已批准"
	case ExhibitionStatusEnrolling:
		return "报名中"
	case ExhibitionStatusEnrollingEnded:
		return "报名结束"
	case ExhibitionStatusRunning:
		return "进行中"
	case ExhibitionStatusEnded:
		return "已结束"
	case ExhibitionStatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}

type OrganizerRole int

const (
	OrganizerRoleUnknown     OrganizerRole = iota // 未知角色
	OrganizerRoleOrganizer                        // 主办方
	OrganizerRoleCoorganizer                      // 联合主办方
	OrganizerRoleSponsor                          // 协办方
)

func GetOrganizerRoleText(role OrganizerRole) string {
	switch role {
	case OrganizerRoleOrganizer:
		return "主办方"
	case OrganizerRoleCoorganizer:
		return "联合主办方"
	case OrganizerRoleSponsor:
		return "协办方"
	default:
		return "未知角色"
	}
}

func GetOrganizerRole(roleText string) (role OrganizerRole) {
	switch roleText {
	case "主办方":
		return OrganizerRoleOrganizer
	case "联合主办方":
		return OrganizerRoleCoorganizer
	case "协办方":
		return OrganizerRoleSponsor
	default:
		return OrganizerRoleUnknown
	}
}

type ExOrganizer struct {
	ID                int64         `json:"id"`
	ExhibitionID      string        `json:"exhibition_id"`
	ServiceProviderID string        `json:"service_provider_id"`
	RoleType          OrganizerRole `json:"role_type"`
	CreateTime        time.Time     `json:"create_time"`
	UpdateTime        time.Time     `json:"update_time"`
}

type Exhibition struct {
	ID           string           `json:"id"`
	Organizers   []*ExOrganizer   `json:"organizers"`
	Title        string           `json:"title"`
	Website      string           `json:"website"`
	Status       ExhibitionStatus `json:"status"`
	Industry     string           `json:"industry"`
	Tags         string           `json:"tags"`
	Country      string           `json:"country"`
	City         string           `json:"city"`
	Venue        string           `json:"venue"`
	VenueAddress string           `json:"venue_address"`
	Description  string           `json:"description"`
	Version      int64            `json:"version"`

	RegistrationStart time.Time `json:"registration_start"`
	RegistrationEnd   time.Time `json:"registration_end"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`

	CreateTime          time.Time `json:"create_time"`
	SubmitForReviewTime time.Time `json:"submit_for_review_time"`
	ApproveTime         time.Time `json:"approve_time"`
	UpdateTime          time.Time `json:"update_time"`

	Files []*File `json:"files"`
}

func ConvertExhibition(in *entity.TExhibition) *Exhibition {
	return &Exhibition{
		ID:           in.ID,
		Title:        in.Title,
		Status:       ExhibitionStatus(in.Status),
		Website:      in.Website,
		Country:      in.Country,
		City:         in.City,
		Venue:        in.Venue,
		VenueAddress: in.VenueAddress,
		Industry:     in.Industry,
		Tags:         in.Tags,
		Description:  in.Description,
		Version:      in.Version,

		RegistrationStart: time.Unix(in.RegistrationStart, 0),
		RegistrationEnd:   time.Unix(in.RegistrationEnd, 0),
		StartTime:         time.Unix(in.StartTime, 0),
		EndTime:           time.Unix(in.EndTime, 0),

		CreateTime:          time.Unix(in.CreateTime, 0),
		SubmitForReviewTime: time.Unix(in.SubmitForReviewTime, 0),
		ApproveTime:         time.Unix(in.ApproveTime, 0),
		UpdateTime:          time.Unix(in.UpdateTime, 0),
	}
}

func ConvertExOrganizer(in *entity.TExOrganizer) *ExOrganizer {
	return &ExOrganizer{
		ID:                in.ID,
		ExhibitionID:      in.ExhibitionID,
		ServiceProviderID: in.ServiceProviderID,
		RoleType:          OrganizerRole(in.RoleType),
		CreateTime:        time.Unix(in.CreateTime, 0),
		UpdateTime:        time.Unix(in.UpdateTime, 0),
	}
}

// 展会与商户关联状态
type ExMerchantStatus int

const (
	ExMerchantStatusPending   ExMerchantStatus = iota // 待审核
	ExMerchantStatusApproved                          // 审核通过
	ExMerchantStatusRejected                          // 审核拒绝
	ExMerchantStatusWithdrawn                         // 已退出
)

func GetExMerchantStatusText(status ExMerchantStatus) string {
	switch status {
	case ExMerchantStatusPending:
		return "待审核"
	case ExMerchantStatusApproved:
		return "审核通过"
	case ExMerchantStatusRejected:
		return "审核拒绝"
	case ExMerchantStatusWithdrawn:
		return "已退出"
	default:
		return "未知状态"
	}
}

// 展会与商户关联事件类型
type ExMerchantEvent uint8

const (
	_                       ExMerchantEvent = iota //
	ExMerchantEventApprove                         // 审核通过
	ExMerchantEventReject                          // 审核拒绝
	ExMerchantEventReApply                         // 重新申请
	ExMerchantEventWithdraw                        // 退出展会
)

func GetExMerchantEventText(event ExMerchantEvent) string {
	switch event {
	case ExMerchantEventApprove:
		return "审核通过"
	case ExMerchantEventReject:
		return "审核拒绝"
	case ExMerchantEventReApply:
		return "重新申请"
	case ExMerchantEventWithdraw:
		return "退出展会"
	default:
		return "未知事件"
	}
}

type ExhibitionMerchant struct {
	ID                  int64            `json:"id"`
	ExhibitionID        string           `json:"exhibition_id"`
	MerchantID          string           `json:"merchant_id"`
	Status              ExMerchantStatus `json:"status"`
	CreateTime          time.Time        `json:"create_time"`
	SubmitForReviewTime time.Time        `json:"submit_for_review_time"`
	ApproveTime         time.Time        `json:"approve_time"`
	UpdateTime          time.Time        `json:"update_time"`
}

func ConvertExhibitionMerchant(in *entity.TExhibitionMerchant) *ExhibitionMerchant {
	return &ExhibitionMerchant{
		ID:                  in.ID,
		ExhibitionID:        in.ExhibitionID,
		MerchantID:          in.MerchantID,
		Status:              ExMerchantStatus(in.Status),
		CreateTime:          time.Unix(in.CreateTime, 0),
		SubmitForReviewTime: time.Unix(in.SubmitForReviewTime, 0),
		ApproveTime:         time.Unix(in.ApproveTime, 0),
		UpdateTime:          time.Unix(in.UpdateTime, 0),
	}
}
