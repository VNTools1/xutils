// Package xerror
/*
 * @Date: 2023-07-20 09:34:46
 * @LastEditTime: 2023-07-20 10:00:55
 * @Description:
 */
package xerror

import (
	"sort"
	"sync"
)

// 定义错误
var (
	ErrCaptchaWrong     = New(425, "error captcha")
	ErrBadParameter     = New(426, "error bad parameter")
	ErrInProcess        = New(427, "processing, please try again later")
	ErrLoginFailed      = New(401, "login failed, please try again")
	ErrTokenNotFound    = New(401, "please login first")
	ErrInvalidToken     = New(401, "login invalid, please login again")
	ErrTokenExpired     = New(401, "login expired, please login again")
	ErrUserTokenExpired = New(401, "user login expired, please login again")
	ErrNoPermmission    = New(401, "no access permission")
	ErrInvalidContent   = New(87014, "illegal content")
	ErrIsEmpty          = New(4041, "data does not exist")
	ErrUserIsEmpty      = New(4042, "user does not exist")
	ErrUserStatusWrong  = New(4044, "user logged out")
	ErrNotFound         = New(404, "resource not found")
	ErrMethodNotAllow   = New(405, "method not allowed")
	ErrAppConfigErr     = New(426, "service configuration error")
	ErrTooManyRequests  = New(429, "requests are too frequent")
	ErrInternalServer   = New(500, "server error")
)

var (
	errors = make(map[int]*Error)
	mu     sync.Mutex
)

// New ...
func New(code int, msg string) *Error {
	err := &Error{
		Code: code,
		Msg:  msg,
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := errors[code]; !ok {
		errors[code] = err
	}
	return err
}

// Errors 返回所有 error 列表
// 按 code 升序
// 通过 WithMsg 声明的错误码不重复显示
func Errors() []*Error {
	res := make([]*Error, 0)
	if len(errors) == 0 {
		return res
	}
	codes := make([]int, 0)
	for code := range errors {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for _, code := range codes {
		res = append(res, errors[code])
	}
	return res
}
