package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/util/log"
	"renting/PostRet/handler"
	pb "renting/PostRet/proto"
	"renting/PostRet/subscriber"
)

const (
	ServerName = "go.micro.srv.PostRet" // server name
)

func main() {
	reg := consul.NewRegistry()
	// New Service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
		micro.Registry(reg),
	)

	// Initialise service
	service.Init()

	// Register Handler
	if err := pb.RegisterPostRetHandler(service.Server(), new(handler.PostRet)); err != nil {
		log.Fatal(err)
	}

	// Register Struct as Subscriber
	if err := micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), new(subscriber.PostRet)); err != nil {
		log.Fatal(err)
	}

	// Register Function as Subscriber
	if err := micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), subscriber.Handler); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
