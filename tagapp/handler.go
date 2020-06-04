package tagapp

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
	"github.com/leyle/ginbase/util"
)

func GetTagsHandler(ctx *apicontext.Context, c *gin.Context) {
	var tags []*common.Tag
	ds := ctx.Ds

	name := c.Query("name")
	page, size, skip := util.GetPageAndSize(c)
	var total int
	err := ds.Where("name LIKE ? AND delete_t = ?", "%"+name+"%", 0).Find(&tags).Count(&total).Limit(size).Offset(skip).Order("create_t desc").Error
	middleware.StopExec(err)
	retData := returnfun.QueryListData{Data: tags, Page: page, Size: size, Total: total}
	returnfun.ReturnOKJson(c, retData)
	return
}

type AddTagForm struct {
	Name string `json:"name" binding:"required,min=2"`
	Type string `json:"type"`
}

func CreateTagsHandler(ctx *apicontext.Context, c *gin.Context) {
	var form AddTagForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)
	tag := common.Tag{Name: form.Name, Type: form.Type}
	ds := ctx.Ds
	err = ds.Create(&tag).Error
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, tag)
	return
}

func UpdateTagsHandler(ctx *apicontext.Context, c *gin.Context) {
	tagId := c.Param("id")
	var tag common.Tag
	ds := ctx.Ds
	err := ds.Where("id = ?", tagId).Find(&tag).Error
	middleware.StopExec(err)
	var form AddTagForm
	err = c.BindJSON(&form)
	middleware.StopExec(err)
	tag.Name = form.Name
	tag.Type = form.Type
	ds.Save(&tag)
	returnfun.ReturnOKJson(c, "")
	return
}

func GetCateHandler(ctx *apicontext.Context, c *gin.Context) {
	var cates []*common.Category
	ds := ctx.Ds

	name := c.Query("name")
	page, size, skip := util.GetPageAndSize(c)
	var total int
	ds.Where("name LIKE ? AND delete_t = ?", "%"+name+"%", 0).Find(&cates).Count(&total).Limit(size).Offset(skip).Order("create_t desc")
	retData := returnfun.QueryListData{Data: cates, Page: page, Size: size, Total: total}
	returnfun.ReturnOKJson(c, retData)
	return
}

type AddCateForm struct {
	Name string `json:"name" binding:"required,min=2"`
}

func CreateCateHandler(ctx *apicontext.Context, c *gin.Context) {
	var form AddCateForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)
	cate := common.Category{Name: form.Name}
	ds := ctx.Ds
	err = ds.Create(&cate).Error
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, cate)
	return
}

func UpdateCateHandler(ctx *apicontext.Context, c *gin.Context) {
	tagId := c.Param("id")
	var cate common.Category
	ds := ctx.Ds
	err := ds.Where("id = ?", tagId).Find(&cate).Error
	middleware.StopExec(err)
	var form AddTagForm
	err = c.BindJSON(&form)
	middleware.StopExec(err)
	cate.Name = form.Name
	ds.Save(&cate)
	returnfun.ReturnOKJson(c, "")
	return
}
