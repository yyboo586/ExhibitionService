package model

import (
	"errors"
	"time"
)

var ErrConcurrentUpdate error = errors.New("concurrent update failed")

type AuthorRequired struct {
	Authorization string `p:"Authorization" in:"header" v:"required#Authorization不能为空" dc:"Bearer {{token}}"`
}

type AuthorNotRequired struct {
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
}

type PageReq struct {
	Page int `p:"page" d:"1" dc:"页码"`
	Size int `p:"size" d:"10" dc:"每页数量"`
}

type PageRes struct {
	Total       int `json:"total" dc:"总数"`
	CurrentPage int `json:"current_page" dc:"当前页码"`
}

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateTime)
}
