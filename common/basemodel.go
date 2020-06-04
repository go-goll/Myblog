package common

import (
	"github.com/jinzhu/gorm"
	"github.com/leyle/ginbase/util"
	"time"
)

type BaseModel struct {
	ID      string `json:"id"`
	CreateT int64  `json:"createT"`
	UpdateT int64  `json:"updateT"`
	DeleteT int64  `json:"deleteT"`
}

func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("id", util.GenerateDataId())
	scope.SetColumn("create_t", time.Now().UnixNano()/1e6)

	return nil
}

func (base *BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("update_t", time.Now().UnixNano()/1e6)
	return nil
}
