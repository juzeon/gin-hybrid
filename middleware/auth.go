package middleware

import (
	"gin-hybrid/pkg/app"
	"gin-hybrid/pkg/util"
	"strings"
)

func Auth(aw *app.Wrapper) app.Result {
	authHeader := aw.Ctx.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = authHeader[7:]
	}
	claims, err := util.ParseJWT(authHeader)
	if err != nil {
		return aw.Error("Login Required")
	}
	aw.Ctx.Set("userClaims", claims)
	aw.Ctx.Set("jwt", authHeader)
	return aw.OK()
}
