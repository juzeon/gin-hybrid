package service

import (
	"gin-hybrid/pkg/app"
	"gin-hybrid/pkg/util"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
func (u UserService) Login(aw *app.Wrapper) app.Result {
	type UserLoginReq struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	var req UserLoginReq
	if err := aw.Ctx.ShouldBind(&req); err != nil {
		return aw.Error(err.Error())
	}
	if req.Username != "admin" || req.Password != "123456" {
		return aw.Error("Wrong username or password (tips: admin, 123456)")
	}
	jwt := util.GenerateJWT(1, 5, "administrator")
	aw.Ctx.SetCookie("hybrid_authorization", jwt, 60*60*24*365, "/", "", false, true)
	return aw.Success(jwt)
}
func (u UserService) Me(aw *app.Wrapper) app.Result {
	uc := aw.ExtractUserClaims()
	return aw.Success(uc)
}
func (u UserService) GetData(aw *app.Wrapper) app.Result {
	return aw.Success("This is an example API call.")
}
