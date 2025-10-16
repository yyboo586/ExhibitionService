package system

type CompanyInfo struct {
	ID          string `json:"id" dc:"公司ID"`
	Name        string `json:"name" v:"required#公司名称不能为空" dc:"公司名称(营业执照上的公司名称)"`
	Country     string `json:"country" v:"required#国家不能为空" dc:"国家"`
	City        string `json:"city" v:"required#城市不能为空" dc:"城市"`
	Address     string `json:"address" v:"required#地址不能为空" dc:"地址"`
	Email       string `json:"email" v:"required#邮箱不能为空" dc:"邮箱"`
	Description string `json:"description" v:"required#公司描述不能为空" dc:"公司描述"`

	BusinessLicense       string `json:"business_license" dc:"营业执照"`
	SocialCreditCode      string `json:"social_credit_code" v:"required#统一社会信用代码不能为空" dc:"统一社会信用代码"`
	LegalPersonName       string `json:"legal_person_name" v:"required#法人姓名不能为空" dc:"法人姓名"`
	LegalPersonCardNumber string `json:"legal_person_card_number" v:"required#法人证件号不能为空" dc:"法人证件号"`
	LegalPersonPhoto      string `json:"legal_person_photo" dc:"法人证件照"`

	Files []*FileInfo `json:"files" v:"required#公司相关文件不能为空" dc:"公司相关文件"`
}
