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

type Exhibition struct {
	ID                string           `json:"id"`
	ServiceProviderID string           `json:"service_provider_id"`
	Title             string           `json:"title"`
	Status            ExhibitionStatus `json:"status"`
	Industry          string           `json:"industry"`
	Tags              string           `json:"tags"`
	Website           string           `json:"website"`
	Venue             string           `json:"venue"`
	VenueAddress      string           `json:"venue_address"`
	Country           string           `json:"country"`
	City              string           `json:"city"`
	Description       string           `json:"description"`
	RegistrationStart time.Time        `json:"registration_start"`
	RegistrationEnd   time.Time        `json:"registration_end"`
	StartTime         time.Time        `json:"start_time"`
	EndTime           time.Time        `json:"end_time"`
	Version           int64            `json:"version"`
	CreateTime        time.Time        `json:"create_time"`
	UpdateTime        time.Time        `json:"update_time"`
}

func ConvertExhibition(in *entity.TExhibition) *Exhibition {
	return &Exhibition{
		ID:                in.ID,
		ServiceProviderID: in.ServiceProviderID,
		Title:             in.Title,
		Status:            ExhibitionStatus(in.Status),
		Industry:          in.Industry,
		Tags:              in.Tags,
		Website:           in.Website,
		Venue:             in.Venue,
		VenueAddress:      in.VenueAddress,
		Country:           in.Country,
		City:              in.City,
		Description:       in.Description,
		RegistrationStart: time.Unix(in.RegistrationStart, 0),
		RegistrationEnd:   time.Unix(in.RegistrationEnd, 0),
		StartTime:         time.Unix(in.StartTime, 0),
		EndTime:           time.Unix(in.EndTime, 0),
		Version:           in.Version,
		CreateTime:        time.Unix(in.CreateTime, 0),
		UpdateTime:        time.Unix(in.UpdateTime, 0),
	}
}

type GetExhibitionReq struct {
	ID string `json:"id"`
}

type GetExhibitionRes struct {
	Exhibition *Exhibition `json:"exhibition"`
}

type ListExhibitionsReq struct {
	Name    string   `json:"name"`
	PageReq *PageReq `json:"page_req"`
}

type ListExhibitionsRes struct {
	Exhibitions []*Exhibition `json:"exhibitions"`
	PageRes     *PageRes      `json:"page_res"`
}

type OrganizerRole int

const (
	OrganizerRoleOrganizer   OrganizerRole = iota // 主办方
	OrganizerRoleCoorganizer                      // 联合主办方
	OrganizerRoleSponsor                          // 协办方
)

type OrganizerReq struct {
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	LogoURL            string        `json:"logo_url"`
	RoleType           OrganizerRole `json:"role_type"`
	ContactPersonName  string        `json:"contact_person_name"`
	ContactPersonPhone string        `json:"contact_person_phone"`
	ContactPersonEmail string        `json:"contact_person_email"`
	Website            string        `json:"website"`
}
