package helper

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity string, name string, second int) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyzeToken token的解析
func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, nil
}

func MailSendCode(mail, code string) error {
	config := viper.New()
	//在项目中查找配置文件的路径，可以使用相对路径，也可以使用绝对路径
	config.AddConfigPath("./etc")
	//配置文件名（不带扩展名）
	config.SetConfigName("core-api")
	//设置文件类型，这里是yaml文件
	config.SetConfigType("yaml")
	//查找并读取配置文件
	err := config.ReadInConfig()
	if err != nil { // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	username := config.GetString("mail.username")
	password := config.GetString("mail.password")

	e := email.NewEmail()
	e.From = "Get <" + username + ">"
	e.To = []string{mail}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码为:<h1>" + code + "</h1>")
	err = e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", username, password, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	return nil
}

func RandCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func GetUUID() string {
	return uuid.NewV4().String()
}

func OssUpload(r *http.Request) string {
	// 获取配置文件信息
	config := viper.New()
	config.SetConfigName("core-api")
	config.SetConfigType("yaml")
	config.AddConfigPath("../core/etc")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	client, err := oss.New(config.GetString("Oss.endpoint"), config.GetString("Oss.accessKeyId"), config.GetString("Oss.accessKeySecret"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	bucket, err := client.Bucket(config.GetString("Oss.bucketName"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	_, fileHeader, err := r.FormFile("file")

	key := "CloudDisk/" + fileHeader.Filename
	tempFile, _ := fileHeader.Open()

	err = bucket.PutObject(key, tempFile)
	if err != nil {
		log.Println(err)
	}

	return "https://recommendation-c.oss-cn-beijing.aliyuncs.com/" + key
}
