package handler

import (
	"net/http"
	"strconv"
	"web_service_auth/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UpdateCarbonateAcidAPI godoc
// @Summary      Обновить количество кислоты
// @Description  Обновляет массу кислоты в текущей заявке-черновике пользователя
// @Tags         carbonate-acids
// @Accept       json
// @Produce      json
// @Param        id       path      int  true  "ID кислоты"
// @Param        request  body      object  true  "Масса кислоты"  SchemaExample({"mass": 10.5})
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/carbonate-acids/{id} [put]
func (h *Handler) UpdateCarbonateAcidAPI(c *gin.Context) {
	acidIDStr := c.Param("id")
	acidID, err := strconv.Atoi(acidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	var req struct {
		Mass float32 `json:"mass" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := h.GetCurrentUserID(c)

	carbonateID, err := h.Repository.GetDraftCarbonate(userID, true)
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	// Заявка существует?
	carbonate, err := h.Repository.GetCarbonateByID(uint(carbonateID))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	// Пользователь - создатель заявки?
	if carbonate.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own carbonates"})
		return
	}

	// Заявка - черновик?
	if carbonate.Status != "черновик" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft carbonates can be modified"})
		return
	}

	// Обновление количества кислоты в заявке
	if err := h.Repository.UpdateCarbonateAcidAmount(uint(carbonateID), uint(acidID), req.Mass); err != nil {
		logrus.Error("Failed to update carbonate acid mass:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acid quantity"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Acid quantity updated successfully"})
}

// DeleteCarbonateAcidAPI godoc
// @Summary      Удалить кислоту из заявки
// @Description  Удаляет кислоту из текущей заявки-черновика пользователя
// @Tags         carbonate-acids
// @Produce      json
// @Param        id   path      int  true  "ID кислоты"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/carbonate-acids/{id} [delete]
func (h *Handler) DeleteCarbonateAcidAPI(c *gin.Context) {
	acidIDStr := c.Param("id")
	acidID, err := strconv.Atoi(acidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	userID := h.GetCurrentUserID(c)

	carbonateID, err := h.Repository.GetDraftCarbonate(userID, false)
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	// Заявка существует?
	carbonate, err := h.Repository.GetCarbonateByID(uint(carbonateID))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	// Только создатель может модифицировать заявку
	if carbonate.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify your own carbonates"})
		return
	}

	// Заявка - черновик?
	if carbonate.Status != "черновик" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft carbonates can be modified"})
		return
	}

	// Удаление кислоты из заявки
	if err := h.Repository.RemoveAcidFromCarbonate(uint(carbonateID), uint(acidID)); err != nil {
		logrus.Error("Failed to remove acid from carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove acid from carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Acid removed from carbonate successfully"})
}
