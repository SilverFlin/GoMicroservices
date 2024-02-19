package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18 ?

type User struct {
	Id       int
	Username string
	Password string
	Email    string
}

var Users []User

func main() {
	flag.Parse()
	router := httprouter.New()

	router.GET("/v1/user/:id", UserIdHandler)
	router.GET("/v1/user", UserAllHandler)
	router.POST("/v1/user", UserHandler)

	http.ListenAndServe(*addr, router)

}

func UserHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	bodyReadCloser := req.Body

	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReadCloser)
	bodyReadCloser.Close()

	content := buf.String()

	var user User

	err := json.Unmarshal([]byte(content), &user)

	if err != nil {
		w.Write([]byte("Invalid Body"))
	}

	user.Id = len(Users)
	Users = append(Users, user)

	userStr, err := json.Marshal(user)
	w.Write(userStr)
}

func UserAllHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	usersStr, err := json.Marshal(Users)

	if err != nil {
		return
	}
	w.Write(usersStr)
}
func UserIdHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		w.Write([]byte("Invalid Param"))
	}

	if id+1 > len(Users) {
		w.Write([]byte("User not found"))
		return
	}
	fmt.Println(Users[id])
	//json.Marshal()
}
