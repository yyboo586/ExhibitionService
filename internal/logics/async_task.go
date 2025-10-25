package logics

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	asyncTask "github.com/yyboo586/common/AsyncTask"
)

var (
	asyncTaskOnce    sync.Once
	asyncTaskManager asyncTask.Manager
)

func NewAsyncTask() asyncTask.Manager {
	asyncTaskOnce.Do(func() {
		var err error
		config := &asyncTask.Config{
			DSN:      g.Config().MustGet(context.Background(), "database.default.link").String(),
			Database: "exhibition_service",
		}
		asyncTaskManager, err = asyncTask.NewAsyncTaskManager(config)
		if err != nil {
			panic(err)
		}

	})
	return asyncTaskManager
}
