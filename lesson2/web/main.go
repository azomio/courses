package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

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

	http.ListenAndServe(":8082", r)
}

type MainPage struct {
	Movies *[]Movie
	User   *User
}

type User struct {
	Name string
}

type Movie struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Poster   string `json:"poster"`
	MovieUrl string `json:"movie_url"`
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
	}

	err = TT.MovieList.Execute(w, page)
	if err != nil {
		log.Printf("Get user error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getMovies() (*[]Movie, error) {
	mm := &[]Movie{}
	err := get("http://localhost:8080/movie", mm)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func getUser(r *http.Request) (*User, error) {
	ses, err := r.Cookie("session")
	if ses == nil {
		return nil, err
	}

	res := &struct {
		User
		Error string
	}{}
	err = get("http://localhost:8081/user?token="+ses.Value, res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, fmt.Errorf(res.Error)
	}

	return &User{
		Name: res.Name,
	}, nil
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
