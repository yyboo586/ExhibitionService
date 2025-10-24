package entity

// 商户表实体
type TMerchant struct {
	ID                  string `orm:"id"`
	CompanyID           string `orm:"company_id"`
	Name                string `orm:"name"`
	Status              int    `orm:"status"`
	Website             string `orm:"website"`
	ContactPersonName   string `orm:"contact_person_name"`
	ContactPersonPhone  string `orm:"contact_person_phone"`
	ContactPersonEmail  string `orm:"contact_person_email"`
	Description         string `orm:"description"`
	Version             int64  `orm:"version"`
	CreateTime          int64  `orm:"create_time"`
	SubmitForReviewTime int64  `orm:"submit_for_review_time"`
	ApproveTime         int64  `orm:"approve_time"`
	UpdateTime          int64  `orm:"update_time"`
}
