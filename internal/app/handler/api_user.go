package handler

import (
	"net/http"
	"web_service/internal/app/auth"
	"web_service/internal/app/ds"
	"web_service/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// POST /api/users/register - Регистрация пользователя
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		logrus.Error("Failed to check existing user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this login already exists"})
		return
	}

	user := &ds.Users{
		Login:       req.Login,
		Password:    req.Password,
		IsModerator: false,
	}

	if err := h.Repository.CreateUser(user); err != nil {
		logrus.Error("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := dto.UserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	c.JSON(http.StatusCreated, response)
}

// POST /api/users/login - Аутентификация пользователя
func (h *Handler) LoginUserAPI(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		logrus.Error("Failed to get user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !(req.Password == user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	response := dto.LoginResponse{
		User: dto.UserResponse{
			ID:          user.ID,
			Login:       user.Login,
			IsModerator: user.IsModerator,
		},
	}

	c.JSON(http.StatusOK, response)
}

// POST /api/users/logout - Деавторизация пользователя
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Logged out successfully"})
}

// GET /api/users/profile - Получение полей пользователя
func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	userID := auth.GetCurrentUserID()

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		logrus.Error("Failed to get user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := dto.UserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	c.JSON(http.StatusOK, response)
}

// PUT /api/users/profile - Обновление профиля пользователя
func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	userID := auth.GetCurrentUserID()

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.Repository.GetUserByID(userID)
	if err != nil {
		logrus.Error("Failed to get user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Login != "" && req.Login != existingUser.Login {
		userWithLogin, err := h.Repository.GetUserByLogin(req.Login)
		if err != nil {
			logrus.Error("Failed to check existing user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
			return
		}

		if userWithLogin != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this login already exists"})
			return
		}
	}

	updates := make(map[string]interface{})
	if req.Login != "" {
		updates["login"] = req.Login
	}

	if req.Password != "" {
		updates["password"] = req.Password
	}

	if err := h.Repository.UpdateUser(userID, updates); err != nil {
		logrus.Error("Failed to update user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	updatedUser, err := h.Repository.GetUserByID(userID)
	if err != nil {
		logrus.Error("Failed to get updated user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user"})
		return
	}

	response := dto.UserResponse{
		ID:          updatedUser.ID,
		Login:       updatedUser.Login,
		IsModerator: updatedUser.IsModerator,
	}

	c.JSON(http.StatusOK, response)
}
