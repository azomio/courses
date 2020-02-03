package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	srv := grpc.NewServer()

	RegisterMovieServer(srv, &Service{})

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting server on localhost:1234")
	srv.Serve(listener)
}

type Service struct{}

func (s *Service) GetMovie(
	c context.Context,
	req *GetMovieRequest,
) (
	resp *GetMovieResponse,
	err error,
) {

	m := MovieList()[req.MovieId]
	resp = &GetMovieResponse{
		MovieId: m.ID,
		Name:    m.Name,
		Poster:  m.Poster,
		Url:     m.MovieUrl,
	}
	return resp, nil
}

type Movie struct {
	ID       int64
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
