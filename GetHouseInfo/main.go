package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"renting/GetHouseInfo/handler"
	pb "renting/GetHouseInfo/proto"
	"renting/GetHouseInfo/subscriber"
)

const (
	ServerName = "go.micro.srv.GetHouseInfo" // server name
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
	if err := pb.RegisterGetHouseInfoHandler(service.Server(), new(handler.GetHouseInfo)); err != nil {
		logger.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber(ServerName, service.Server(), new(subscriber.GetHouseInfo)); err != nil {
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
