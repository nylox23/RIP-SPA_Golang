package handler

import (
	"net/http"
	"strconv"
	"time"
	"web_service/internal/app/auth"
	"web_service/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/carbonates/current - Получить информацию для кнопки перехода к заявке
func (h *Handler) GetCurrentCarbonateAPI(c *gin.Context) {
	userID := auth.GetCurrentUserID()

	carbonateID, _ := h.Repository.GetDraftCarbonate(userID)
	acidCount := h.Repository.GetAcidCount()

	response := dto.CarbonateIconsResponse{
		CarbonateID: carbonateID,
		AcidCount:   acidCount,
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/carbonates - Получить список заявок с фильтром
func (h *Handler) GetCarbonatesAPI(c *gin.Context) {
	var filter dto.CarbonateFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	carbonates, err := h.Repository.GetCarbonatesWithFilter(filter.Status, filter.DateFrom, filter.DateTo)
	if err != nil {
		logrus.Error("Failed to get carbonates:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonates"})
		return
	}

	response := dto.CarbonateListResponse{
		Carbonates: make([]dto.CarbonateListEntry, len(carbonates)),
	}

	for i, carbonate := range carbonates {
		response.Carbonates[i] = dto.CarbonateListEntry{
			CarbonateResponse: dto.CarbonateResponse{
				ID:         carbonate.ID,
				Status:     carbonate.Status,
				DateCreate: carbonate.DateCreate,
				DateUpdate: carbonate.DateUpdate,
				Creator:    carbonate.Creator.Login,
				Moderator:  carbonate.Moderator.Login,
				Mass:       carbonate.Mass,
			},
			Calculated: h.Repository.GetCalculated(carbonate.ID),
		}

		if carbonate.DateFinish.Valid {
			response.Carbonates[i].DateFinish = &carbonate.DateFinish.Time
		}
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/carbonates/:id - Получить заявку со списком кислот
func (h *Handler) GetCarbonateAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	acids, err := h.Repository.GetCarbonateAcids(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate acids:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate acids"})
		return
	}

	response := dto.CarbonateDetailResponse{
		CarbonateResponse: dto.CarbonateResponse{
			ID:         carbonate.ID,
			Status:     carbonate.Status,
			DateCreate: carbonate.DateCreate,
			DateUpdate: carbonate.DateUpdate,
			Creator:    carbonate.Creator.Login,
			Moderator:  carbonate.Moderator.Login,
			Mass:       carbonate.Mass,
		},
		Acids: make([]dto.CarbonateAcidResponse, len(acids)),
	}

	if carbonate.DateFinish.Valid {
		response.DateFinish = &carbonate.DateFinish.Time
	}

	for i, acid := range acids {
		response.Acids[i] = dto.CarbonateAcidResponse{
			ID:          acid.ID,
			AcidID:      acid.AcidID,
			CarbonateID: acid.CarbonateID,
			Mass:        acid.Mass,
			Result:      acid.Result,
			Acid: dto.AcidResponse{
				ID:        acid.Acid.ID,
				NameExt:   acid.Acid.NameExt,
				Info:      acid.Acid.Info,
				Name:      acid.Acid.Name,
				Hplus:     acid.Acid.Hplus,
				MolarMass: acid.Acid.MolarMass,
				Img:       acid.Acid.Img,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

// PUT /api/carbonates/:id - Обновление полей заявки
func (h *Handler) UpdateCarbonateAPI(c *gin.Context) {
	userID := auth.GetCurrentUserID()

	id, _ := h.Repository.GetDraftCarbonate(userID)

	var req dto.CarbonateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.Mass > 0 {
		updates["mass"] = req.Mass
	}
	updates["date_update"] = time.Now()

	if err := h.Repository.UpdateCarbonate(uint(id), updates); err != nil {
		logrus.Error("Failed to update carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate updated successfully"})
}

// PUT /api/carbonates/:id/form - Формирование заявки создателем
func (h *Handler) FormCarbonateAPI(c *gin.Context) {
	userID := auth.GetCurrentUserID()

	carbonateID, _ := h.Repository.GetDraftCarbonate(userID)

	carbonate, err := h.Repository.GetCarbonateByID(uint(carbonateID))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate.Mass <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mass is required and must be greater than 0"})
		return
	}

	updates := map[string]interface{}{
		"status":      "сформирован",
		"date_update": time.Now(),
	}

	if err := h.Repository.UpdateCarbonate(uint(carbonateID), updates); err != nil {
		logrus.Error("Failed to form carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to form carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate formed successfully"})
}

// PUT /api/carbonates/:id/status - Завершение/Отклонение заявки модератором
func (h *Handler) SetCarbonateStatusAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=завершен отклонен"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	if carbonate.Status != "сформирован" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only formed carbonates can be finalized or rejected"})
		return
	}

	userID := auth.GetCurrentUserID()

	updates := map[string]interface{}{
		"status":       req.Status,
		"moderator_id": userID,
		"date_finish":  time.Now(),
		"date_update":  time.Now(),
	}

	if err := h.Repository.UpdateCarbonate(uint(id), updates); err != nil {
		logrus.Error("Failed to update carbonate status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update carbonate status"})
		return
	}

	if req.Status == "завершен" {
		if err := h.calculateCarbonateResults(uint(id)); err != nil {
			logrus.Error("Failed to calculate carbonate results:", err)
		}
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate status updated successfully"})
}

// DELETE /api/carbonates/:id - Удаление заявки
func (h *Handler) DeleteCarbonateAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	userID := auth.GetCurrentUserID()

	if carbonate.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own carbonates"})
		return
	}

	if carbonate.Status != "черновик" && carbonate.Status != "отклонен" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft and rejected carbonates can be deleted"})
		return
	}

	if err := h.Repository.DeleteCarbonate(uint(id)); err != nil {
		logrus.Error("Failed to delete carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate deleted successfully"})
}

func (h *Handler) calculateCarbonateResults(carbonateID uint) error {
	acids, err := h.Repository.GetCarbonateAcids(carbonateID)
	if err != nil {
		return err
	}

	for _, acid := range acids {
		result := 22.4 * min(acid.Carbonate.Mass/100, acid.Mass*float32(acid.Acid.Hplus)/2/acid.Acid.MolarMass)

		if err := h.Repository.UpdateCarbonateAcidResult(acid.ID, result); err != nil {
			return err
		}
	}

	return nil
}
