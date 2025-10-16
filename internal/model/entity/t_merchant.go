package entity

// 展商表实体
type TMerchant struct {
	ID                 string `orm:"id"`
	CompanyID          string `orm:"company_id"`
	ExhibitionID       string `orm:"exhibition_id"`
	Name               string `orm:"name"`
	Description        string `orm:"description"`
	BoothNumber        string `orm:"booth_number"`
	ContactPersonName  string `orm:"contact_person_name"`
	ContactPersonPhone string `orm:"contact_person_phone"`
	ContactPersonEmail string `orm:"contact_person_email"`
	Status             int    `orm:"status"`
	Version            int64  `orm:"version"`
	CreateTime         int64  `orm:"create_time"`
	UpdateTime         int64  `orm:"update_time"`
}
