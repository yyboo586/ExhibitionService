package entity

// 展会与商户关联表实体
type TExhibitionMerchant struct {
	ID                  int64  `orm:"id"`
	ExhibitionID        string `orm:"exhibition_id"`
	MerchantID          string `orm:"merchant_id"`
	Status              int    `orm:"status"`
	CreateTime          int64  `orm:"create_time"`
	SubmitForReviewTime int64  `orm:"submit_for_review_time"`
	ApproveTime         int64  `orm:"approve_time"`
	UpdateTime          int64  `orm:"update_time"`
}
