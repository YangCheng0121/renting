package main

import (
	"github.com/asim/go-micro/v3"
	"renting/GetUserInfo/handler"
	pb "renting/GetUserInfo/proto"
	"renting/GetUserInfo/subscriber"

	"github.com/micro/micro/v3/service/logger"
)

const (
	ServerName = "go.micro.srv.GetUserInfo" // server name
)

func main() {
	// Create service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	// Register handler
	if err := pb.RegisterGetUserInfoHandler(service.Server(), new(handler.GetUserInfo)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.GetUserInfo)); err != nil {
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
