package errs

import "fmt"

// Error 错误
type Error struct {
	Code int
	Msg  string
}

// Error 错误信息
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// New 创建
func New(code int, msg string) error {
	err := &Error{
		Code: code,
		Msg:  msg,
	}

	return err
}

// Newf 格式化创建
func Newf(code int, format string, params ...interface{}) error {
	err := &Error{
		Code: int(int32(code)),
		Msg:  fmt.Sprintf(format, params...),
	}
	return err
}

const (
	CodeSuccess = 0
	CodeUnknown = 999
)

// Code 获取错误码
func Code(err error) int {
	if err == nil {
		return CodeSuccess
	}
	e, ok := err.(*Error)
	if !ok {
		return CodeUnknown
	}
	if e == (*Error)(nil) {
		return CodeSuccess
	}
	return e.Code
}

// Msg 获取错误信息
func Msg(err error) string {
	if err == nil {
		return "success"
	}
	e, ok := err.(*Error)
	if !ok {
		return err.Error()
	}
	if e == (*Error)(nil) {
		return "success"
	}
	return e.Msg
}
