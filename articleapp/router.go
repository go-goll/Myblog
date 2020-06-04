package articleapp

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
)

func ArticleRouter(ctx *apicontext.Context, g *gin.RouterGroup) {
	JR := g.Group("", func(c *gin.Context) {
	})
	{
		JR.GET("api/article", func(c *gin.Context) {
			GetArticleHandler(ctx, c)
		})
		JR.GET("api/article/:id", func(c *gin.Context) {
			GetArticleInfoHandler(ctx, c)
		})
	}
	auth := g.Group("", func(c *gin.Context) {
		common.Auth(c)
	})
	{
		auth.POST("api/article", func(c *gin.Context) {
			CreateArticleHandler(ctx, c)
		})
		auth.PUT("api/article/:id", func(c *gin.Context) {
			UpdateArticleHandler(ctx, c)
		})
	}
}
