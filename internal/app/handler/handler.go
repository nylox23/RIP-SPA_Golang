package handler

import (
	"net/http"
	"web_service_auth/internal/app/config"
	"web_service_auth/internal/app/dto"
	"web_service_auth/internal/app/redis"
	"web_service_auth/internal/app/repository"
	"web_service_auth/internal/app/role"
	"web_service_auth/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository   *repository.Repository
	MinioService *service.MinioService
	Config       *config.Config
	Redis        *redis.Client
}

func NewHandler(r *repository.Repository, minioService *service.MinioService, cfg *config.Config, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:   r,
		MinioService: minioService,
		Config:       cfg,
		Redis:        redisClient,
	}
}

func (h *Handler) RegisterAPIHandler(router *gin.Engine) {

	api := router.Group("/api")
	{
		acids := api.Group("/acids")
		{
			acids.GET("", h.GetAcidsAPI)
			acids.GET("/:id", h.GetAcidAPI)
			acids.POST("/:id/toCarbonate", h.WithAuthCheck(role.User, role.Admin), h.AddAcidToCarbonateAPI)
			acids.POST("", h.WithAuthCheck(role.Admin), h.CreateAcidAPI)
			acids.PUT("/:id", h.WithAuthCheck(role.Admin), h.UpdateAcidAPI)
			acids.DELETE("/:id", h.WithAuthCheck(role.Admin), h.DeleteAcidAPI)
			acids.POST("/:id/image", h.WithAuthCheck(role.Admin), h.AddAcidImageAPI)
		}

		carbonates := api.Group("/carbonates")
		{
			carbonates.GET("/current", h.WithAuthCheck(role.User, role.Admin), h.GetCurrentCarbonateAPI)
			carbonates.GET("", h.WithAuthCheck(role.User, role.Admin), h.GetCarbonatesAPI)
			carbonates.GET("/:id", h.WithAuthCheck(role.User, role.Admin), h.GetCarbonateAPI)
			carbonates.PUT("", h.WithAuthCheck(role.User, role.Admin), h.UpdateCarbonateAPI)
			carbonates.PUT("/form", h.WithAuthCheck(role.User, role.Admin), h.FormCarbonateAPI)
			carbonates.PUT("/:id/status", h.WithAuthCheck(role.Admin), h.SetCarbonateStatusAPI)
			carbonates.DELETE("/:id", h.WithAuthCheck(role.User, role.Admin), h.DeleteCarbonateAPI)
		}

		carbonateAcids := api.Group("/carbonate-acids")
		{
			carbonateAcids.PUT("/:id", h.WithAuthCheck(role.User, role.Admin), h.UpdateCarbonateAcidAPI)
			carbonateAcids.DELETE("/:id", h.WithAuthCheck(role.User, role.Admin), h.DeleteCarbonateAcidAPI)
			carbonateAcids.PUT("/update-results", h.UpdateCalculationResultAPI)
		}

		users := api.Group("/users")
		{
			users.POST("/register", h.RegisterUserAPI)
			users.POST("/login", h.LoginUserAPI)
			users.POST("/logout", h.WithAuthCheck(role.User, role.Admin), h.LogoutUserAPI)
			users.GET("/profile", h.WithAuthCheck(role.User, role.Admin), h.GetUserProfileAPI)
			users.PUT("/profile", h.WithAuthCheck(role.User, role.Admin), h.UpdateUserProfileAPI)
		}
	}
}

const InternalToken = "12345678"

func (h *Handler) UpdateCalculationResultAPI(c *gin.Context) {
	token := c.GetHeader("X-Internal-Token")
	if token != InternalToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid internal token"})
		return
	}

	var req dto.CalcCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, item := range req.Results {
		if err := h.Repository.UpdateCarbonateAcidResult(item.ID, item.Result); err != nil {
			logrus.Errorf("Failed to update async result for ID %d: %v", item.ID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
