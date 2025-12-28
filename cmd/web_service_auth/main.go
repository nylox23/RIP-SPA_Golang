package main

import (
	"context"
	"fmt"
	"web_service_auth/internal/app/redis"

	"web_service_auth/internal/app/config"
	"web_service_auth/internal/app/dsn"
	"web_service_auth/internal/app/handler"
	"web_service_auth/internal/app/repository"
	"web_service_auth/internal/app/service"
	"web_service_auth/internal/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Carb_Calc
// @version 1.0
// @description Calculates the volume of gas emitted in a reaction

// @contact.name API Support
// @contact.url https://github.com/nylox23/DIA
// @contact.email kovalyovea@student.bmstu.ru

// @license.name AS IS (NO WARRANTY)

// @host localhost:8080
// @schemes https http
// @BasePath /
//components:
//	securitySchemes:
//		bearerAuth:
//		type: http
//		scheme: bearer
//		bearerFormat: JWT
//		description: JWT Authorization header using the Bearer scheme.

func main() {
	ctx := context.Background()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

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

	red, errRed := redis.New(ctx, conf.Redis)
	if errRed != nil {
		logrus.Fatalf("error initializing Redis service: %v", errRed)
	}

	hand := handler.NewHandler(rep, minioService, conf, red)

	application := pkg.NewApp(router, hand)
	application.RunApp()
}
