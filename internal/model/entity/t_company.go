package entity

type TCompany struct {
	ID          string `orm:"id"`
	Name        string `orm:"name"`
	Country     string `orm:"country"`
	Status      int    `orm:"status"`
	Phone       string `orm:"phone"`
	Email       string `orm:"email"`
	Address     string `orm:"address"`
	Description string `orm:"description"`

	BusinessLicense       string `orm:"business_license"`
	SocialCreditCode      string `orm:"social_credit_code"`
	LegalPersonName       string `orm:"legal_person_name"`
	LegalPersonCardNumber string `orm:"legal_person_card_number"`
	LegalPersonPhotoUrl   string `orm:"legal_person_photo_url"`
	LegalPersonPhone      string `orm:"legal_person_phone"`

	ApplyTime   int64 `orm:"apply_time"`
	ApproveTime int64 `orm:"approve_time"`
	CreateTime  int64 `orm:"create_time"`
	UpdateTime  int64 `orm:"update_time"`
}
