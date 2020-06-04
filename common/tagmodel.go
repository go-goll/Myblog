package common

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Name string `json:"name"`
	Type string `json:"type" gorm:"default:'OTHER'"`
	BaseModel
}

type Category struct {
	Name     string     `json:"name"`
	Articles []*Article `json:"articles,omitempty" gorm:"FOREIGNKEY:CategoryId;ASSOCIATION_FOREIGNKEY:ID"`
	BaseModel
}

func ExistTag(Db *gorm.DB, value string) (*Tag, error) {
	var tag Tag
	err := Db.Where("id = ?", value).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func ExistCategory(Db *gorm.DB, value string) (*Category, error) {
	var cate Category
	err := Db.Where("id = ?", value).First(&cate).Error
	if err != nil {
		return nil, err
	}

	return &cate, nil
}
