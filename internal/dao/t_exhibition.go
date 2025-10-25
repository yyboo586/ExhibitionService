package dao

import (
	"ExhibitionService/internal/dao/internal"
)

var Exhibition = exhibitionDao{
	internal.NewTExhibition(),
}

type exhibitionDao struct {
	*internal.TExhibition
}

var ExOrganizer = exOrganizerDao{
	internal.NewTExOrganizer(),
}

type exOrganizerDao struct {
	*internal.TExOrganizer
}

var (
	ExhibitionMerchant = exhibitionMerchantDao{
		internal.NewTExhibitionMerchant(),
	}
)

type exhibitionMerchantDao struct {
	*internal.TExhibitionMerchant
}
