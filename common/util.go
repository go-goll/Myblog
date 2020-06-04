package common

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"regexp"
	"strconv"
)

func CheckPhone(phone string) bool {
	isorno, _ := regexp.MatchString(`^(1[3|4|6|7|5|8|9|2][0-9]\d{4,8})$`, phone)
	if isorno {
		return true
	}
	return false
}

func CheckEmail(email string) bool {
	isEmail, _ := regexp.MatchString(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`, email)
	if isEmail {
		return true
	}
	return false
}

func IsEmpty(params interface{}) bool {
	//初始化变量
	var (
		flag          bool = true
		default_value reflect.Value
	)

	r := reflect.ValueOf(params)

	//获取对应类型默认值
	default_value = reflect.Zero(r.Type())
	//由于params 接口类型 所以default_value也要获取对应接口类型的值 如果获取不为接口类型 一直为返回false
	if !reflect.DeepEqual(r.Interface(), default_value.Interface()) {
		flag = false
	}
	return flag
}

var MAX_ONE_PAGE_SIZE = 100

func GetPageAndSize(c *gin.Context) (page, size, skip int) {
	p := c.Query("page")
	s := c.Query("size")

	if p != "" {
		page, _ = strconv.Atoi(p)
	} else {
		page = 1
	}

	if s != "" {
		size, _ = strconv.Atoi(s)
	} else {
		size = 10
	}

	if page < 1 {
		page = 1
	}

	if size > MAX_ONE_PAGE_SIZE {
		size = MAX_ONE_PAGE_SIZE
	}

	skip = (page - 1) * size

	return
}
