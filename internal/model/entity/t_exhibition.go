package entity

// 展会表实体
type TExhibition struct {
	ID           string `orm:"id"`
	Title        string `orm:"title"`
	Website      string `orm:"website"`
	Status       int    `orm:"status"`
	Industry     string `orm:"industry"`
	Tags         string `orm:"tags"`
	Country      string `orm:"country"`
	City         string `orm:"city"`
	Venue        string `orm:"venue"`
	VenueAddress string `orm:"venue_address"`
	Description  string `orm:"description"`
	Version      int64  `orm:"version"`

	RegistrationStart int64 `orm:"registration_start"`
	RegistrationEnd   int64 `orm:"registration_end"`
	StartTime         int64 `orm:"start_time"`
	EndTime           int64 `orm:"end_time"`

	CreateTime          int64 `orm:"create_time"`
	SubmitForReviewTime int64 `orm:"submit_for_review_time"`
	ApproveTime         int64 `orm:"approve_time"`
	UpdateTime          int64 `orm:"update_time"`
}

type TExOrganizer struct {
	ID                int64  `orm:"id"`
	ExhibitionID      string `orm:"exhibition_id"`
	ServiceProviderID string `orm:"service_provider_id"`
	RoleType          int    `orm:"role_type"`
	CreateTime        int64  `orm:"create_time"`
	UpdateTime        int64  `orm:"update_time"`
}
