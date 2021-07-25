package client

import (
	"net/http"

	"go.uber.org/zap"
)

type Header map[string]string

type Request struct {
	URL                 string
	Method              string
	XRequestID          string
	Header              Header
	HideLogRequestBody  bool
	HideLogResponseBody bool
	Logger              *zap.Logger
	Body                []byte
}

func (r *Request) addHeader(key, value string) {
	if r.Header == nil {
		r.Header = Header{}
	}
	r.Header[key] = value
}

type Response struct {
	HTTPResponse *http.Response
	Body         []byte
}
