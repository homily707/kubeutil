package error

import "fmt"

var (
	Success     = NewError(200, "success")
	ServerError = NewError(400, "server error")
)

type Error struct {
	code    int      `json:"code"`
	message string   `json:"message"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("error code %d exist already", code))
	}
	codes[code] = msg
	return &Error{code: code, message: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.message
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.message, args...)
}

func (e *Error) Details() []string {
	return e.details
}
