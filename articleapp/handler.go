package articleapp

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
	"github.com/leyle/ginbase/util"
)

type AddArticleForm struct {
	Title       string   `json:"title" binding:"required,min=4"`
	Content     string   `json:"content" binding:"required,min=4"`
	Description string   `json:"description"`
	Image       string   `json:"image" binding:"required"`
	CategoryID  string   `json:"categoryId" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

func CreateArticleHandler(ctx *apicontext.Context, c *gin.Context) {
	var form AddArticleForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)
	ds := ctx.Ds
	user := common.GetUserInfoByToken(c)
	article, err := CheckArticle(&form, ds)
	middleware.StopExec(err)
	article.UserId = user.ID
	article.CreateUser = user
	err = ds.Create(&article).Error
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, article)
	return
}

func UpdateArticleHandler(ctx *apicontext.Context, c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		returnfun.ReturnErrJson(c, "not found")
		return
	}
	var article common.Article
	ds := ctx.Ds
	err := ds.Debug().Where("id = ?", id).Preload("CreateUser").First(&article).Error

	middleware.StopExec(err)

	var form AddArticleForm
	err = c.BindJSON(&form)
	middleware.StopExec(err)
	articleUp, err := CheckArticle(&form, ds)
	middleware.StopExec(err)
	articleUp.CreateT = article.CreateT
	articleUp.DeleteT = article.DeleteT
	articleUp.ID = article.ID
	articleUp.UserId = article.UserId
	articleUp.CreateUser = article.CreateUser
	err = ds.Save(&articleUp).Error
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, articleUp)
	return

}

func GetArticleHandler(ctx *apicontext.Context, c *gin.Context) {
	var articles []*common.Article
	name := c.Query("title")
	page, size, skip := util.GetPageAndSize(c)
	var total int
	ds := ctx.Ds
	nice := ds.Debug().Model(&articles).Where("title LIKE ? AND delete_t = ?", "%"+name+"%", 0).Preload("Tags")
	err := nice.Preload("CreateUser").Preload("Category").Find(&articles).Count(&total).Limit(size).Offset(skip).Order("create_t desc").Error
	middleware.StopExec(err)
	retData := returnfun.QueryListData{Data: articles, Page: page, Size: size, Total: total}
	returnfun.ReturnOKJson(c, retData)
	return
}

func GetArticleInfoHandler(ctx *apicontext.Context, c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		returnfun.ReturnErrJson(c, "not found")
		return
	}
	var article common.Article
	ds := ctx.Ds
	err := ds.Debug().Where("id = ?", id).Preload("Tags").Preload("Category").Preload("CreateUser").First(&article).Error

	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, article)
	return
}
