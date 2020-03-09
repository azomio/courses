package main

import (
	"github.com/azomio/courses/lesson6/pkg/grpc/user"
	"google.golang.org/grpc"
)

func main() {

	connUser, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer connUser.Close()
	userCli := user.NewUserClient(connUser)

}
