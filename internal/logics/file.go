package logics

import (
	"ExhibitionService/internal/dao"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/model/entity"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
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
		return
	}
	return
}

func (f *file) UpdateFileStatus(ctx context.Context, fileID string, status model.FileStatus) (err error) {
	_, err = dao.File.Ctx(ctx).Data(map[string]any{
		dao.File.Columns().Status:     int(status),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.File.Columns().FileID, fileID).Update()
	return
}

func (f *file) UpdateFileCompanyInfo(ctx context.Context, tx gdb.TX, fileID string, companyID string, typ model.FileType) (err error) {
	_, err = dao.File.Ctx(ctx).TX(tx).Data(map[string]any{
		dao.File.Columns().CompanyID:  companyID,
		dao.File.Columns().Type:       int(typ),
		dao.File.Columns().UpdateTime: time.Now().Unix(),
	}).Where(dao.File.Columns().FileID, fileID).Update()
	return
}

func (f *file) GetFile(ctx context.Context, fileID string) (out *model.File, err error) {
	var t entity.TFile
	err = dao.File.Ctx(ctx).Where(dao.File.Columns().FileID, fileID).Scan(&t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("文件不存在: %s", fileID)
		}
		return nil, err
	}

	out = model.ConvertFile(&t)

	if out.Status != model.FileStatusUploadSuccess {
		return nil, fmt.Errorf("文件状态异常: %s", model.GetFileStatusText(out.Status))
	}

	return out, nil
}
