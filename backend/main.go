package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/user_manage?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println("Connection Failed", err)
	} else {
		fmt.Println("Connection Established")
	}

	db.AutoMigrate(&User{})
}
