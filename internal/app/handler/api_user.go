package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
	"web_service_auth/internal/app/ds"
	"web_service_auth/internal/app/dto"
	"web_service_auth/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// RegisterUserAPI godoc
// @Summary      Зарегистрировать пользователя
// @Description  Создает нового пользователя с предоставленными логином и паролем
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserRegisterRequest  true  "Данные для регистрации пользователя"
// @Success      201      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/users/register [post]
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is empty"})
		return
	}

	if req.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login is empty"})
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
		UUID:     uuid.New(),
		Login:    req.Login,
		Password: generateHashString(req.Password),
		Role:     role.User,
	}

	if err := h.Repository.CreateUser(user); err != nil {
		logrus.Error("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := dto.UserResponse{
		ID:    user.UUID,
		Login: user.Login,
		Role:  user.Role,
	}

	c.JSON(http.StatusCreated, response)
}

// LoginUserAPI godoc
// @Summary      Aутентификация пользователя
// @Description  Аутентифицирует пользователя с помощью логина и пароля, возвращает JWT токен
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserLoginRequest  true  "Данные для входа в систему"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/users/login [post]
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

	if user.Password != generateHashString(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(h.Config.JWT.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(h.Config.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "carb-calc-admin",
		},
		UserUUID: user.UUID,
		Role:     user.Role,
	})

	strToken, err := token.SignedString([]byte(h.Config.JWT.Token))
	if err != nil {
		logrus.Error("Failed to create token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	response := dto.LoginResponse{
		ExpiresIn:   float32(h.Config.JWT.ExpiresIn.Seconds()),
		AccessToken: strToken,
		TokenType:   "Bearer",
	}

	c.JSON(http.StatusOK, response)
}

// LogoutUserAPI godoc
// @Summary      Разлогинить пользователя
// @Description  Добавляет JWT токен в блеклист, деактивируя его
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/users/logout [post]
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	jwtStr := c.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header"})
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		logrus.Error("Failed to parse token:", err)
		return
	}

	err = h.Redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, h.Config.JWT.ExpiresIn)
	if err != nil {
		logrus.Error("Failed to write JWT to blacklist:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Logged out successfully"})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUserProfileAPI godoc
// @Summary      Получить данные пользователя
// @Description  Возвращает информацию из профиля пользователя
// @Tags         users
// @Produce      json
// @Success      200      {object}  dto.UserResponse
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/users/profile [get]
func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	userID := h.GetCurrentUserID(c)

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
		ID:    user.UUID,
		Login: user.Login,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserProfileAPI godoc
// @Summary      Обновить данные пользователя
// @Description  Обновляет информацию в профиле пользователя (логин и/или пароль)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserUpdateRequest  true  "Данные для обновления пользователя"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/users/profile [put]
func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	userID := h.GetCurrentUserID(c)

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
		updates["password"] = generateHashString(req.Password)
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
		ID:    updatedUser.UUID,
		Login: updatedUser.Login,
		Role:  updatedUser.Role,
	}

	c.JSON(http.StatusOK, response)
}
