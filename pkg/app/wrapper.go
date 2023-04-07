package app

import (
	"gin-hybrid/data/dto"
	"github.com/gin-gonic/gin"
	"runtime"
	"time"
)

type Result struct {
	Code                int           `json:"code"`
	Msg                 string        `json:"msg,omitempty"`
	Line                int           `json:"line,omitempty"`
	File                string        `json:"file,omitempty"`
	Data                interface{}   `json:"data,omitempty"`
	Duration            time.Duration `json:"duration,omitempty"`
	wrapper             *Wrapper
	ResponseContentType string   `json:"-"`
	Redirect            redirect `json:"-"`
}
type redirect struct {
	Code     int
	Location string
}

func (r Result) SendJSON() {
	r.wrapper.Ctx.JSON(200, r)
}
func (r Result) IsSuccessful() bool {
	return r.Code == 0
}
func (r Result) ScanData(data any) {
	data = r.Data
}
func (r Result) GetResponseCode() int {
	if r.Code != 0 && r.Code != -1 {
		return r.Code
	}
	return 200
}

type Wrapper struct {
	Ctx *gin.Context
}

func NewWrapper(c *gin.Context) *Wrapper {
	return &Wrapper{Ctx: c}
}

func (w Wrapper) Redirect(url string, isPermanent bool) Result {
	code := 302
	if isPermanent {
		code = 301
	}
	return Result{Code: 0, Redirect: redirect{
		Code:     code,
		Location: url,
	}}
}
func (w Wrapper) OK() Result {
	return Result{
		Code:    0,
		Msg:     "",
		Data:    nil,
		wrapper: &w,
	}
}
func (w Wrapper) SuccessWithRawData(data []byte, contentType string) Result {
	return Result{
		Code:                0,
		Msg:                 "",
		Data:                data,
		wrapper:             &w,
		ResponseContentType: contentType,
	}
}
func (w Wrapper) Success(data interface{}) Result {
	return Result{
		Code:    0,
		Msg:     "",
		Data:    data,
		wrapper: &w,
	}
}
func (w Wrapper) Error(msg string) Result {
	_, file, n, _ := runtime.Caller(1)
	return Result{
		Code:    -1,
		Msg:     msg,
		Line:    n,
		File:    file,
		Data:    nil,
		wrapper: &w,
	}
}
func (w Wrapper) ErrorWithCode(code int, msg string) Result {
	_, file, n, _ := runtime.Caller(1)
	return Result{
		Code:    code,
		Msg:     msg,
		Line:    n,
		File:    file,
		Data:    nil,
		wrapper: &w,
	}
}
func (w Wrapper) GetIP() string {
	return w.Ctx.ClientIP()
}
func (w Wrapper) ExtractUserClaims() *dto.UserClaims {
	raw, exist := w.Ctx.Get("userClaims")
	if !exist {
		panic("userClaims not exists")
	}
	uc, ok := raw.(*dto.UserClaims)
	if !ok {
		panic("userClaims failed to convert")
	}
	return uc
}