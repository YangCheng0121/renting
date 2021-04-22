package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"renting/PostAvatar/handler"
	pb "renting/PostAvatar/proto"
	"renting/PostAvatar/subscriber"
)

const (
	ServerName = "go.micro.srv.PostAvatar" // server name
)

func main() {
	// Create service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	// Register handler
	if err := pb.RegisterPostAvatarHandler(service.Server(), new(handler.PostAvatar)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.PostAvatar)); err != nil {
		logger.Fatal(err)
	}

	// Register Function as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), subscriber.Handler); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
