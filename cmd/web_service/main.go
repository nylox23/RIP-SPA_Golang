package main

import (
	"fmt"

	"web_service/internal/app/config"
	"web_service/internal/app/dsn"
	"web_service/internal/app/handler"
	"web_service/internal/app/repository"
	"web_service/internal/app/service"
	"web_service/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	minioService, errMinio := service.NewMinioService()
	if errMinio != nil {
		logrus.Fatalf("error initializing Minio service: %v", errMinio)
	}

	hand := handler.NewHandler(rep, minioService)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
