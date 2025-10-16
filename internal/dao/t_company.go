package dao

import "ExhibitionService/internal/dao/internal"

// companyDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type companyDao struct {
	*internal.CompanyDao
}

var (
	// Company is globally public accessible object for table t_company operations.
	Company = companyDao{
		internal.NewCompanyDao(),
	}
)
