package ecode

const (
	StatusCode200 = 200
	StatusCode400 = 400
	StatusCode404 = 404
	Code500       = 500
	Code0         = 0
)

type ECode struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *ECode) Error() string {
	return e.Message
}

func New(Code int, Message string) *ECode {
	return &ECode{
		StatusCode: StatusCode200,
		Code:       Code,
		Message:    Message,
	}
}
