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

	// migrate table user
	db.AutoMigrate(&User{})

	handleRequests()

}
func handleRequests() {
	fmt.Println("Running...")
	myRoute := mux.NewRouter().StrictSlash(true)

	myRoute.HandleFunc("/user", viewUser).Methods("GET")
	myRoute.HandleFunc("/user", addUser).Methods("POST")
	myRoute.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	myRoute.HandleFunc("/user/{limit}/{offset}", getUsers).Methods("GET")
	myRoute.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRoute.HandleFunc("/user/{id}", getUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":2121", myRoute))
}

func addUser(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var user User
	json.Unmarshal(payloads, &user)

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		db.Create(&user)
		// message success
		res2 := Result{Code: 200, Data: user, Message: "User Created"}
		result, err2 := json.Marshal(res2)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var userUpdates User
	json.Unmarshal(payloads, &userUpdates)

	validate := validator.New()
	err := validate.Struct(userUpdates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		var user User
		db.First(&user, userId)
		db.Model(&user).Updates(userUpdates)
		// message success
		res2 := Result{Code: 200, Data: user, Message: "User Updated"}
		result, err2 := json.Marshal(res2)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	var user User
	db.First(&user, userId)
	db.Delete(&user)

	// message success
	res2 := Result{Code: 200, Data: user, Message: "User Deleted"}
	result, err2 := json.Marshal(res2)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func viewUser(w http.ResponseWriter, r *http.Request) {

	var users = []User{}
	db.Find(&users)
	res := Result{Code: 200, Data: users, Message: "Success Get Users"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	varr := mux.Vars(r)
	limit := varr["limit"]

	var users = []User{}
	db.Limit(limit).Find(&users)
	res := Result{Code: 200, Data: users, Message: "Success Get Users"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	var user User
	db.First(&user, userId)

	res := Result{Code: 200, Data: user, Message: "Success Get User"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
