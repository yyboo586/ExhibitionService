package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	fileOnce   sync.Once
	fileDomain *file
)

type file struct{}

func NewFile() interfaces.IFile {
	fileOnce.Do(func() {
		fileDomain = &file{}
	})
	return fileDomain
}

func (f *file) Create(ctx context.Context, fileID string, fileName string, fileLink string) (err error) {
	data := map[string]any{
		dao.File.Columns().FileID:     fileID,
		dao.File.Columns().FileName:   fileName,
		dao.File.Columns().FileLink:   fileLink,
		dao.File.Columns().Status:     model.FileStatusInit,
		dao.File.Columns().CreateTime: time.Now().Unix(),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}
	_, err = dao.File.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return gerror.Newf("create file failed, err: %s", err.Error())
	}
	return nil
}

func (f *file) UpdateStatus(ctx context.Context, fileID string, status model.FileStatus) (err error) {
	_, err = dao.File.Ctx(ctx).Data(map[string]any{
		dao.File.Columns().Status:     int(status),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.File.Columns().FileID, fileID).Update()
	if err != nil {
		return gerror.Newf("update file status failed, err: %s", err.Error())
	}

	return nil
}

func (f *file) UpdateCustomInfo(ctx context.Context, tx gdb.TX, fileID string, fileModule model.FileModule, customID string, typ model.FileType) (err error) {
	_, err = f.Get(ctx, fileID)
	if err != nil {
		return err
	}

	_, err = dao.File.Ctx(ctx).TX(tx).Data(map[string]any{
		dao.File.Columns().Module:     int(fileModule),
		dao.File.Columns().CustomID:   customID,
		dao.File.Columns().Type:       int(typ),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.File.Columns().FileID, fileID).Update()
	if err != nil {
		return gerror.Newf("update file custom info failed, err: %s", err.Error())
	}

	return nil
}

func (f *file) Get(ctx context.Context, fileID string) (out *model.File, err error) {
	var t entity.TFile
	err = dao.File.Ctx(ctx).Where(dao.File.Columns().FileID, fileID).Scan(&t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.Newf("get file failed, file not found, file id: %s", fileID)
		}
		return nil, gerror.Newf("get file failed, err: %s", err.Error())
	}

	out = model.ConvertFile(&t)

	if out.Status != model.FileStatusUploadSuccess {
		return nil, gerror.Newf("get file failed, file status abnormal, file id: %s, status: %s", fileID, model.GetFileStatusText(out.Status))
	}

	return out, nil
}

func (f *file) ListByModuleAndCustomID(ctx context.Context, module model.FileModule, customID string) (out []*model.File, err error) {
	var t []*entity.TFile
	err = dao.File.Ctx(ctx).
		Where(dao.File.Columns().Module, int(module)).
		Where(dao.File.Columns().CustomID, customID).
		Where(dao.File.Columns().Status, model.FileStatusUploadSuccess).
		Scan(&t)
	if err != nil {
		return nil, gerror.Newf("list files by module and custom id failed, err: %s", err.Error())
	}

	for _, v := range t {
		out = append(out, model.ConvertFile(v))
	}
	return out, nil
}

func (f *file) ClearCustomInfo(ctx context.Context, tx gdb.TX, module model.FileModule, customID string) (err error) {
	_, err = dao.File.Ctx(ctx).TX(tx).Data(map[string]any{
		dao.File.Columns().Module:     int(module),
		dao.File.Columns().CustomID:   customID,
		dao.File.Columns().Type:       int(model.FileTypeUnknown),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.File.Columns().Module, int(module)).Where(dao.File.Columns().CustomID, customID).Update()
	if err != nil {
		return gerror.Newf("clear file custom info failed, err: %s", err.Error())
	}
	return nil
}

func (f *file) IsUploadSuccess(ctx context.Context, fileInfo *model.File) (err error) {
	if fileInfo.Status != model.FileStatusUploadSuccess {
		return gerror.Newf("check file upload success failed, file status abnormal, file id: %s, status: %s", fileInfo.FileID, model.GetFileStatusText(fileInfo.Status))
	}
	return nil
}

func (f *file) CheckFileUploadSuccess(ctx context.Context, fileIDs []string) (err error) {
	for _, v := range fileIDs {
		exists, err := dao.File.Ctx(ctx).
			Where(dao.File.Columns().FileID, v).
			Where(dao.File.Columns().Status, model.FileStatusUploadSuccess).
			Exist()
		if err != nil {
			return err
		}
		if !exists {
			return gerror.Newf("check file upload success failed, file not upload success, file id: %s", v)
		}
	}
	return nil
}
