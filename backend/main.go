package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error


type Users struct {
	id int 'json:"id"'
	username string 'json:"username"'
	password string 'json:"password"'
	name string 'json:"name"'
}

type Result struct {
	code int 'json:"code"'
	data interface{} 'json:"data"'
	message string 'json:"message"'
}

func main() {
	db, err = gorm.Open("mysql", "root:@/user_manage?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println("Connection Failed", err)
	} else {
		fmt.Println("Connection Established")
	}

	// db.AutoMigrate(&Users{})
}
