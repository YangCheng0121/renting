package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
	"renting/GetImageCd/handler"
	pb "renting/GetImageCd/proto"
	"renting/GetImageCd/subscriber"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
)

const (
	ServerName = "go.micro.srv.GetImageCd" // server name
)

func main() {
	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	// Create service
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name(ServerName),
		micro.Version("latest"),
	)
	// Initialise service
	service.Init()

	// Register handler
	if err := pb.RegisterGetImageCdHandler(service.Server(), new(handler.GetImageCd)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.GetImageCd)); err != nil {
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
