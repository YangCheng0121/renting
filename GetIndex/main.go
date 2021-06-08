package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"renting/GetIndex/handler"
	pb "renting/GetIndex/proto"
	"renting/GetIndex/subscriber"
)

const (
	ServerName = "go.micro.srv.GetIndex" // server name
)

func main() {
	reg := consul.NewRegistry()
	// Create service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
		micro.Registry(reg),
	)

	// Initialise service
	service.Init()

	// Register handler
	if err := pb.RegisterGetIndexHandler(service.Server(), new(handler.GetIndex)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.GetIndex)); err != nil {
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
