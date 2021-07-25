package client

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"lending-engine/common"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	contentType     string = "Content-Type"
	applicationJSON string = "application/json"
	applicationXML  string = "application/xml"
	textXML         string = "text/xml"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: http.DefaultMaxIdleConnsPerHost,
			MaxConnsPerHost:     0,
			IdleConnTimeout:     0,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
		Timeout: viper.GetDuration("client.timeout"),
	}
	return &Client{
		HTTPClient: httpClient,
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	if req.XRequestID == "" {
		req.XRequestID = uuid.New().String()
	}
	req.addHeader(contentType, applicationJSON)
	req.addHeader(common.XRequestID, req.XRequestID)

	body := string(req.Body)
	req.Logger.Debug(body)
	if req.HideLogRequestBody {
		body = ""
	}
	req.Logger.Info("http client request info",
		zap.String("method", req.Method),
		zap.String("url", req.URL),
		zap.Reflect("header", req.Header),
		zap.String("body", body))
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	for key, value := range req.Header {
		httpReq.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body = string(respBody)
	req.Logger.Debug(body)
	if req.HideLogResponseBody {
		body = ""
	}
	req.Logger.Info("http client response info",
		zap.String("latency", time.Since(start).String()),
		zap.String("status", resp.Status),
		zap.Reflect("header", resp.Header),
		zap.String("body", body),
		zap.String("url", resp.Request.URL.String()))

	return &Response{
		HTTPResponse: resp,
		Body:         respBody,
	}, nil
}
