package tagapp

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
)

func TagRouter(ctx *apicontext.Context, g *gin.RouterGroup) {
	JR := g.Group("", func(c *gin.Context) {
	})
	{
		JR.GET("api/tag", func(c *gin.Context) {
			GetTagsHandler(ctx, c)
		})
		JR.GET("api/category", func(c *gin.Context) {
			GetCateHandler(ctx, c)
		})
	}
	auth := g.Group("", func(c *gin.Context) {
		common.Auth(c)
	})
	{
		auth.POST("api/tag", func(c *gin.Context) {
			CreateTagsHandler(ctx, c)
		})
		auth.POST("api/category", func(c *gin.Context) {
			CreateCateHandler(ctx, c)
		})

		auth.PUT("api/tag/:id", func(c *gin.Context) {
			UpdateTagsHandler(ctx, c)
		})
		auth.PUT("api/category/:id", func(c *gin.Context) {
			UpdateCateHandler(ctx, c)
		})
	}
}
