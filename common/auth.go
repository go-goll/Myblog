package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/lanfeng6/Myblog/apicontext"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
	"github.com/leyle/ginbase/util"
	"strconv"
	"strings"
	"time"
)

var AesKey = util.Md5("wangguo") // 32 byte 使用加密方法就是 aes-256-cfb
const TokenRedisPrefix = "USER:TOKEN:USERID"
const combineKeyLength = 2

type TokenVal struct {
	Token string        `json:"token"`
	User  *User         `json:"user"`
	T     *util.CurTime `json:"t"`
}

func CombineRawString(userId string) string {
	t := time.Now().Unix()
	text := fmt.Sprintf("%s|%d", userId, t)
	return text
}
func GenerateToken(userId string) (string, error) {
	text := CombineRawString(userId)

	token, err := util.Encrypt([]byte(AesKey), text)
	if err != nil {
		Logger.Errorf("", "给用户[%s]生成token时，调用aes加密失败, %s", userId, err.Error())
		return "", err
	}

	// 在用 base64 编码
	b64Token := base64.StdEncoding.EncodeToString([]byte(token))

	return b64Token, nil
}

// 存储为 key 是 userid， 值是 tokenvalue
func SaveToken(r *redis.Client, token string, user *User) error {
	tkVal := &TokenVal{
		Token: token,
		User:  user,
		T:     util.GetCurTime(),
	}

	tkDump, _ := jsoniter.Marshal(&tkVal)

	key := generateTokenKey(user.ID)
	_, err := r.Set(key, tkDump, 0).Result()
	if err != nil {
		Logger.Errorf("", "存储用户[%s]的token到redis失败, %s", user.ID, err.Error())
		return err
	}

	return nil
}
func generateTokenKey(userId string) string {
	return fmt.Sprintf("%s:%s", TokenRedisPrefix, userId)
}

// 解析 token
func ParseToken(token string) (string, int64, error) {
	// 先 base64 解码
	de64Token, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		Logger.Errorf("", "base64解码token[%s]失败, %s", token, err.Error())
		return "", 0, err
	}

	// 再 aes 解密
	text, err := util.Decrypt([]byte(AesKey), string(de64Token))
	if err != nil {
		Logger.Errorf("", "aes解密token[%s]失败, %s", de64Token, err.Error())
		return "", 0, err
	}

	userId, t, err := ParseCombinedRawString(text)

	return userId, t, err
}

func ParseCombinedRawString(text string) (string, int64, error) {
	Logger.Debugf("", "ParseCombinedRawString[%s]", text)
	infos := strings.Split(text, "|")
	if len(infos) != combineKeyLength {
		return "", 0, errors.New("Invalid token, maybe old token,please logout and login again")
	}
	userId := infos[0]
	st := infos[1]

	t, _ := strconv.ParseInt(st, 10, 64)

	return userId, t, nil
}

func DeleteToken(r *redis.Client, userId string) error {
	key := generateTokenKey(userId)
	_, err := r.Del(key).Result()
	if err != nil && err != redis.Nil {
		Logger.Errorf("", "移除用户[%s]token失败, %s", userId, err.Error())
		return err
	}

	Logger.Infof("", "移除用户[%s]token成功", userId)
	return nil
}

// 角色验证的数据库相关
var AuthOption = &apicontext.Context{}

const AuthResultCtxKey = "AUTHRESULT"

type AuthResult struct {
	Result int   `json:"result"` // 验证结果，见上面字典
	User   *User `json:"user"`   // 用户信息
}

// 验证 token
func CheckToken(r *redis.Client, token string) (*TokenVal, error) {
	// 先解析 token
	userId, t, err := ParseToken(token)
	if err != nil {
		return nil, err
	}
	Logger.Debugf("", "CheckToken 时，parsetoken成功，用户[%s]，token生成时间[%s]", userId, util.FmtTimestampTime(t))

	// 从 redis 中读取 tokenval 信息
	key := generateTokenKey(userId)
	data, err := r.Get(key).Result()
	if err != nil {
		Logger.Errorf("", "CheckToken 时，从redis读取指定用户[%s]的tokenval失败, %s", userId, err.Error())
		return nil, err
	}

	var tkVal *TokenVal
	err = jsoniter.UnmarshalFromString(data, &tkVal)
	if err != nil {
		Logger.Errorf("", "CheckToken 时，反序列化从 redis 读取回来的用户[%s]的数据失败, %s", userId, err.Error())
		return nil, err
	}

	if tkVal.Token != token {
		// token 被重新生成了，原 token 失效
		err = fmt.Errorf("token失效")
		Logger.Infof("", "验证用户[%s][%s]的token[%s]时，传递token与redis保存token不一致，待验证token已失效", tkVal.User.ID, tkVal.User.Username, token)
		return nil, err
	}

	Logger.Debugf("", "CheckToken 成功，用户[%s][%s]", tkVal.User.ID, tkVal.User.Username)

	return tkVal, nil
}

func AuthToken(ao *redis.Client, token string) (*User, error) {
	tkVal, err := CheckToken(ao, token)
	if err != nil {
		Logger.Errorf("", "AuthToken 时，token验证失败, %s", err.Error())
		return nil, err
	}
	return tkVal.User, nil
}

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		Logger.Error(middleware.GetReqId(c), "请求接口中无token值")
		returnfun.Return401Json(c, "No token")
		return
	}

	// 判断是否是 api token，如果是，直接返回有权限的操作

	user, err := AuthToken(AuthOption.R, token)
	if err != nil {
		returnfun.Return401Json(c, "Invalid token")
		return
	}
	c.Set(AuthResultCtxKey, user)

	c.Next()
}

func GetUserInfoByToken(c *gin.Context) *User {
	ar, exist := c.Get(AuthResultCtxKey)
	if !exist {
		return nil
	}
	result := ar.(*User)
	return result
}
