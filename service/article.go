package service

import (
	"fmt"
	"gin-hybrid/data/dto"
	"gin-hybrid/pkg/app"
	"gin-hybrid/rest"
	"time"
)

type ArticleService struct {
	userService *rest.Service
}

func NewArticleService(userService *rest.Service) *ArticleService {
	return &ArticleService{userService: userService}
}
func (a ArticleService) PostArticle(aw *app.Wrapper) app.Result {
	uc := aw.ExtractUserClaims()
	var req dto.CreateArticleReq
	if err := aw.Ctx.ShouldBind(&req); err != nil {
		return aw.Error(err.Error())
	}
	var exampleGet string
	a.userService.MustCall(&exampleGet, "get", "/example", map[string]any{"example": uc.UserID},
		aw.ExtractJWT())
	type ExamplePostData struct {
		Example time.Time `form:"example"`
	}
	var examplePost string
	a.userService.MustCall(&examplePost, "post", "/example", ExamplePostData{Example: uc.LoginTime},
		aw.ExtractJWT())
	return aw.Success(fmt.Sprintf("post an article with title %v and content %v, example_get: %v, example_post: %v",
		req.Title, req.Content, exampleGet, examplePost))
}
func (a ArticleService) ListArticle(aw *app.Wrapper) app.Result {
	return aw.Success("result of listed articles")
}
