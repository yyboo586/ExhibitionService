package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type TExhibition struct {
	table   string
	group   string // group is the database configuration group name of current DAO.
	columns TExhibitionColumns
}

type TExhibitionColumns struct {
	ID                string // 展会ID
	ServiceProviderID string // 服务提供商ID
	Title             string // 展会标题
	Status            string // 展会状态
	Industry          string // 所属行业
	Tags              string // 展会标签
	Website           string // 展会官网
	Venue             string // 展会地点
	VenueAddress      string // 展会详细地址
	Country           string // 国家
	City              string // 城市
	Description       string // 展会描述
	RegistrationStart string // 报名开始时间
	RegistrationEnd   string // 报名结束时间
	StartTime         string // 展会开始时间
	EndTime           string // 展会结束时间
	Version           string // 版本号
	CreateTime        string // 创建时间
	UpdateTime        string // 更新时间
}

var exhibitionColumns = TExhibitionColumns{
	ID:                "id",
	ServiceProviderID: "service_provider_id",
	Title:             "title",
	Status:            "status",
	Industry:          "industry",
	Tags:              "tags",
	Website:           "website",
	Venue:             "venue",
	VenueAddress:      "venue_address",
	Country:           "country",
	City:              "city",
	Description:       "description",
	RegistrationStart: "registration_start",
	RegistrationEnd:   "registration_end",
	StartTime:         "start_time",
	EndTime:           "end_time",
	Version:           "version",
	CreateTime:        "create_time",
	UpdateTime:        "update_time",
}

func NewTExhibition() *TExhibition {
	return &TExhibition{
		group:   "default",
		table:   "t_exhibition",
		columns: exhibitionColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TExhibition) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TExhibition) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TExhibition) Columns() TExhibitionColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *TExhibition) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TExhibition) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TExhibition) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
