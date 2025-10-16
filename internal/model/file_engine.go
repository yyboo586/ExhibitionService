package model

type PreUploadReq struct {
	FileName    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	BucketID    string `json:"bucket_id"`
}

type PreUploadRes struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	VisitURL     string `json:"visit_url"`
	UploadURL    string `json:"upload_url"`
	ExpiresAt    string `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
}

type PreDownloadRes struct {
	DownloadURL string `json:"download_url"`
	ExpiresAt   string `json:"expires_at"`
	ExpiresIn   int64  `json:"expires_in"`
}
