package entity

type TCompany struct {
	ID          string `orm:"id"`
	Name        string `orm:"name"`
	Type        int    `orm:"type"`
	Country     string `orm:"country"`
	City        string `orm:"city"`
	Address     string `orm:"address"`
	Email       string `orm:"email"`
	Description string `orm:"description"`
	Version     int64  `orm:"version"`

	SocialCreditCode      string `orm:"social_credit_code"`
	LegalPersonName       string `orm:"legal_person_name"`
	LegalPersonCardNumber string `orm:"legal_person_card_number"`

	CreateTime int64 `orm:"create_time"`
	UpdateTime int64 `orm:"update_time"`
}
