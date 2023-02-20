package router

import (
	"gin-hybrid/middleware"
	"gin-hybrid/service"
)

func GetUserAPIRouters() []APIRouter {
	srv := service.ExUser
	routers := []APIRouter{
		{
			Method:   "post",
			Path:     "/login",
			Handlers: AssembleHandlers(srv.Login),
		},
		{
			Method:   "get",
			Path:     "/me",
			Handlers: AssembleHandlers(middleware.Auth, srv.Me),
		},
		{
			Method:   "get",
			Path:     "/get_data",
			Handlers: AssembleHandlers(srv.GetData),
		},
	}
	return routers
}
