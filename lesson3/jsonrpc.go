package main

import (
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
)

func main() {
	srv := rpc.NewServer()

	srv.RegisterCodec(json.NewCodec(), "application/json")

	srv.RegisterService(new(Service), "Service")
	http.Handle("/rpc", srv)

	log.Println("Starting server on localhost:1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

type Service struct{}

func (h *Service) GetMovie(r *http.Request, in *struct{ Id int }, out *Movie) error {
	*out = MovieList()[in.Id]
	return nil
}

type Movie struct {
	ID       int
	Name     string
	Poster   string
	MovieUrl string
}

func MovieList() []Movie {
	return []Movie{
		Movie{0, "Бойцовский клуб", "/static/posters/fightclub.jpg", "https://youtu.be/qtRKdVHc-cE"},
		Movie{1, "Крестный отец", "/static/posters/father.jpg", "https://youtu.be/ar1SHxgeZUc"},
		Movie{2, "Криминальное чтиво", "/static/posters/pulpfiction.jpg", "https://youtu.be/s7EdQ4FqbhY"},
	}
}
