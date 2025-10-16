package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type TServiceProvider struct {
	table   string
	group   string // group is the database configuration group name of current DAO.
	columns TServiceProviderColumns
}

type TServiceProviderColumns struct {
	ID                 string // 服务提供商ID
	CompanyID          string // 关联的公司ID
	Name               string // 服务提供商名称
	Status             string // 状态
	Website            string // 官网
	ContactPersonName  string // 联系人姓名
	ContactPersonPhone string // 联系人电话
	ContactPersonEmail string // 联系人邮箱
	Description        string // 服务提供商描述
	Version            string // 版本号

	CreateTime          string // 创建时间
	SubmitForReviewTime string // 提交审核时间
	ApproveTime         string // 审核通过时间
	UpdateTime          string // 更新时间
}

var serviceProviderColumns = TServiceProviderColumns{
	ID:                 "id",
	CompanyID:          "company_id",
	Name:               "name",
	Status:             "status",
	Website:            "website",
	ContactPersonName:  "contact_person_name",
	ContactPersonPhone: "contact_person_phone",
	ContactPersonEmail: "contact_person_email",
	Description:        "description",
	Version:            "version",

	CreateTime:          "create_time",
	SubmitForReviewTime: "submit_for_review_time",
	ApproveTime:         "approve_time",
	UpdateTime:          "update_time",
}

func NewTServiceProvider() *TServiceProvider {
	return &TServiceProvider{
		group:   "default",
		table:   "t_service_provider",
		columns: serviceProviderColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TServiceProvider) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TServiceProvider) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TServiceProvider) Columns() TServiceProviderColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *TServiceProvider) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TServiceProvider) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TServiceProvider) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
