package handler

import (
	"net/http"
	"strconv"
	"web_service/internal/app/auth"
	"web_service/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PUT /api/carbonate-acids/:id/ - Обновление количества кислоты в заявке
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

	userID := auth.GetCurrentUserID()

	carbonateID, err := h.Repository.GetDraftCarbonate(userID)
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

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

	if carbonate.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own carbonates"})
		return
	}

	if carbonate.Status != "черновик" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft carbonates can be modified"})
		return
	}

	if err := h.Repository.UpdateCarbonateAcidAmount(uint(carbonateID), uint(acidID), req.Mass); err != nil {
		logrus.Error("Failed to update carbonate acid mass:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acid quantity"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Acid quantity updated successfully"})
}

// DELETE /api/carbonate-acids/:id/ - Удаление кислоты из заявки
func (h *Handler) DeleteCarbonateAcidAPI(c *gin.Context) {
	acidIDStr := c.Param("id")
	acidID, err := strconv.Atoi(acidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	userID := auth.GetCurrentUserID()

	carbonateID, err := h.Repository.GetDraftCarbonate(userID)
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

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

	if carbonate.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify your own carbonates"})
		return
	}

	if carbonate.Status != "черновик" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft carbonates can be modified"})
		return
	}

	if err := h.Repository.RemoveAcidFromCarbonate(uint(carbonateID), uint(acidID)); err != nil {
		logrus.Error("Failed to remove acid from carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove acid from carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Acid removed from carbonate successfully"})
}
