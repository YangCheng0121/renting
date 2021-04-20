package main

import (
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"renting/web/handler"
	_ "renting/web/models"
)

const (
	ServerName = "go.micro.web.renting" // server name
)

func main() {
	// create new web Server
	srv := httpServer.NewServer(
		server.Name(ServerName),
		server.Address(":8080"),
	)

	// register router
	rou := httprouter.New()
	// 映射静态页面
	rou.NotFound = http.FileServer(http.Dir("html"))

	// 获取地区信息
	rou.GET("/api/v1.0/areas", handler.GetArea)
	// 获取图片验证码
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)
	// 获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile",handler.GetSmsCd)
	// 获取session
	rou.GET("/api/v1.0/session",handler.GetSession)
	// 注册
	rou.POST("/api/v1.0/users",handler.PostRet)
	// 登录
	rou.POST("/api/v1.0/sessions", handler.PostLogin)

	hd := srv.NewHandler(rou)

	if err := srv.Handle(hd); err != nil {
		log.Fatal(err)
	}

	// Create service
	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(registry.NewRegistry()),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
