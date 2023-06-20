package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"xorm.io/xorm"
)

var Engine = Init()

func Init() *xorm.Engine {
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

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.GetString("mysql.username"),
		config.GetString("mysql.password"),
		config.GetString("mysql.host"),
		config.GetString("mysql.port"),
		config.GetString("mysql.DB"),
		config.GetString("mysql.charset"),
	)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Printf("Xorm New Engine Error:%v", err)
		return nil
	}
	return engine
}
