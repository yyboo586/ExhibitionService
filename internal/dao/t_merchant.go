package dao

import (
	"ExhibitionService/internal/dao/internal"
)

var Merchant = merchantDao{
	internal.NewTMerchant(),
}

type merchantDao struct {
	*internal.TMerchant
}
