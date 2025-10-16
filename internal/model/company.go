package model

import (
	"ExhibitionService/internal/model/entity"
	"errors"
	"time"
)

type CompanyType int

const (
	_                          CompanyType = iota
	CompanyTypeServiceProvider             // 服务提供商
	CompanyTypeMerchant                    // 商户
)

func GetCompanyTypeText(typ CompanyType) string {
	switch typ {
	case CompanyTypeServiceProvider:
		return "服务提供商"
	case CompanyTypeMerchant:
		return "商户"
	default:
		return "未知类型"
	}
}

type CompanyStatus int

const (
	CompanyStatusInit         CompanyStatus = iota // 初始状态
	CompanyStatusInReview                          // 待审核
	CompanyStatusApproved                          // 审核通过
	CompanyStatusRejected                          // 审核驳回
	CompanyStatusDisabled                          // 禁用
	CompanyStatusUnregistered                      // 注销
)

func GetCompanyStatusText(status CompanyStatus) string {
	switch status {
	case CompanyStatusInit:
		return "初始状态"
	case CompanyStatusInReview:
		return "待审核"
	case CompanyStatusApproved:
		return "审核通过"
	case CompanyStatusRejected:
		return "审核驳回"
	case CompanyStatusDisabled:
		return "禁用"
	case CompanyStatusUnregistered:
		return "注销"
	default:
		return "未知状态"
	}
}

func GetCompanyStatus(statusText string) (status CompanyStatus, err error) {
	switch statusText {
	case "初始状态":
		status = CompanyStatusInit
	case "待审核":
		status = CompanyStatusInReview
	case "审核通过":
		status = CompanyStatusApproved
	case "审核驳回":
		status = CompanyStatusRejected
	case "禁用":
		status = CompanyStatusDisabled
	case "注销":
		status = CompanyStatusUnregistered
	default:
		err = errors.New("未知状态")
	}
	return
}

type Company struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        CompanyType `json:"type"`
	Country     string      `json:"country"`
	City        string      `json:"city"`
	Address     string      `json:"address"`
	Email       string      `json:"email"`
	Description string      `json:"description"`
	Version     int64       `json:"version"`

	BusinessLicense       string `json:"business_license"`
	SocialCreditCode      string `json:"social_credit_code"`
	LegalPersonName       string `json:"legal_person_name"`
	LegalPersonCardNumber string `json:"legal_person_card_number"`
	LegalPersonPhoto      string `json:"legal_person_photo_url"`

	Files []*File `json:"files"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func ConvertCompany(in *entity.TCompany) (out *Company) {
	out = &Company{
		ID:          in.ID,
		Name:        in.Name,
		Type:        CompanyType(in.Type),
		Country:     in.Country,
		City:        in.City,
		Address:     in.Address,
		Email:       in.Email,
		Description: in.Description,
		Version:     in.Version,

		SocialCreditCode:      in.SocialCreditCode,
		LegalPersonName:       in.LegalPersonName,
		LegalPersonCardNumber: in.LegalPersonCardNumber,

		CreateTime: time.Unix(in.CreateTime, 0),
		UpdateTime: time.Unix(in.UpdateTime, 0),
	}

	return
}
