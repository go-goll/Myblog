package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	Username string     `json:"username" gorm:"unique_index;not null"`
	Avatar   string     `json:"avatar"`
	Salt     string     `json:"-"`
	Password string     `json:"-"`
	Phone    string     `json:"phone" gorm:"unique_index"`
	Email    string     `json:"email" gorm:"unique_index"`
	Articles []*Article `json:"articles,omitempty" gorm:"FOREIGNKEY:UserId;ASSOCIATION_FOREIGNKEY:ID"`
	BaseModel
}

type LoginHistory struct {
	Username string `json:"username"`
	BaseModel
}

func ExistUserByUniqueField(Db *gorm.DB, field, value string) (bool, error) {
	var user User
	err := Db.Select("id").Where(gorm.ToColumnName(field)+" = ?", value).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return false, err
	}
	if user.ID != "" {
		fmt.Printf("field:%s,value:%s.\n", field, value)
		return true, nil
	}
	return false, nil
}

func (user *User) SaveLoginHistory(Db *gorm.DB) error {
	his := LoginHistory{Username: user.Username}
	err := Db.Debug().Create(&his).Error
	if err != nil {
		return err
	}
	return nil
}
