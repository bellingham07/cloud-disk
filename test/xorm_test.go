package test

import (
	"bytes"
	"cloud-disk/core/models"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"xorm.io/xorm"
)

func TestXormTest(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "cao:cao123456@tcp(43.143.225.87:3306)/cloud-disk?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	data := make([]*models.UserBasic, 0)
	err = engine.Find(&data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
	// 将struct变成byte数组
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	// 将byte数组转换为一个buffer 然后将byte buffer以json的形式打印
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dst)
}
