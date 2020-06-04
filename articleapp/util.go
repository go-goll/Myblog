package articleapp

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/lanfeng6/Myblog/common"
)

func CheckArticle(form *AddArticleForm, ds *gorm.DB) (*common.Article, error) {
	var article common.Article
	if form.Description == "" {
		form.Description = form.Content
	}
	rs := []rune(form.Description)
	if len(rs) > 50 {
		article.Description = string(rs[0:50]) + "..."
	}
	article.Content = form.Content
	article.Image = form.Image
	article.Title = form.Title
	for _, i := range form.Tags {
		tag, err := common.ExistTag(ds, i)
		if err != nil {
			err = errors.New("校验tag[" + i + "]时出错:" + err.Error())
			return nil, err
		}
		article.Tags = append(article.Tags, tag)
	}
	cate, err := common.ExistCategory(ds, form.CategoryID)
	if err != nil {
		err = errors.New("校验categoryId[" + form.CategoryID + "]时出错:" + err.Error())
		return nil, err
	}
	article.CategoryId = form.CategoryID
	article.Category = cate
	return &article, nil
}
