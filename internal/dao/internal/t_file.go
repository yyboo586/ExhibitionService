package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FileDao is the data access object for table t_file.
type FileDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns FileColumns // columns contains all the column names of Table for convenient usage.
}

// FileColumns defines and stores column names for table t_file.
type FileColumns struct {
	ID         string // 文件ID
	Module     string // 模块类型
	CustomID   string // 自定义ID
	Type       string // 文件类型
	FileID     string // 文件ID
	FileName   string // 文件名称
	FileLink   string // 媒体URL
	Status     string // 文件上传状态
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
}

// fileColumns holds the columns for table t_file.
var fileColumns = FileColumns{
	ID:         "id",
	Module:     "module",
	CustomID:   "custom_id",
	Type:       "type",
	FileID:     "file_id",
	FileName:   "file_name",
	FileLink:   "file_link",
	Status:     "status",
	CreateTime: "create_time",
	UpdateTime: "update_time",
}

// NewFileDao creates and returns a new DAO object for table data access.
func NewFileDao() *FileDao {
	return &FileDao{
		group:   "default",
		table:   "t_file",
		columns: fileColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FileDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FileDao) Columns() FileColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current DAO.
func (dao *FileDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FileDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
