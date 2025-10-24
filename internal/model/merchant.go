package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

type MerchantStatus int

const (
	MerchantStatusPending      MerchantStatus = iota // 待审核
	MerchantStatusApproved                           // 已审核
	MerchantStatusRejected                           // 审核驳回
	MerchantStatusDisabled                           // 已禁用
	MerchantStatusUnregistered                       // 注销
)

func GetMerchantStatusText(status MerchantStatus) string {
	switch status {
	case MerchantStatusPending:
		return "待审核"
	case MerchantStatusApproved:
		return "已审核"
	case MerchantStatusRejected:
		return "审核驳回"
	case MerchantStatusDisabled:
		return "已禁用"
	case MerchantStatusUnregistered:
		return "注销"
	default:
		return "未知状态"
	}
}

// 展商事件类型
type MerchantEvent uint8

const (
	MerchantEventReCommit   MerchantEvent = iota // 重新提交审核
	MerchantEventApprove                         // 审核通过
	MerchantEventReject                          // 审核驳回
	MerchantEventDisable                         // 禁用
	MerchantEventEnable                          // 启用
	MerchantEventUnregister                      // 注销
)

func GetMerchantEventText(event MerchantEvent) string {
	switch event {
	case MerchantEventApprove:
		return "审核通过"
	case MerchantEventReject:
		return "审核驳回"
	case MerchantEventReCommit:
		return "重新提交审核"
	case MerchantEventDisable:
		return "禁用"
	case MerchantEventEnable:
		return "启用"
	case MerchantEventUnregister:
		return "注销"
	default:
		return "未知事件"
	}
}

type Merchant struct {
	ID                  string         `json:"id"`
	CompanyID           string         `json:"company_id"`
	Name                string         `json:"name"`
	Status              MerchantStatus `json:"status"`
	Website             string         `json:"website"`
	ContactPersonName   string         `json:"contact_person_name"`
	ContactPersonPhone  string         `json:"contact_person_phone"`
	ContactPersonEmail  string         `json:"contact_person_email"`
	Description         string         `json:"description"`
	Version             int64          `json:"version"`
	CreateTime          time.Time      `json:"create_time"`
	SubmitForReviewTime time.Time      `json:"submit_for_review_time"`
	ApproveTime         time.Time      `json:"approve_time"`
	UpdateTime          time.Time      `json:"update_time"`

	// 关联信息
	CompanyInfo *Company `json:"company_info,omitempty"`
	Files       []*File  `json:"files,omitempty"`
}

func ConvertMerchant(in *entity.TMerchant) *Merchant {
	return &Merchant{
		ID:                 in.ID,
		CompanyID:          in.CompanyID,
		Name:               in.Name,
		Status:             MerchantStatus(in.Status),
		Website:            in.Website,
		ContactPersonName:  in.ContactPersonName,
		ContactPersonPhone: in.ContactPersonPhone,
		ContactPersonEmail: in.ContactPersonEmail,
		Description:        in.Description,
		Version:            in.Version,

		CreateTime:          time.Unix(in.CreateTime, 0),
		SubmitForReviewTime: time.Unix(in.SubmitForReviewTime, 0),
		ApproveTime:         time.Unix(in.ApproveTime, 0),
		UpdateTime:          time.Unix(in.UpdateTime, 0),
	}
}
