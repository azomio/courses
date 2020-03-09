package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/azomio/courses/lesson4/pkg/grpc/user"
	"github.com/azomio/courses/lesson4/pkg/jwt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	srv := grpc.NewServer()

	user.RegisterUserServer(srv, &UserService{})

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting server on localhost: 1234")
	srv.Serve(listener)
}

type UserService struct{}

func (s *UserService) Login(
	ctx context.Context,
	in *user.LoginRequest,
) (*user.LoginResponse, error) {
	u := UU.GetByEmail(in.GetEmail())

	if u == nil {
		return &user.LoginResponse{Error: "Пользователь не найден"}, nil
	}

	if u.Pwd != in.GetPwd() {
		return &user.LoginResponse{Error: "Неправильный email или пароль"}, nil
	}

	t, err := jwt.Make(jwt.Payload{u.ID, u.Name, u.IsPaid})

	if err != nil {
		log.Printf("Token error: %v", err)
		return &user.LoginResponse{Error: "Внутренняя ошибка сервиса"}, nil
	}

	return &user.LoginResponse{Jwt: t}, nil
}

type User struct {
	ID     int
	Email  string
	Name   string
	IsPaid bool
	Pwd    string
	Token  string
}

type UserStorage []User

var UU = UserStorage{
	User{1, "bob@mail.ru", "Bob", true, "god", "1"},
	User{2, "alice@mail.ru", "Alice", false, "secret", "2"},
}

func (uu UserStorage) GetByEmail(email string) *User {
	for _, u := range uu {
		if u.Email == email {
			return &u
		}
	}
	return nil
}

type LoginRespErr struct {
	Error string `json:"error"`
}
