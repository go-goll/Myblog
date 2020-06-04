package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
)

func UserRouter(ctx *apicontext.Context, g *gin.RouterGroup) {
	JR := g.Group("", func(c *gin.Context) {
	})
	{
		JR.POST("api/user/register", func(c *gin.Context) {
			RegisterHandler(ctx, c)
		})
		JR.POST("api/user/login", func(c *gin.Context) {
			LoginByPasswdHandler(ctx, c)
		})
	}

	auth := g.Group("", func(c *gin.Context) {
		common.Auth(c)
	})
	{
		auth.GET("api/user/info", func(c *gin.Context) {
			GetUserInfoHandler(ctx, c)
		})
		auth.POST("api/user/logout", func(c *gin.Context) {
			LogoutHandler(ctx, c)
		})
		auth.POST("api/user/changepassword", func(c *gin.Context) {
			UpdatePasswdHandler(ctx, c)
		})
	}
}
