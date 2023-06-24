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

func FileSplitUpload(r *http.Request) {
	client, err := oss.New("end", "key id", "secret")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	bucketName := "bucket name"
	objectName := "n.mp4"
	// 填写本地文件的完整路径。如果未指定本地路径，则默认从示例程序所属项目对应本地路径中上传文件。

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		return
	}

	//locaFilename := "D:\\GoProject\\cloud-disk\\test\\img\\neymar.mp4"
	//locaFilename, _ := fileHeader.Open()
	locaFilename := fileHeader.Filename
	fmt.Println(fileHeader.Filename)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 将本地文件分片，且分片数量指定为3。
	chunks, err := oss.SplitFileByPartNum(locaFilename, 3)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fd, err := os.Open(locaFilename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	defer fd.Close()

	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		// oss.ContentDisposition("attachment;filename=FileName.txt"),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		// oss.ObjectStorageClass(oss.StorageStandard),
	}

	// 步骤1：初始化一个分片上传事件。
	imur, err := bucket.InitiateMultipartUpload(objectName, options...)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("imur:", imur.UploadID)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		parts = append(parts, part)
	}

	// 指定Object的读写权限为私有，默认为继承Bucket的读写权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 步骤3：完成分片上传。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("cmur:", cmur)
}
