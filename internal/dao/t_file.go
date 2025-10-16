package dao

import "ExhibitionService/internal/dao/internal"

// fileDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type fileDao struct {
	*internal.FileDao
}

var (
	// File is globally public accessible object for table t_file operations.
	File = fileDao{
		internal.NewFileDao(),
	}
)
