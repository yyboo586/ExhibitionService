package dao

import (
	"ExhibitionService/internal/dao/internal"
)

var ServiceProvider = serviceProviderDao{
	internal.NewTServiceProvider(),
}

type serviceProviderDao struct {
	*internal.TServiceProvider
}
