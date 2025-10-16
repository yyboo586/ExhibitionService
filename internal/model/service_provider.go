package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

type ServiceProviderStatus int

const (
	ServiceProviderStatusPending      ServiceProviderStatus = iota // 待审核
	ServiceProviderStatusApproved                                  // 已审核
	ServiceProviderStatusRejected                                  // 已驳回
	ServiceProviderStatusDisabled                                  // 已禁用
	ServiceProviderStatusUnregistered                              // 已注销
)

func GetServiceProviderStatusText(status ServiceProviderStatus) string {
	switch status {
	case ServiceProviderStatusPending:
		return "待审核"
	case ServiceProviderStatusApproved:
		return "已审核"
	case ServiceProviderStatusRejected:
		return "已驳回"
	case ServiceProviderStatusDisabled:
		return "已禁用"
	case ServiceProviderStatusUnregistered:
		return "已注销"
	default:
		return "未知状态"
	}
}

// 服务提供商事件类型
type ServiceProviderEvent uint8

const (
	ServiceProviderEventReCommit   ServiceProviderEvent = iota // 重新提交审核
	ServiceProviderEventApprove                                // 审核通过
	ServiceProviderEventReject                                 // 审核驳回
	ServiceProviderEventDisable                                // 禁用
	ServiceProviderEventEnable                                 // 启用
	ServiceProviderEventUnregister                             // 注销
)

func GetServiceProviderEventText(event ServiceProviderEvent) string {
	switch event {
	case ServiceProviderEventReCommit:
		return "重新提交审核"
	case ServiceProviderEventApprove:
		return "审核通过"
	case ServiceProviderEventReject:
		return "审核驳回"
	case ServiceProviderEventDisable:
		return "禁用"
	case ServiceProviderEventEnable:
		return "启用"
	case ServiceProviderEventUnregister:
		return "注销"
	default:
		return "未知事件"
	}
}

type ServiceProvider struct {
	ID                 string                `json:"id"`
	CompanyID          string                `json:"company_id"`
	Name               string                `json:"name"`
	Status             ServiceProviderStatus `json:"status"`
	Website            string                `json:"website"`
	ContactPersonName  string                `json:"contact_person_name"`
	ContactPersonPhone string                `json:"contact_person_phone"`
	ContactPersonEmail string                `json:"contact_person_email"`
	Description        string                `json:"description"`

	Version int64 `json:"version"`

	Files []*File `json:"files"`

	CreateTime          time.Time `json:"create_time"`
	SubmitForReviewTime time.Time `json:"submit_for_review_time"`
	ApproveTime         time.Time `json:"approve_time"`
	UpdateTime          time.Time `json:"update_time"`

	CompanyInfo *Company `json:"company_info"`
}

func ConvertServiceProvider(in *entity.TServiceProvider) *ServiceProvider {
	return &ServiceProvider{
		ID:                 in.ID,
		CompanyID:          in.CompanyID,
		Name:               in.Name,
		Status:             ServiceProviderStatus(in.Status),
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
