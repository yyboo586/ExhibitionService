package model

import (
	"ExhibitionService/internal/model/entity"
	"errors"
	"time"
)

type CompanyStatus int

const (
	CompanyStatusPending      CompanyStatus = iota // 待审核
	CompanyStatusApproved                          // 审核通过
	CompanyStatusDisabled                          // 禁用
	CompanyStatusUnregistered                      // 注销
)

func GetCompanyStatusText(status CompanyStatus) string {
	switch status {
	case CompanyStatusPending:
		return "待审核"
	case CompanyStatusApproved:
		return "审核通过"
	case CompanyStatusDisabled:
		return "禁用"
	case CompanyStatusUnregistered:
		return "审核驳回"
	default:
		return "未知状态"
	}
}

func GetCompanyStatus(statusText string) (status CompanyStatus, err error) {
	switch statusText {
	case "待审核":
		status = CompanyStatusPending
	case "审核通过":
		status = CompanyStatusApproved
	case "审核驳回":
		status = CompanyStatusUnregistered
	case "禁用":
		status = CompanyStatusDisabled
	default:
		err = errors.New("未知状态")
	}
	return
}

type Company struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Country     string        `json:"country"`
	Status      CompanyStatus `json:"status"`
	Phone       string        `json:"phone"`
	Email       string        `json:"email"`
	Address     string        `json:"address"`
	Description string        `json:"description"`

	BusinessLicense       string `json:"business_license"`
	SocialCreditCode      string `json:"social_credit_code"`
	LegalPersonName       string `json:"legal_person_name"`
	LegalPersonCardNumber string `json:"legal_person_card_number"`
	LegalPersonPhotoUrl   string `json:"legal_person_photo_url"`
	LegalPersonPhone      string `json:"legal_person_phone"`

	ApplyTime   time.Time `json:"apply_time"`
	ApproveTime time.Time `json:"approve_time"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func ConvertCompany(in *entity.TCompany) *Company {
	return &Company{
		ID:          in.ID,
		Name:        in.Name,
		Country:     in.Country,
		Status:      CompanyStatus(in.Status),
		Phone:       in.Phone,
		Email:       in.Email,
		Address:     in.Address,
		Description: in.Description,

		BusinessLicense:       in.BusinessLicense,
		SocialCreditCode:      in.SocialCreditCode,
		LegalPersonName:       in.LegalPersonName,
		LegalPersonCardNumber: in.LegalPersonCardNumber,
		LegalPersonPhotoUrl:   in.LegalPersonPhotoUrl,
		LegalPersonPhone:      in.LegalPersonPhone,

		ApplyTime:   time.Unix(in.ApplyTime, 0),
		ApproveTime: time.Unix(in.ApproveTime, 0),
		CreateTime:  time.Unix(in.CreateTime, 0),
		UpdateTime:  time.Unix(in.UpdateTime, 0),
	}
}
