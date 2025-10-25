package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type TExhibitionMerchant struct {
	table   string
	group   string // group is the database configuration group name of current DAO.
	columns TExhibitionMerchantColumns
}

type TExhibitionMerchantColumns struct {
	ID                  string // 主键
	ExhibitionID        string // 展会ID
	MerchantID          string // 商户ID
	Status              string // 状态
	CreateTime          string // 创建时间
	SubmitForReviewTime string // 提交审核时间
	ApproveTime         string // 审核通过时间
	UpdateTime          string // 更新时间
}

var exhibitionMerchantColumns = TExhibitionMerchantColumns{
	ID:                  "id",
	ExhibitionID:        "exhibition_id",
	MerchantID:          "merchant_id",
	Status:              "status",
	CreateTime:          "create_time",
	SubmitForReviewTime: "submit_for_review_time",
	ApproveTime:         "approve_time",
	UpdateTime:          "update_time",
}

func NewTExhibitionMerchant() *TExhibitionMerchant {
	return &TExhibitionMerchant{
		group:   "default",
		table:   "t_exhibition_merchant",
		columns: exhibitionMerchantColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TExhibitionMerchant) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TExhibitionMerchant) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TExhibitionMerchant) Columns() TExhibitionMerchantColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *TExhibitionMerchant) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TExhibitionMerchant) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TExhibitionMerchant) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
