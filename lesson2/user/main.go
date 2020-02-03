package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", UserHandler).Methods("GET")
	http.ListenAndServe(":8081", r)
}

type LoginRespErr struct {
	Error string `json:"error"`
}

type User struct {
	Name  string
	pwd   string
	Token string
}

var UU = []User{
	User{"Bob", "god", "asdf11"},
	User{"Alice", "secret", "qwer12"},
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := r.Form.Get("token")

	fmt.Println(token)

	var usr User
	for _, u := range UU {
		if token == u.Token {
			usr = u
			break
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if usr.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(LoginRespErr{"Пользователь не найден"})
		return
	}

	json.NewEncoder(w).Encode(usr)
	return
}
