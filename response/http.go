package response

import (
	"github.com/easyhutu/kratosV2-handles/ecode"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func JSON(ctx http.Context, v interface{}, err error) error {
	return ctx.JSON(ecode.StatusCode200, NewResponse(v, err))
}

func NewResponse(data interface{}, err error) *Response {
	e, ok := err.(*ecode.ECode)
	if !ok {
		if e != nil {
			e = ecode.New(ecode.Code500, err.Error())
		} else {
			e = ecode.New(ecode.Code0, "")
		}

	}
	return &Response{
		Code:    e.Code,
		Data:    data,
		Message: e.Message,
	}
}
