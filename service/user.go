package service

import (
	"gin-hybrid/data/dto"
	"gin-hybrid/pkg/app"
	"gin-hybrid/pkg/util"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
func (u UserService) Login(aw *app.Wrapper) app.Result {
	var req dto.UserLoginReq
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
func (u UserService) ExampleGet(aw *app.Wrapper) app.Result {
	return aw.Success("This is an example GET call: " + aw.Ctx.Query("example"))
}
func (u UserService) ExamplePost(aw *app.Wrapper) app.Result {
	return aw.Success("This is an example POST call: " + aw.Ctx.PostForm("example"))
}
