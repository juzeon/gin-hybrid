package dto

type CreateArticleReq struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}