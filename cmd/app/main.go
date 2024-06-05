package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	v1 "sso-service/api/http/v1"
	v2 "sso-service/api/http/v2"
	"sso-service/config"
	_ "sso-service/docs"
	"sso-service/internal/application"
	"syscall"
	"time"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @description sso-service API
// @description Это сваггер-документация для сервиса авторизации, аутентификации и управления аккаунтами на платформе KForge
// @description Все тела запросов, необходимые токены и возможные ошибки указаны в описании методов.

const (
	ConfigPath = "././config.yaml"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	if configPath == "" {
		panic("no config path")
	}

	cfg, err := config.LoadConfig(configPath)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	time.Sleep(10 * time.Second)
	app := application.NewApp(cfg)

	err = app.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer app.Shutdown()

	restServer := v1.InitRouter(app.GetServiceManager())
	grpcServer := v2.NewGrpcServer(app.GetServiceManager())

	go func() {
		err = restServer.Run(cfg.Server.Rest)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err = grpcServer.Run(cfg.Server.Grpc)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
