package service

import (
	"ExhibitionService/api/v1/system"
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"context"
)

type fileService struct {
	fileDomain interfaces.IFile
	fileEngine interfaces.IFileEngine
}

func NewFileService(fileDomain interfaces.IFile, fileEngine interfaces.IFileEngine) *fileService {
	return &fileService{
		fileDomain: fileDomain,
		fileEngine: fileEngine,
	}
}

func (s *fileService) PreUploadFile(ctx context.Context, req *system.UploadFileReq) (*system.UploadFileRes, error) {
	out, err := s.fileEngine.PreUpload(ctx, &model.PreUploadReq{
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Size:        req.FileSize,
		BucketID:    "public-bucket",
	})
	if err != nil {
		return nil, err
	}

	err = s.fileDomain.Create(ctx, out.ID, req.FileName, out.VisitURL)
	if err != nil {
		return nil, err
	}

	return &system.UploadFileRes{FileID: out.ID, OriginalName: out.OriginalName, UploadURL: out.UploadURL}, nil
}

func (s *fileService) ReportFileUploadResult(ctx context.Context, req *system.ReportFileUploadResultReq) (*system.ReportFileUploadResultRes, error) {
	err := s.fileEngine.ReportUploadResult(ctx, req.FileID, req.Success)
	if err != nil {
		return nil, err
	}

	err = s.fileDomain.UpdateStatus(ctx, req.FileID, model.GetFileStatusFromBool(req.Success))
	if err != nil {
		return nil, err
	}

	return &system.ReportFileUploadResultRes{}, nil
}
