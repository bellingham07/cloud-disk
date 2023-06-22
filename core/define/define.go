package define

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

var JwtKey = "cloud-disk-key"

// CodeLength 验证码长度
var CodeLength = 6

// LoginCodePrefix 验证码前缀
var LoginCodePrefix = "LoginCode:"

// CodeExpireTime 验证码过期时间(s)
var CodeExpireTime = 300 * time.Second

// PageSize 分页默认参数
var PageSize = 20
