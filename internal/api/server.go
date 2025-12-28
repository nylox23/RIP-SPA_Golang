package api

import (
	"log"
	"spa_InMemory/internal/app/handler"
	"spa_InMemory/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/acid_select", handler.GetAcids)
	r.GET("/acid/:id", handler.GetAcid)
	r.GET("/carbonate/:id", handler.GetSelected)

	r.Run()
	log.Println("Server down")
}
