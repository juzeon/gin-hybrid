package router

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/xeonx/timeago"
	"html/template"
	"strings"
	"time"
)

func GetWebRouters() []WebRouter {
	routers := []WebRouter{
		{
			Name:  "index",
			Title: "Index",
		},
		{
			Name:  "user/login",
			Title: "Login",
		},
		{
			Name:    "user/me",
			Title:   "User Information",
			UseAPIs: AssemblePaths("/user/me"),
		},
	}
	return routers
}
func GetWebRoutersCommonAPIs() map[string]APIRouter {
	return map[string]APIRouter{
		"user": AssemblePaths("/user/me")[0],
	}
}
func GetWebRoutersFuncs() map[string]any {
	merged := map[string]any{}
	for key, item := range map[string]any(sprig.FuncMap()) {
		merged[key] = item
	}
	custom := map[string]any{
		"raw": func(str string) template.HTML {
			return template.HTML(str)
		},
		"concat": func(values ...any) string {
			v := ""
			for range values {
				v += "%v"
			}
			return fmt.Sprintf(v, values...)
		},
		"ago": func(value time.Time) string {
			return timeago.NoMax(timeago.Chinese).Format(value)
		},
	}
	for key, item := range custom {
		merged[key] = item
	}
	return merged
}
func AssemblePaths(paths ...string) []APIRouter {
	var routers []APIRouter
	for _, path := range paths {
		if !strings.HasPrefix(path, "/") {
			panic("path must start with /: " + path)
		}
		if !strings.HasPrefix(path, "/api") {
			path = "/api" + path
		}
		router, ok := PathAPIRouterMap[path]
		if !ok {
			panic("router path " + path + " not exist")
		}
		routers = append(routers, router)
	}
	return routers
}
