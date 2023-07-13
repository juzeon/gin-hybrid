package router

import (
	"fmt"
	"gin-hybrid/conf"
	"gin-hybrid/data/dto"
	"gin-hybrid/pkg/app"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type APIRouter struct {
	Method   string
	Path     string
	Handlers []func(aw *app.Wrapper) app.Result
	RPCOnly  bool
}

var PathAPIRouterMap = map[string]APIRouter{}

type WebRouter struct {
	Name           string               // name of router
	OverwritePath  string               // use this to rewrite relativePath if it's not null
	UseAPIs        []APIRouter          // APIRouters to call
	Process        func(map[string]any) // additionally process renderMap
	Title          string
	GetTitle       func(map[string]any) string // use GetTitle instead of Title if this function exists
	GetKeywords    func(map[string]any) string
	GetDescription func(map[string]any) string
}

func RegisterAPIRouters[T any](apiRouters []APIRouter, api *gin.RouterGroup, conf *conf.ServiceConfig[T]) {
	if !strings.HasPrefix(api.BasePath(), "/") {
		panic("BasePath must start with /: " + api.BasePath())
	}
	g := api.Group("/" + conf.InitConf.Name)
	for _, apiRouter := range apiRouters {
		apiRouter := apiRouter
		if !strings.HasPrefix(apiRouter.Path, "/") {
			panic("Path must start with /: " + apiRouter.Path)
		}
		apiRouter.Method = strings.ToLower(apiRouter.Method)
		commonHandler := func(ctx *gin.Context) {
			aw := app.NewWrapper(ctx)
			var result app.Result
			if apiRouter.RPCOnly {
				rpcKey := ctx.GetHeader("X-RPC-Key")
				if rpcKey != conf.ParentConf.RPCKey {
					ctx.JSON(401, "direct API Call sent to RPC-only routes")
					return
				}
			}
			t := time.Now()
			for _, handler := range apiRouter.Handlers {
				result = handler(aw)
				if !result.IsSuccessful() {
					break
				}
			}
			result.Duration = time.Now().Sub(t)
			if result.Reader != nil {
				ctx.DataFromReader(200, result.ResponseContentLength, result.ResponseContentType, result.Reader,
					result.ExtraHeaders)
			} else {
				ctx.JSON(result.GetResponseCode(), result)
			}
		}
		switch apiRouter.Method {
		case "get":
			g.GET(apiRouter.Path, commonHandler)
		case "post":
			g.POST(apiRouter.Path, commonHandler)
		case "put":
			g.PUT(apiRouter.Path, commonHandler)
		case "delete":
			g.DELETE(apiRouter.Path, commonHandler)
		case "patch":
			g.PATCH(apiRouter.Path, commonHandler)
		case "head":
			g.HEAD(apiRouter.Path, commonHandler)
		case "options":
			g.OPTIONS(apiRouter.Path, commonHandler)
		default:
			panic("method " + apiRouter.Method + " not found")
		}
		PathAPIRouterMap[apiRouter.Method+":"+g.BasePath()+apiRouter.Path] = apiRouter
	}
}
func RegisterWebRouters(webRouters []WebRouter, e *gin.Engine) {
	for _, webRouter := range webRouters {
		webRouter := webRouter
		if strings.HasPrefix(webRouter.Name, "/") {
			panic("WebRouter.Name should not start with /: " + webRouter.Name)
		}
		webRouterNameArr := strings.Split(webRouter.Name, "/")
		templateName := webRouterNameArr[len(webRouterNameArr)-1]
		webRouterHandler := func(ctx *gin.Context) {
			t := time.Now()
			renderMap := map[string]any{}
			var result app.Result
			var firstAPIResult app.Result
			aw := app.NewWrapper(ctx)

			// call APIs specified by templates
			for ix, apiRouter := range webRouter.UseAPIs {
				for _, apiHandler := range apiRouter.Handlers {
					result = apiHandler(aw)
					if !result.IsSuccessful() {
						break
					}
				}
				if !result.IsSuccessful() {
					ctx.HTML(result.GetResponseCode(), "error.gohtml", result)
					return
				}
				if ix == 0 {
					firstAPIResult = result
					renderMap["d"] = result.Data
					renderMap["code"] = result.Code
					renderMap["msg"] = result.Msg
				} else {
					renderMap["d"+strconv.Itoa(ix)] = result.Data
					renderMap["code"+strconv.Itoa(ix)] = result.Code
					renderMap["msg"+strconv.Itoa(ix)] = result.Msg
				}
			}

			// call common APIs
			for name, apiRouter := range GetWebRoutersCommonAPIs() {
				for _, apiHandler := range apiRouter.Handlers {
					result = apiHandler(aw)
					if !result.IsSuccessful() {
						break
					}
				}
				renderMap[name] = result.Data
			}
			if uc, exist := aw.Ctx.Get("userClaims"); exist {
				renderMap["role"] = uc.(*dto.UserClaims).RoleID
			} else {
				renderMap["role"] = 0
			}

			// modify renderMap with custom functions
			renderMap["title"] = webRouter.Title
			if webRouter.GetTitle != nil {
				renderMap["title"] = webRouter.GetTitle(renderMap)
			}
			if webRouter.GetDescription != nil {
				renderMap["description"] = webRouter.GetDescription(renderMap)
			}
			if webRouter.GetKeywords != nil {
				renderMap["keywords"] = webRouter.GetKeywords(renderMap)
			}
			if webRouter.Process != nil {
				webRouter.Process(renderMap)
			}

			renderMap["duration"] = time.Now().Sub(t)
			if firstAPIResult.Redirect.Code != 0 {
				ctx.Redirect(firstAPIResult.Redirect.Code, firstAPIResult.Redirect.Location)
			} else {
				ctx.HTML(200, templateName+".gohtml", renderMap)
			}
		}
		relativePath := "/" + webRouter.Name
		if webRouter.OverwritePath != "" {
			relativePath = webRouter.OverwritePath
		}
		e.GET(relativePath, webRouterHandler)
		if webRouter.Name == "index" {
			e.GET("/", webRouterHandler)
		}
	}
}

func Setup(e *gin.Engine, enableWebRouters bool) {
	e.Use(func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") != "" {
			return
		}
		if token, err := ctx.Cookie("hybrid_authorization"); err == nil {
			ctx.Request.Header.Set("Authorization", token)
		}
	})
	if enableWebRouters {
		e.HTMLRender = loadTemplates()
		e.Static("/static", "web/static")
		RegisterWebRouters(GetWebRouters(), e)
	}
}
func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	bases, err := filepath.Glob("web/template/base/*.gohtml")
	if err != nil {
		panic(err)
	}
	heads, err := filepath.Glob("web/template/head/*.gohtml")
	if err != nil {
		panic(err)
	}
	pages, err := filepath.Glob("web/template/page/**/*.gohtml")
	if err != nil {
		panic(err)
	}
	pagesMore, err := filepath.Glob("web/template/page/*.gohtml")
	if err != nil {
		panic(err)
	}
	pages = append(pages, pagesMore...)
	for _, include := range pages {
		baseCopy := make([]string, len(bases))
		copy(baseCopy, bases)
		headCopy := make([]string, len(heads))
		copy(headCopy, heads)
		files := append(baseCopy, include)
		files = append(files, headCopy...)
		r.AddFromFilesFuncs(filepath.Base(include), GetWebRoutersFuncs(), files...)
	}
	standaloneArr := []string{"error.gohtml"}
	for _, standalone := range standaloneArr {
		r.AddFromFilesFuncs(standalone, GetWebRoutersFuncs(), "web/template/standalone/"+standalone)
	}
	return r
}
func AssembleHandlers(handlers ...func(aw *app.Wrapper) app.Result) []func(aw *app.Wrapper) app.Result {
	var result []func(aw *app.Wrapper) app.Result
	for _, handler := range handlers {
		result = append(result, handler)
	}
	return result
}
func RecoveryFunc(c *gin.Context, err any) {
	aw := app.NewWrapper(c)
	msg := "Internal error: " + fmt.Sprintf("%v", err)
	if strings.HasPrefix(aw.Ctx.Request.RequestURI, "/api") {
		aw.Ctx.Header("Content-Type", "application/json; charset=utf-8")
		aw.ErrorWithCode(500, msg).SendJSON()
	} else {
		c.HTML(500, "error.gohtml", aw.Error(msg))
	}
}
