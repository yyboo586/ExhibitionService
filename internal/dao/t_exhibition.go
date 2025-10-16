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
