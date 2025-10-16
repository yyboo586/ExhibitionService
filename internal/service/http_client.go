package service

import (
	"ExhibitionService/internal/interfaces"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type httpClient struct {
	c *http.Client
}

var (
	cOnce sync.Once
	c     *httpClient
)

func NewHTTPClient() interfaces.HTTPClient {
	cOnce.Do(func() {
		c = &httpClient{
			c: &http.Client{
				Timeout: time.Second * 10,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return nil
				},
			},
		}
	})
	return c
}

func (c *httpClient) GET(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	c.setHeader(req, header)

	return c.do(req)
}

func (c *httpClient) POST(ctx context.Context, url string, header map[string]interface{}, body interface{}) (status int, respBody []byte, err error) {
	var dataByte []byte
	switch data := body.(type) {
	case []byte:
		dataByte = data
	case string:
		dataByte = []byte(data)
	default:
		dataByte, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	g.Log().Debug(ctx, "------------------------------------------------------------------------------------------------")
	g.Log().Debugf(ctx, "Request URL: %s", url)
	g.Log().Debugf(ctx, "Request Header: %+v", header)
	g.Log().Debugf(ctx, "Request Body: %s", string(dataByte))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(dataByte))
	if err != nil {
		return
	}
	c.setHeader(req, header)
	status, respBody, err = c.do(req)
	if err != nil {
		return
	}

	g.Log().Debugf(ctx, "Response Status: %d", status)
	g.Log().Debugf(ctx, "Response Header: %+v", req.Header)
	g.Log().Debugf(ctx, "Response Body: %s", string(respBody))
	g.Log().Debug(ctx, "------------------------------------------------------------------------------------------------")
	return
}

func (c *httpClient) DELETE(ctx context.Context, url string, header map[string]interface{}) (status int, respBody []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	c.setHeader(req, header)
	return c.do(req)
}

func (c *httpClient) do(req *http.Request) (status int, respBody []byte, err error) {
	resp, err := c.c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return resp.StatusCode, respBody, nil
}

func (c *httpClient) setHeader(req *http.Request, header map[string]interface{}) {
	for k, v := range header {
		req.Header.Set(k, v.(string))
	}
}
