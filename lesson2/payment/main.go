package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Addr     string
	UserAddr string
	WebAddr  string
}

var cfg = Config{
	Addr:     ":8083",
	UserAddr: "http://localhost:8081",
	WebAddr:  "http://localhost:8080",
}

var (
	CheckoutFormTemplate *template.Template
	MsgTemplate          *template.Template
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/checkout", checkoutFormHandler).Methods("GET")
	r.HandleFunc("/checkout", checkoutHandler).Methods("POST")

	var err error
	CheckoutFormTemplate, err = template.ParseFiles("payform.html")
	if err != nil {
		log.Fatalf("Parse template error: %v", err)
	}

	MsgTemplate, err = template.ParseFiles("msg.html")
	if err != nil {
		log.Fatalf("Parse template error: %v", err)
	}

	err = http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		log.Fatalf("Can't start server: %v", err)
	}
}

func checkoutFormHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Пользователь не указан"))
		return
	}

	err := CheckoutFormTemplate.Execute(w, struct{ Uid string }{uid})
	if err != nil {
		log.Printf("Reder template error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uid := r.FormValue("uid")
	pan := r.FormValue("pan")
	date := r.FormValue("date")
	cvc := r.FormValue("cvc")

	params := url.Values{
		"id":      []string{uid},
		"is_paid": []string{"true"},
	}

	request, err := http.NewRequest(
		http.MethodPatch,
		cfg.UserAddr+"/user",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		log.Printf("Build request error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("User response error: %v", err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read response %v", err)
		return
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Can't parse response: %v", err)
		return
	}

	err = MsgTemplate.Execute(w, struct {
		Msg     string
		BackURL string
	}{
		"Платеж успешно совершен",
		cfg.WebAddr,
	})
	if err != nil {
		log.Printf("Reder template error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
