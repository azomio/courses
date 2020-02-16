package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	Addr         string
	UserGRPCAddr string
	UserAddr     string
	MovieAddr    string
	PaymentAddr  string
}

var cfg = Config{
	Addr:         ":8080",
	UserGRPCAddr: ":1234",
	MovieAddr:    "http://localhost:8081",
	UserAddr:     "http://localhost:8082",
	PaymentAddr:  "http://localhost:8083/",
}

var TT struct {
	MovieList *template.Template
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", MainHandler)

	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	var err error
	TT.MovieList, err = template.ParseFiles("template/main.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}

type MainPage struct {
	Movies *[]Movie
	User   User
	PayURL string
}

type User struct {
	ID     int
	Name   string
	IsPaid bool
}

type Movie struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Poster   string `json:"poster"`
	MovieUrl string `json:"movie_url"`
	IsPaid   bool   `json:"is_paid"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	page := MainPage{}

	var err error
	page.Movies, err = getMovies()
	if err != nil {
		log.Printf("Get movie error: %v", err)
	}

	page.User, err = getUser(r)
	if err != nil {
		log.Printf("Get user error: %v", err)
	} else {
		page.PayURL = cfg.PaymentAddr + "/checkout?uid=" + strconv.Itoa(page.User.ID)
	}

	err = TT.MovieList.Execute(w, page)
	if err != nil {
		log.Printf("Get user error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getMovies() (*[]Movie, error) {
	mm := &[]Movie{}
	err := get(cfg.MovieAddr+"/movies", mm)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func getUser(r *http.Request) (usr User, err error) {
	ses, err := r.Cookie("session")
	if ses == nil {
		return usr, err
	}

	res := &struct {
		User
		Error string
	}{}
	err = get(cfg.UserAddr+"/user?token="+ses.Value, res)
	if err != nil {
		return usr, err
	}

	log.Printf("user service response: %+v", res)

	if res.Error != "" {
		return usr, fmt.Errorf(res.Error)
	}

	usr.ID = res.ID
	usr.Name = res.Name
	usr.IsPaid = res.IsPaid

	return usr, nil
}

func get(url string, out interface{}) error {
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("make request error: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response error: %w", err)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return fmt.Errorf("parse response error '%s': %w", body, err)
	}

	return nil
}
