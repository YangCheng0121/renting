package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"renting/PostOrders/handler"
	pb "renting/PostOrders/proto"
	"renting/PostOrders/subscriber"
)

const (
	ServerName = "go.micro.srv.PostOrders" // server name
)

func main() {
	reg := consul.NewRegistry()
	// Create service
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register handler
	if err := pb.RegisterPostOrdersHandler(service.Server(), new(handler.PostOrders)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.PostOrders)); err != nil {
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
