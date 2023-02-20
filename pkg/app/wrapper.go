package app

import (
	"gin-hybrid/data/dto"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	wrapper *Wrapper
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

func (w Wrapper) OK() Result {
	return Result{
		Code:    0,
		Msg:     "",
		Data:    nil,
		wrapper: &w,
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
	return Result{
		Code:    -1,
		Msg:     msg,
		Data:    nil,
		wrapper: &w,
	}
}
func (w Wrapper) ErrorWithCode(code int, msg string) Result {
	return Result{
		Code:    code,
		Msg:     msg,
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
