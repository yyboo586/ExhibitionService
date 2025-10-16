package entity

// 展会表实体
type TExhibition struct {
	ID                string `orm:"id"`
	ServiceProviderID string `orm:"service_provider_id"`
	Title             string `orm:"title"`
	Status            int    `orm:"status"`
	Industry          string `orm:"industry"`
	Tags              string `orm:"tags"`
	Website           string `orm:"website"`
	Venue             string `orm:"venue"`
	VenueAddress      string `orm:"venue_address"`
	Country           string `orm:"country"`
	City              string `orm:"city"`
	Description       string `orm:"description"`
	RegistrationStart int64  `orm:"registration_start"`
	RegistrationEnd   int64  `orm:"registration_end"`
	StartTime         int64  `orm:"start_time"`
	EndTime           int64  `orm:"end_time"`
	Version           int64  `orm:"version"`
	CreateTime        int64  `orm:"create_time"`
	UpdateTime        int64  `orm:"update_time"`
}
