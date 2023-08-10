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
	ErrCaptchaWrong     = New(425, "验证码错误")
	ErrBadParameter     = New(426, "请求参数错误")
	ErrInProcess        = New(427, "正在处理中，请稍后再试")
	ErrLoginFailed      = New(401, "登陆失败，请重试")
	ErrTokenNotFound    = New(401, "请先登陆")
	ErrInvalidToken     = New(401, "登陆失效，请重新登陆")
	ErrTokenExpired     = New(401, "登陆过期，请重新登陆")
	ErrUserTokenExpired = New(401, "用户登陆过期，请重新登陆")
	ErrNoPerm           = New(401, "无访问权限")
	ErrInvalidContent   = New(87014, "非法内容")
	ErrIsEmpty          = New(4041, "数据不存在")
	ErrUserIsEmpty      = New(4042, "用户不存在")
	ErrUserStatusWrong  = New(4044, "用户已注销")
	ErrNotFound         = New(404, "资源不存在")
	ErrMethodNotAllow   = New(405, "方法不被允许")
	ErrAppConfigErr     = New(426, "服务配置错误")
	ErrTooManyRequests  = New(429, "请求过于频繁")
	ErrInternalServer   = New(500, "服务器发生错误")
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
