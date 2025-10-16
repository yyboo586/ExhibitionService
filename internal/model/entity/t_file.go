package entity

type TFile struct {
	ID         string `orm:"id"`
	Module     int    `orm:"module"`
	CustomID   string `orm:"custom_id"`
	Type       int    `orm:"type"`
	FileID     string `orm:"file_id"`
	FileName   string `orm:"file_name"`
	FileLink   string `orm:"file_link"`
	Status     int    `orm:"status"`
	CreateTime int64  `orm:"create_time"`
	UpdateTime int64  `orm:"update_time"`
}
