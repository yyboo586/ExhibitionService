package interfaces

import "context"

type HTTPClient interface {
	GET(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error)
	POST(ctx context.Context, url string, header map[string]interface{}, body interface{}) (status int, respBody []byte, err error)
	DELETE(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error)
}
