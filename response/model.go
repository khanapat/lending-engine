package response

type Response struct {
	Code  uint64      `json:"code" example:"200"`
	Title string      `json:"title" example:"Register key success."`
	Data  interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Code        uint64 `json:"code" example:"400"`
	Title       string `json:"title" example:"Cannot register public key."`
	Description string `json:"description" example:"Please contact administrator for more information."`
	Error       string `json:"error,omitempty" example:"'<Field>' must be REQUIRED field but the input is ''"`
}

func NewResponse(resp Response, data interface{}) *Response {
	return &Response{
		Code:  resp.Code,
		Title: resp.Title,
		Data:  data,
	}
}

func NewErrResponse(resp ErrResponse, err string) *ErrResponse {
	return &ErrResponse{
		Code:        resp.Code,
		Title:       resp.Title,
		Description: resp.Description,
		Error:       err,
	}
}
