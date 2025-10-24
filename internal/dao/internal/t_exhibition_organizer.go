package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type TExOrganizer struct {
	table   string
	group   string // group is the database configuration group name of current DAO.
	columns TExOrganizerColumns
}

type TExOrganizerColumns struct {
	ID                string // 展会ID
	ExhibitionID      string // 展会ID
	ServiceProviderID string // 服务提供商ID
	RoleType          string // 角色类型
	CreateTime        string
	UpdateTime        string
}

var exOrganizerColumns = TExOrganizerColumns{
	ID:                "id",
	ExhibitionID:      "exhibition_id",
	ServiceProviderID: "service_provider_id",
	RoleType:          "role_type",
	CreateTime:        "create_time",
	UpdateTime:        "update_time",
}

func NewTExOrganizer() *TExOrganizer {
	return &TExOrganizer{
		group:   "default",
		table:   "t_exhibition_organizer",
		columns: exOrganizerColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TExOrganizer) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TExOrganizer) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TExOrganizer) Columns() TExOrganizerColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *TExOrganizer) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TExOrganizer) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TExOrganizer) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
