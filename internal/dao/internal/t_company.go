package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CompanyDao is the data access object for table t_company.
type CompanyDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns CompanyColumns // columns contains all the column names of Table for convenient usage.
}

// CompanyColumns defines and stores column names for table t_company.
type CompanyColumns struct {
	ID          string // 主键
	Name        string // 公司名称
	Type        string // 公司类型
	Country     string // 国家
	City        string // 城市
	Address     string // 地址
	Email       string // 邮箱
	Description string // 公司描述
	Version     string // 版本号

	SocialCreditCode      string // 统一社会信用代码
	LegalPersonName       string // 法人姓名
	LegalPersonCardNumber string // 法人证件号

	CreateTime string // 创建时间
	UpdateTime string // 更新时间
}

// companyColumns holds the columns for table t_company.
var companyColumns = CompanyColumns{
	ID:          "id",
	Name:        "name",
	Type:        "type",
	Country:     "country",
	City:        "city",
	Address:     "address",
	Email:       "email",
	Description: "description",
	Version:     "version",

	SocialCreditCode:      "social_credit_code",
	LegalPersonName:       "legal_person_name",
	LegalPersonCardNumber: "legal_person_card_number",

	CreateTime: "create_time",
	UpdateTime: "update_time",
}

// NewAsyncTaskDao creates and returns a new DAO object for table data access.
func NewCompanyDao() *CompanyDao {
	return &CompanyDao{
		group:   "default",
		table:   "t_company",
		columns: companyColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CompanyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CompanyDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CompanyDao) Columns() CompanyColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *CompanyDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CompanyDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CompanyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
