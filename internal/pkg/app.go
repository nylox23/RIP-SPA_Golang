package pkg

import (
	"fmt"

	"web_service_auth/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("API Server starting up")

	a.Handler.RegisterAPIHandler(a.Router)

	serverAddress := fmt.Sprintf("%s:%d", a.Handler.Config.ServiceHost, a.Handler.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("API Server shut down")
}
