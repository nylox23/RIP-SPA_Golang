package handler

import (
	"web_service/internal/app/repository"
	"web_service/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository   *repository.Repository
	MinioService *service.MinioService
}

func NewHandler(r *repository.Repository, minioService *service.MinioService) *Handler {
	return &Handler{
		Repository:   r,
		MinioService: minioService,
	}
}

func (h *Handler) RegisterAPIHandler(router *gin.Engine) {

	api := router.Group("/api")
	{
		acids := api.Group("/acids")
		{
			acids.GET("", h.GetAcidsAPI)
			acids.GET("/:id", h.GetAcidAPI)
			acids.POST("", h.CreateAcidAPI)
			acids.PUT("/:id", h.UpdateAcidAPI)
			acids.DELETE("/:id", h.DeleteAcidAPI)
			acids.POST("/:id/toCarbonate", h.AddAcidToCarbonateAPI)
			acids.POST("/:id/image", h.AddAcidImageAPI)
		}

		carbonates := api.Group("/carbonates")
		{
			carbonates.GET("/current", h.GetCurrentCarbonateAPI)
			carbonates.GET("", h.GetCarbonatesAPI)
			carbonates.GET("/:id", h.GetCarbonateAPI)
			carbonates.PUT("", h.UpdateCarbonateAPI)
			carbonates.PUT("/form", h.FormCarbonateAPI)
			carbonates.PUT("/:id/status", h.SetCarbonateStatusAPI)
			carbonates.DELETE("/:id", h.DeleteCarbonateAPI)
		}

		carbonateAcids := api.Group("/carbonate-acids")
		{
			carbonateAcids.PUT("/:id", h.UpdateCarbonateAcidAPI)
			carbonateAcids.DELETE("/:id", h.DeleteCarbonateAcidAPI)
		}

		users := api.Group("/users")
		{
			users.POST("/register", h.RegisterUserAPI)
			users.POST("/login", h.LoginUserAPI)
			users.POST("/logout", h.LogoutUserAPI)
			users.GET("/profile", h.GetUserProfileAPI)
			users.PUT("/profile", h.UpdateUserProfileAPI)
		}
	}
}
