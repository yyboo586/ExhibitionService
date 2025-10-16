package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type TMerchant struct {
	table   string
	group   string // group is the database configuration group name of current DAO.
	columns TMerchantColumns
}

type TMerchantColumns struct {
	ID                 string // 展商ID
	CompanyID          string // 关联的公司ID
	ExhibitionID       string // 关联的展会ID
	Name               string // 展商名称
	Description        string // 展商描述
	BoothNumber        string // 展位号
	ContactPersonName  string // 联系人姓名
	ContactPersonPhone string // 联系人电话
	ContactPersonEmail string // 联系人邮箱
	Status             string // 状态
	Version            string // 版本号
	CreateTime         string // 创建时间
	UpdateTime         string // 更新时间
}

var merchantColumns = TMerchantColumns{
	ID:                 "id",
	CompanyID:          "company_id",
	ExhibitionID:       "exhibition_id",
	Name:               "name",
	Description:        "description",
	BoothNumber:        "booth_number",
	ContactPersonName:  "contact_person_name",
	ContactPersonPhone: "contact_person_phone",
	ContactPersonEmail: "contact_person_email",
	Status:             "status",
	Version:            "version",
	CreateTime:         "create_time",
	UpdateTime:         "update_time",
}

func NewTMerchant() *TMerchant {
	return &TMerchant{
		group:   "default",
		table:   "t_merchant",
		columns: merchantColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TMerchant) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TMerchant) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TMerchant) Columns() TMerchantColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *TMerchant) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TMerchant) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TMerchant) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
