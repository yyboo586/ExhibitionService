package model

import asyncTask "github.com/yyboo586/common/AsyncTask"

const (
	_                                    asyncTask.TaskType = iota
	TaskTypeExhibitionAutoStartEnrolling                    // 展会自动开始报名
	TaskTypeExhibitionAutoEndEnrolling                      // 展会自动结束报名
	TaskTypeExhibitionAutoStartRunning                      // 展会自动开始进行
	TaskTypeExhibitionAutoEnd                               // 展会自动结束
)
