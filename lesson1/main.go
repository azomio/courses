package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

var App struct {
	TT struct {
		MovieList    *template.Template
		Movie        *template.Template
		LoginHandler *template.Template
		UserPage     *template.Template
	}
}

func prepareTemplates() {
	tDir := "web/template/"
	tLayout := tDir + "layout.html"

	tt := []struct {
		T    **template.Template
		Path string
	}{
		{&App.TT.MovieList, "movie_list.html"},
		{&App.TT.Movie, "movie.html"},
		{&App.TT.Login, "login.html"},
		{&App.TT.UserPage, "user_page.html"},
	}

	var err error
	for _, t := range tt {
		*t.T, err = template.ParseFiles(tLayout, tDir+t.Path)
		if err != nil {
			onErr(err, "Can't load template "+t.Path)
		}
	}

	return
}

func main() {

	prepareTemplates()

	r := mux.NewRouter()
	r.HandleFunc("/", MovieListHandler)
	r.HandleFunc("/movie/{id:[0-9]+}", MovieHandler)
	r.HandleFunc("/login", LoginFormHandler).Methods("GET")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/logout", LogoutHandler).Methods("POST")

	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	onErr(srv.ListenAndServe(), "Can't run server")
}

func MovieListHandler(w http.ResponseWriter, r *http.Request) {

	data := struct{ Movies []*Movie }{MovieList()}

	App.TT.MovieList.Execute(w, data)
	return
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		HandleErr(err, w, r)
		return
	}

	data := MovieList()[id]

	App.TT.Movie.Execute(w, data)
	return
}

func HandleErr(e error, w http.ResponseWriter, r *http.Request) {
	return
}

type LoginPage struct {
	Hint string
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {

	u := getUserFromRequest(r)
	if u == nil {
		App.TT.Login.Execute(w, nil)
	} else {
		App.TT.UserPage.Execute(w, u)
	}

	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:   "session",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	if len(r.PostForm["login"]) == 0 ||
		len(r.PostForm["password"]) == 0 {
		App.TT.Login.Execute(w, nil)
		return
	}

	login := r.PostForm["login"][0]
	pwd := r.PostForm["password"][0]

	u := LoadUserByLogin(login)

	if u == nil {
		App.TT.Login.Execute(w, &LoginPage{"Не правильные логин или пароль"})
		return
	}

	if u.Password != pwd {
		App.TT.Login.Execute(w, &LoginPage{"Не правильные логин или пароль"})
		return
	}

	cookie := http.Cookie{Name: "session", Value: u.Session, Expires: time.Now().Add(24 * time.Hour)}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusFound)

	return
}

func getUserFromRequest(r *http.Request) *User {

	ses, _ := r.Cookie("session")
	if ses == nil {
		return nil
	}

	u := LoadUserBySession(ses.Value)
	if u == nil {
		return nil
	}

	return u
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	for _, c := range cookie {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func onErr(err error, msg string) {
	if err != nil {
		fmt.Printf("[ERR] %s: %v\n", msg, err)
		os.Exit(1)
	}
}
