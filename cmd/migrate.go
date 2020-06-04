package main

import (
	"blog/common"
	"github.com/jinzhu/gorm"
)

func MakeMigrate(Ds *gorm.DB) {
	Ds.AutoMigrate(&common.User{}, &common.LoginHistory{}, &common.Tag{}, &common.Category{}, &common.Article{})
	Ds.Model(&common.User{}).AddIndex("idx_user_createT", "create_t")
	Ds.Model(&common.LoginHistory{}).AddIndex("idx_login_history_createT", "create_t")
	Ds.Model(&common.Tag{}).AddIndex("idx_tag_createT", "create_t")
	Ds.Model(&common.Article{}).AddIndex("idx_article_createT", "create_t")
}
