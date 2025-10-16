package model

import (
	"ExhibitionService/internal/model/entity"
	"time"
)

type MerchantStatus int

const (
	MerchantStatusPending  MerchantStatus = iota // 待审核
	MerchantStatusApproved                       // 已审核
	MerchantStatusDisabled                       // 已禁用
)

func GetMerchantStatusText(status MerchantStatus) string {
	switch status {
	case MerchantStatusPending:
		return "待审核"
	case MerchantStatusApproved:
		return "已审核"
	case MerchantStatusDisabled:
		return "已禁用"
	default:
		return "未知状态"
	}
}

type Merchant struct {
	ID                 string         `json:"id"`
	CompanyID          string         `json:"company_id"`
	ExhibitionID       string         `json:"exhibition_id"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	BoothNumber        string         `json:"booth_number"`
	ContactPersonName  string         `json:"contact_person_name"`
	ContactPersonPhone string         `json:"contact_person_phone"`
	ContactPersonEmail string         `json:"contact_person_email"`
	Status             MerchantStatus `json:"status"`
	Version            int64          `json:"version"`
	CreateTime         time.Time      `json:"create_time"`
	UpdateTime         time.Time      `json:"update_time"`
}

func ConvertMerchant(in *entity.TMerchant) *Merchant {
	return &Merchant{
		ID:                 in.ID,
		CompanyID:          in.CompanyID,
		ExhibitionID:       in.ExhibitionID,
		Name:               in.Name,
		Description:        in.Description,
		BoothNumber:        in.BoothNumber,
		ContactPersonName:  in.ContactPersonName,
		ContactPersonPhone: in.ContactPersonPhone,
		ContactPersonEmail: in.ContactPersonEmail,
		Status:             MerchantStatus(in.Status),
		Version:            in.Version,
		CreateTime:         time.Unix(in.CreateTime, 0),
		UpdateTime:         time.Unix(in.UpdateTime, 0),
	}
}

type CreateMerchantReq struct {
	CompanyID          string `json:"company_id" v:"required#公司ID不能为空" dc:"公司ID"`
	ExhibitionID       string `json:"exhibition_id" v:"required#展会ID不能为空" dc:"展会ID"`
	Name               string `json:"name" v:"required#展商名称不能为空" dc:"展商名称"`
	Description        string `json:"description" dc:"展商描述"`
	BoothNumber        string `json:"booth_number" dc:"展位号"`
	ContactPersonName  string `json:"contact_person_name" dc:"联系人姓名"`
	ContactPersonPhone string `json:"contact_person_phone" dc:"联系人电话"`
	ContactPersonEmail string `json:"contact_person_email" dc:"联系人邮箱"`
}

type GetMerchantReq struct {
	ID string `json:"id"`
}

type GetMerchantRes struct {
	Merchant *Merchant `json:"merchant"`
}

type ListMerchantsReq struct {
	ExhibitionID string   `json:"exhibition_id"`
	Name         string   `json:"name"`
	PageReq      *PageReq `json:"page_req"`
}

type ListMerchantsRes struct {
	Merchants []*Merchant `json:"merchants"`
	PageRes   *PageRes    `json:"page_res"`
}
