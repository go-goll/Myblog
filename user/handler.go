package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/common"
	"github.com/leyle/ginbase/dbandmq"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
	"github.com/leyle/ginbase/util"
)

type RegisterInfo struct {
	Username string `json:"username" binding:"required,min=2,max=30"`
	Password string `json:"password" binding:"required,min=2,max=30"`
	Avatar   string `json:"avatar" binding:"required"`
	Phone    string `json:"phone" binding:"required,max=11,min=11"`
	Email    string `json:"email" binding:"required,min=2,max=30"`
}

func RegisterHandler(ctx *apicontext.Context, c *gin.Context) {
	var form RegisterInfo
	var err error
	err = c.BindJSON(&form)
	middleware.StopExec(err)
	ds := ctx.Ds
	// 加锁，防止名字重复
	lockVal, ok := dbandmq.AcquireLock(ctx.R, form.Username, dbandmq.DEFAULT_LOCK_ACQUIRE_TIMEOUT, dbandmq.DEFAULT_LOCK_KEY_TIMEOUT)
	if !ok {
		returnfun.ReturnErrJson(c, "锁定数据失败")
		return
	}
	defer dbandmq.ReleaseLock(ctx.R, form.Username, lockVal)
	var exist bool
	exist, err = common.ExistUserByUniqueField(ds, "username", form.Username)
	if err != nil {
		returnfun.ReturnErrJson(c, "连接数据库查询失败")
		return
	}
	if exist {
		returnfun.ReturnErrJson(c, "用户名已被占用")
		return
	}
	exist, err = common.ExistUserByUniqueField(ds, "phone", form.Phone)
	if err != nil {
		returnfun.ReturnErrJson(c, "连接数据库查询失败")
		return
	}
	if exist {
		returnfun.ReturnErrJson(c, "手机号已被占用")
		return
	}
	isPhone := common.CheckPhone(form.Phone)
	if !isPhone {
		returnfun.ReturnErrJson(c, "手机号校验不通过")
		return
	}
	isEmail := common.CheckEmail(form.Email)
	if !isEmail {
		returnfun.ReturnErrJson(c, "邮箱校验不通过")
		return
	}
	salt := util.GenerateDataId()
	passwd := form.Password + salt
	hashP := util.Sha256(passwd)
	user := common.User{
		Username: form.Username,
		Password: hashP,
		Salt:     salt,
		Avatar:   form.Avatar,
		Phone:    form.Phone,
		Email:    form.Email,
	}

	err = ds.Create(&user).Error
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, user)
}

type LoginIdPasswdForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginByPasswdHandler(ctx *apicontext.Context, c *gin.Context) {
	var form LoginIdPasswdForm
	var err error
	err = c.BindJSON(&form)
	middleware.StopExec(err)
	ds := ctx.Ds
	// 校验
	var user common.User
	err = ds.Where("username=?", form.Username).First(&user).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		middleware.StopExec(err)
	}
	if user.DeleteT != 0 {
		returnfun.ReturnErrJson(c, "该账号已被禁止登录，请联系管理员")
		return
	}
	passwd := form.Password + user.Salt
	if user.Password != util.Sha256(passwd) {
		returnfun.ReturnErrJson(c, "账号或密码有误，请重试")
		return
	}
	// 检查一致，生成 token ，存储到数据库，返回用户token信息
	token, err := common.GenerateToken(user.ID)
	middleware.StopExec(err)

	err = common.SaveToken(ctx.R, token, &user)
	middleware.StopExec(err)

	err = user.SaveLoginHistory(ds)
	middleware.StopExec(err)
	retData := gin.H{
		"token": token,
		"user":  user,
	}
	returnfun.ReturnOKJson(c, retData)
	return
}

func GetUserInfoHandler(ctx *apicontext.Context, c *gin.Context) {
	user := common.GetUserInfoByToken(c)
	returnfun.ReturnOKJson(c, user)
	return
}

func LogoutHandler(ctx *apicontext.Context, c *gin.Context) {
	user := common.GetUserInfoByToken(c)
	if user == nil {
		returnfun.ReturnErrJson(c, "获取用户信息失败")
		return
	}
	err := common.DeleteToken(ctx.R, user.ID)
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, "")
	return
}

// 修改自己的密码
type UpdatePasswdForm struct {
	Password string `json:"password" binding:"required,max=32,min=6"`
}

func UpdatePasswdHandler(ctx *apicontext.Context, c *gin.Context) {
	var form UpdatePasswdForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)
	// 生成新密码
	salt := util.GenerateDataId()
	p := form.Password + salt
	hashP := util.Sha256(p)
	user := common.GetUserInfoByToken(c)
	user.Salt = salt
	user.Password = hashP
	ds := ctx.Ds
	ds.Save(&user)
	// 移除当前token
	err = common.DeleteToken(ctx.R, user.ID)
	middleware.StopExec(err)
	returnfun.ReturnOKJson(c, "")
}
