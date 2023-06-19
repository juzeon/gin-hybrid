package router

import (
	"gin-hybrid/pkg/app"
	"gin-hybrid/rest"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

type ProxyRouter struct {
	restService *rest.Service
}

func NewProxyRouter(restService *rest.Service) *ProxyRouter {
	return &ProxyRouter{restService: restService}
}
func (p *ProxyRouter) Handler(c *gin.Context) {
	aw := app.NewWrapper(c)
	endpoint, err := p.restService.GetEndpointRandomly()
	if err != nil {
		aw.ErrorWithCode(503, err.Error()).SendJSON()
		return
	}
	u, err := url.Parse("http://" + endpoint)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func RegisterReverseProxy(restService *rest.Service, g *gin.RouterGroup) {
	proxyRouter := NewProxyRouter(restService)
	g.Any("/*all", proxyRouter.Handler)
}
