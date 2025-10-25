package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/model"
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	asyncTask "github.com/yyboo586/common/AsyncTask"
)

type ExhibitionAction func(ctx context.Context, exhibition *model.Exhibition, data interface{}) error

type ExhibitionTransition struct {
	State  model.ExhibitionStatus
	Action ExhibitionAction
}

func (e *exhibition) initTransitionMap() {
	e.transitionMap = map[model.ExhibitionStatus]map[model.ExhibitionEvent]ExhibitionTransition{
		model.ExhibitionStatusPreparing: {
			model.ExhibitionEventSubmitForReview: {
				State:  model.ExhibitionStatusPending,
				Action: e.handleExhibitionSubmitForReview,
			},
			model.ExhibitionEventCancel: {
				State:  model.ExhibitionStatusCancelled,
				Action: e.handleExhibitionCancel,
			},
		},
		model.ExhibitionStatusPending: {
			model.ExhibitionEventApprove: {
				State:  model.ExhibitionStatusApproved,
				Action: e.handleExhibitionApprove,
			},
			model.ExhibitionEventReject: {
				State:  model.ExhibitionStatusPreparing,
				Action: e.handleExhibitionReject,
			},
		},
		model.ExhibitionStatusApproved: {
			model.ExhibitionEventStartEnrolling: {
				State:  model.ExhibitionStatusEnrolling,
				Action: e.handleExhibitionStartEnrolling,
			},
			model.ExhibitionEventCancel: {
				State:  model.ExhibitionStatusCancelled,
				Action: e.handleExhibitionCancel,
			},
		},
		model.ExhibitionStatusEnrolling: {
			model.ExhibitionEventEndEnrolling: {
				State:  model.ExhibitionStatusEnrollingEnded,
				Action: e.handleExhibitionEndEnrolling,
			},
			model.ExhibitionEventCancel: {
				State:  model.ExhibitionStatusCancelled,
				Action: e.handleExhibitionCancel,
			},
		},
		model.ExhibitionStatusEnrollingEnded: {
			model.ExhibitionEventStartRunning: {
				State:  model.ExhibitionStatusRunning,
				Action: e.handleExhibitionStartRunning,
			},
			model.ExhibitionEventCancel: {
				State:  model.ExhibitionStatusCancelled,
				Action: e.handleExhibitionCancel,
			},
		},
		model.ExhibitionStatusRunning: {
			model.ExhibitionEventEnd: {
				State:  model.ExhibitionStatusEnded,
				Action: e.handleExhibitionEnd,
			},
		},
	}
}

func (e *exhibition) HandleEvent(ctx context.Context, exhibitionID string, event model.ExhibitionEvent, data interface{}) (err error) {
	// 获取展会信息
	exhibition, err := e.GetExhibition(ctx, exhibitionID)
	if err != nil {
		return err
	}

	// 验证当前状态是否支持该事件
	transition, ok := e.transitionMap[exhibition.Status][event]
	if !ok {
		return gerror.Newf("exhibition current status: %s, not supported event: %s", model.GetExhibitionStatusText(exhibition.Status), model.GetExhibitionEventText(event))
	}

	// 执行状态操作
	err = transition.Action(ctx, exhibition, data)
	if err != nil {
		return gerror.Newf("handle event failed, exhibition id: %s, event: %s, err: %s", exhibitionID, model.GetExhibitionEventText(event), err.Error())
	}

	return nil
}

func (e *exhibition) handleExhibitionSubmitForReview(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusPending),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("submit for review failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("submit for review failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("submit for review failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (e *exhibition) handleExhibitionApprove(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	exInfo := map[string]any{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusApproved),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}
	taskContent := map[string]any{
		"exhibition_id":      exhibition.ID,
		"registration_start": exhibition.RegistrationStart.Unix(),
	}
	taskContentBytes, err := json.Marshal(taskContent)
	if err != nil {
		return gerror.Newf("approve exhibition failed, marshal task content failed, err: %s", err.Error())
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = e.asyncTask.AddScheduledTaskWithTx(ctx, tx, model.TaskTypeExhibitionAutoStartEnrolling, exhibition.ID, taskContentBytes, exhibition.RegistrationStart)
		if err != nil {
			return gerror.Newf("approve exhibition failed, add scheduled task failed, err: %s", err.Error())
		}

		result, err := dao.Exhibition.Ctx(ctx).Data(exInfo).
			Where(dao.Exhibition.Columns().ID, exhibition.ID).
			Where(dao.Exhibition.Columns().Version, exhibition.Version).
			Update()
		if err != nil {
			return gerror.Newf("approve exhibition failed, err: %s", err.Error())
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return gerror.Newf("approve exhibition failed, get rows affected failed, err: %s", err.Error())
		}
		if rowsAffected == 0 {
			return gerror.Newf("approve exhibition failed, %v", model.ErrConcurrentUpdate)
		}

		return nil
	})
	if err != nil {
		return gerror.Newf("approve exhibition failed, err: %s", err.Error())
	}

	e.asyncTask.WakeUp(model.TaskTypeExhibitionAutoStartEnrolling)

	return nil
}

func (e *exhibition) HandleTaskAutoStartEnrolling(ctx context.Context, task *asyncTask.Task) (err error) {
	return e.HandleEvent(ctx, task.CustomID, model.ExhibitionEventStartEnrolling, task)
}

func (e *exhibition) handleExhibitionReject(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusPreparing),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("reject exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("reject exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("reject exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (e *exhibition) handleExhibitionStartEnrolling(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	exInfo := map[string]any{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusEnrolling),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}
	taskContent := map[string]any{
		"exhibition_id":    exhibition.ID,
		"registration_end": exhibition.RegistrationEnd.Unix(),
	}
	taskContentBytes, err := json.Marshal(taskContent)
	if err != nil {
		return gerror.Newf("start enrolling failed, marshal task content failed, err: %s", err.Error())
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = e.asyncTask.AddScheduledTaskWithTx(ctx, tx, model.TaskTypeExhibitionAutoEndEnrolling, exhibition.ID, taskContentBytes, exhibition.RegistrationEnd)
		if err != nil {
			return gerror.Newf("start enrolling failed, add scheduled task failed, err: %s", err.Error())
		}

		result, err := dao.Exhibition.Ctx(ctx).Data(exInfo).
			Where(dao.Exhibition.Columns().ID, exhibition.ID).
			Where(dao.Exhibition.Columns().Version, exhibition.Version).
			Update()
		if err != nil {
			return gerror.Newf("start enrolling failed, err: %s", err.Error())
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return gerror.Newf("start enrolling failed, get rows affected failed, err: %s", err.Error())
		}
		if rowsAffected == 0 {
			return gerror.Newf("start enrolling failed, %v", model.ErrConcurrentUpdate)
		}

		return nil
	})
	if err != nil {
		return gerror.Newf("start enrolling failed, err: %s", err.Error())
	}

	e.asyncTask.WakeUp(model.TaskTypeExhibitionAutoEndEnrolling)

	return nil
}

func (e *exhibition) HandleTaskAutoEndEnrolling(ctx context.Context, task *asyncTask.Task) (err error) {
	return e.HandleEvent(ctx, task.CustomID, model.ExhibitionEventEndEnrolling, task)
}

func (e *exhibition) handleExhibitionEndEnrolling(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	exInfo := map[string]any{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusEnrollingEnded),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}
	taskContent := map[string]any{
		"exhibition_id": exhibition.ID,
		"start_time":    exhibition.StartTime.Unix(),
	}
	taskContentBytes, err := json.Marshal(taskContent)
	if err != nil {
		return gerror.Newf("end enrolling failed, marshal task content failed, err: %s", err.Error())
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = e.asyncTask.AddScheduledTaskWithTx(ctx, tx, model.TaskTypeExhibitionAutoStartRunning, exhibition.ID, taskContentBytes, exhibition.StartTime)
		if err != nil {
			return gerror.Newf("end enrolling failed, add scheduled task failed, err: %s", err.Error())
		}

		result, err := dao.Exhibition.Ctx(ctx).Data(exInfo).
			Where(dao.Exhibition.Columns().ID, exhibition.ID).
			Where(dao.Exhibition.Columns().Version, exhibition.Version).
			Update()
		if err != nil {
			return gerror.Newf("end enrolling failed, err: %s", err.Error())
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return gerror.Newf("end enrolling failed, get rows affected failed, err: %s", err.Error())
		}
		if rowsAffected == 0 {
			return gerror.Newf("end enrolling failed, %v", model.ErrConcurrentUpdate)
		}

		return nil
	})
	if err != nil {
		return gerror.Newf("end enrolling failed, err: %s", err.Error())
	}

	e.asyncTask.WakeUp(model.TaskTypeExhibitionAutoStartRunning)

	return nil
}

func (e *exhibition) HandleTaskAutoStartRunning(ctx context.Context, task *asyncTask.Task) (err error) {
	return e.HandleEvent(ctx, task.CustomID, model.ExhibitionEventStartRunning, task)
}

func (e *exhibition) handleExhibitionStartRunning(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	exInfo := map[string]any{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusRunning),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}
	taskContent := map[string]any{
		"exhibition_id": exhibition.ID,
		"end_time":      exhibition.EndTime.Unix(),
	}
	taskContentBytes, err := json.Marshal(taskContent)
	if err != nil {
		return gerror.Newf("start running failed, marshal task content failed, err: %s", err.Error())
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = e.asyncTask.AddScheduledTaskWithTx(ctx, tx, model.TaskTypeExhibitionAutoEnd, exhibition.ID, taskContentBytes, exhibition.EndTime)
		if err != nil {
			return gerror.Newf("start running failed, add scheduled task failed, err: %s", err.Error())
		}

		result, err := dao.Exhibition.Ctx(ctx).Data(exInfo).
			Where(dao.Exhibition.Columns().ID, exhibition.ID).
			Where(dao.Exhibition.Columns().Version, exhibition.Version).
			Update()
		if err != nil {
			return gerror.Newf("start running failed, err: %s", err.Error())
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return gerror.Newf("start running failed, get rows affected failed, err: %s", err.Error())
		}
		if rowsAffected == 0 {
			return gerror.Newf("start running failed, %v", model.ErrConcurrentUpdate)
		}

		return nil
	})
	if err != nil {
		return gerror.Newf("start running failed, err: %s", err.Error())
	}

	e.asyncTask.WakeUp(model.TaskTypeExhibitionAutoEnd)

	return nil
}

func (e *exhibition) HandleTaskAutoEnd(ctx context.Context, task *asyncTask.Task) (err error) {
	return e.HandleEvent(ctx, task.CustomID, model.ExhibitionEventEnd, task)
}

func (e *exhibition) handleExhibitionEnd(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusEnded),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("end exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("end exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("end exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}

func (e *exhibition) handleExhibitionCancel(ctx context.Context, exhibition *model.Exhibition, data interface{}) (err error) {
	result, err := dao.Exhibition.Ctx(ctx).Data(g.Map{
		dao.Exhibition.Columns().Status:     int(model.ExhibitionStatusCancelled),
		dao.Exhibition.Columns().Version:    exhibition.Version + 1,
		dao.Exhibition.Columns().UpdateTime: time.Now().Unix(),
	}).
		Where(dao.Exhibition.Columns().ID, exhibition.ID).
		Where(dao.Exhibition.Columns().Version, exhibition.Version).
		Update()
	if err != nil {
		return gerror.Newf("cancel exhibition failed, err: %s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return gerror.Newf("cancel exhibition failed, get rows affected failed, err: %s", err.Error())
	}
	if rowsAffected == 0 {
		return gerror.Newf("cancel exhibition failed, %v", model.ErrConcurrentUpdate)
	}

	return nil
}
