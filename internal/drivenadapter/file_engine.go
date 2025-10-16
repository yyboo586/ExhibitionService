package drivenadapter

import (
	"ExhibitionService/internal/interfaces"
	"ExhibitionService/internal/model"
	"ExhibitionService/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	fileEngineOnce     sync.Once
	fileEngineInstance *fileEngine
)

type fileEngine struct {
	addr   string
	client interfaces.HTTPClient
}

func NewFileEngine() interfaces.IFileEngine {
	fileEngineOnce.Do(func() {
		fileEngineInstance = &fileEngine{
			addr:   g.Cfg().MustGet(context.Background(), "server.thirdService.fileEngineService").String(),
			client: service.NewHTTPClient(),
		}
	})
	return fileEngineInstance
}

func (f *fileEngine) PreUpload(ctx context.Context, in *model.PreUploadReq) (out *model.PreUploadRes, err error) {
	url := fmt.Sprintf("%s/api/v1/file-engine/files/upload-tokens", f.addr)
	fileInfo := map[string]interface{}{
		"filename":     in.FileName,
		"content_type": in.ContentType,
		"size":         in.Size,
		"bucket_id":    in.BucketID,
	}
	reqBody := map[string]interface{}{
		"file": fileInfo,
	}

	status, respBody, err := f.client.POST(ctx, url, nil, reqBody)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf("get upload url failed, status: %d, respBody: %s", status, string(respBody))
		return
	}

	var resp map[string]interface{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	out = &model.PreUploadRes{
		ID:           resp["id"].(string),
		OriginalName: resp["original_name"].(string),
		VisitURL:     resp["visit_url"].(string),
		UploadURL:    resp["upload_url"].(string),
		ExpiresAt:    resp["expires_at"].(string),
		ExpiresIn:    int64(resp["expires_in"].(float64)),
	}
	return
}

func (f *fileEngine) PreDownload(ctx context.Context, fileID string) (out *model.PreDownloadRes, err error) {
	url := fmt.Sprintf("%s/api/v1/file-engine/files/%s/download-tokens", f.addr, fileID)

	status, respBody, err := f.client.GET(ctx, url, nil)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf("get download url failed, status: %d, respBody: %s", status, string(respBody))
		return
	}

	var resp map[string]interface{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	out = &model.PreDownloadRes{
		DownloadURL: resp["download_url"].(string),
		ExpiresAt:   resp["expires_at"].(string),
		ExpiresIn:   int64(resp["expires_in"].(float64)),
	}
	return
}

func (f *fileEngine) Delete(ctx context.Context, fileID string) (err error) {
	url := fmt.Sprintf("%s/api/v1/file-engine/files/%s", f.addr, fileID)

	status, respBody, err := f.client.DELETE(ctx, url, nil)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf("delete file failed, status: %d, respBody: %s", status, string(respBody))
		return
	}

	return
}

func (f *fileEngine) ReportUploadResult(ctx context.Context, fileID string, success bool) (err error) {
	url := fmt.Sprintf("%s/api/v1/file-engine/files/%s/status", f.addr, fileID)

	reqBody := map[string]interface{}{
		"success": success,
	}

	status, respBody, err := f.client.POST(ctx, url, nil, reqBody)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf("report upload result failed, status: %d, respBody: %s", status, string(respBody))
		return
	}

	return
}
