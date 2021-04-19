package main

import (
	"github.com/asim/go-micro/v3"
	"renting/PostRet/handler"
	pb "renting/PostRet/proto"
	"renting/PostRet/subscriber"

	"github.com/micro/micro/v3/service/logger"
)

const (
	ServerName = "go.micro.srv.PostRet" // server name
)

func main() {
	// Create service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	// Register handler
	if err := pb.RegisterPostRetHandler(service.Server(), new(handler.PostRet)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.PostRet)); err != nil {
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
