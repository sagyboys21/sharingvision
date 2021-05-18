package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type User struct {
	Id       int    `validate:"omitempty,uuid"`
	Username string `validate:"required,gte=3"`
	Password string `validate:"required,gte=7"`
	Name     string `validate:"required,gte=3"`
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

	handleRequests()
	// migrate table user
	db.AutoMigrate(&User{})

}
func handleRequests() {
	fmt.Println("Running at port 9999")
	myRoute := mux.NewRouter().StrictSlash(true)

	myRoute.HandleFunc("/getuser", getUser)
	myRoute.HandleFunc("/adduser", addUser).Methods("POST")
	myRoute.HandleFunc("/getuser", getUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":9999", myRoute))
}

func addUser(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var user User
	json.Unmarshal(payloads, &user)

	// user := User{
	// 	Username: "",
	// 	Password: "",
	// 	Name:     "agungcoba",
	// }

	validate := validator.New()
	err := validate.Struct(user)

	// message error
	res1 := Result{Code: 212, Data: user, Message: err.Error()}
	resulterr, err1 := json.Marshal(res1)
	// message success
	res2 := Result{Code: 200, Data: user, Message: "User Created"}
	result, err := json.Marshal(res2)

	if err1 == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resulterr)
	}

	db.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	var users = []User{}

	db.Find(&users)
	res := Result{Code: 200, Data: users, Message: "Success Get Products"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
