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
	ExhibitionEventStartRunning                           // 开始进行(由定时任务触发)
	ExhibitionEventEnd                                    // 结束展会(由定时任务触发)
	ExhibitionEventCancel                                 // 取消展会
)

type ExhibitionStatus int

const (
	ExhibitionStatusPreparing ExhibitionStatus = iota // 筹备中
	ExhibitionStatusPending                           // 待审核
	ExhibitionStatusApproved                          // 已批准
	ExhibitionStatusEnrolling                         // 报名中
	ExhibitionStatusRunning                           // 进行中
	ExhibitionStatusEnded                             // 已结束
	ExhibitionStatusCancelled                         // 已取消
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
