package handler

import (
	"net/http"
	"strconv"
	"web_service/internal/app/auth"
	"web_service/internal/app/ds"
	"web_service/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/acids - Список кислот с фильтром
func (h *Handler) GetAcidsAPI(c *gin.Context) {
	search := c.DefaultQuery("search", "")

	acids, err := h.Repository.GetAcidsWithFilter(search)
	if err != nil {
		logrus.Error("Failed to get acids:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acids"})
		return
	}

	response := dto.AcidListResponse{
		Acids: make([]dto.AcidResponse, len(acids)),
	}

	for i, acid := range acids {
		response.Acids[i] = dto.AcidResponse{
			ID:        acid.ID,
			NameExt:   acid.NameExt,
			Info:      acid.Info,
			Name:      acid.Name,
			Hplus:     acid.Hplus,
			MolarMass: acid.MolarMass,
			Img:       acid.Img,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/acids/:id - Получение одной из кислот
func (h *Handler) GetAcidAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	acid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		logrus.Error("Failed to get acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acid"})
		return
	}

	if acid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acid not found"})
		return
	}

	response := dto.AcidResponse{
		ID:        acid.ID,
		NameExt:   acid.NameExt,
		Info:      acid.Info,
		Name:      acid.Name,
		Hplus:     acid.Hplus,
		MolarMass: acid.MolarMass,
		Img:       acid.Img,
	}

	c.JSON(http.StatusOK, response)
}

// POST /api/acids - Создание кислоты (без изображения)
func (h *Handler) CreateAcidAPI(c *gin.Context) {
	var req dto.AcidCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acid := &ds.Acid{
		NameExt:   req.NameExt,
		Info:      req.Info,
		Name:      req.Name,
		Hplus:     req.Hplus,
		MolarMass: req.MolarMass,
	}

	if err := h.Repository.CreateAcid(acid); err != nil {
		logrus.Error("Failed to create acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acid"})
		return
	}

	response := dto.AcidResponse{
		ID:        acid.ID,
		NameExt:   acid.NameExt,
		Info:      acid.Info,
		Name:      acid.Name,
		Hplus:     acid.Hplus,
		MolarMass: acid.MolarMass,
	}

	c.JSON(http.StatusCreated, response)
}

// PUT /api/acids/:id - Обновление кислоты
func (h *Handler) UpdateAcidAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	var req dto.AcidUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingAcid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		logrus.Error("Failed to get acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acid"})
		return
	}

	if existingAcid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acid not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.NameExt != "" {
		updates["name_ext"] = req.NameExt
	}
	if req.Info != "" {
		updates["info"] = req.Info
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Hplus > 0 {
		updates["hplus"] = req.Hplus
	}

	if req.MolarMass > 0 {
		updates["molar_mass"] = req.MolarMass
	}

	if err := h.Repository.UpdateAcid(id, updates); err != nil {
		logrus.Error("Failed to update acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acid"})
		return
	}

	updatedAcid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		logrus.Error("Failed to get updated acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated acid"})
		return
	}

	response := dto.AcidResponse{
		ID:        updatedAcid.ID,
		NameExt:   updatedAcid.NameExt,
		Info:      updatedAcid.Info,
		Name:      updatedAcid.Name,
		Hplus:     updatedAcid.Hplus,
		MolarMass: updatedAcid.MolarMass,
		Img:       updatedAcid.Img,
	}

	c.JSON(http.StatusOK, response)
}

// DELETE /api/acids/:id - Удаление кислоты (смена статуса)
func (h *Handler) DeleteAcidAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	existingAcid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		logrus.Error("Failed to get acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acid"})
		return
	}

	if existingAcid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acid not found"})
		return
	}

	if existingAcid.Img != "" {
		if err := h.MinioService.DeleteImage(existingAcid.Img); err != nil {
			logrus.Warn("Failed to delete acid image:", err)
		}
	}

	if err := h.Repository.DeleteAcid(id); err != nil {
		logrus.Error("Failed to delete acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete acid"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Acid deleted successfully"})
}

// POST /api/acids/:id/toCarbonate - Добавление кислоты в заявку
func (h *Handler) AddAcidToCarbonateAPI(c *gin.Context) {
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

	acid, err := h.Repository.GetAcidByID(acidID)
	if err != nil {
		logrus.Error("Failed to get acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acid"})
		return
	}

	if acid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acid not found"})
		return
	}

	if err := h.Repository.AddToCarbonate(uint(acidID), uint(carbonateID)); err != nil {
		logrus.Error("Failed to add acid to carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add acid to carbonate"})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "Acid added to carbonate successfully"})
}

// POST /api/acids/:id/image - Добавление изображения к кислоте
func (h *Handler) AddAcidImageAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	existingAcid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		logrus.Error("Failed to get acid:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve acid"})
		return
	}

	if existingAcid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acid not found"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	if existingAcid.Img != "" {
		if err := h.MinioService.DeleteImage(existingAcid.Img); err != nil {
			logrus.Warn("Failed to delete old acid image:", err)
		}
	}

	objectName, err := h.MinioService.UploadImage(file, id)
	if err != nil {
		logrus.Error("Failed to upload image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	imageURL := h.MinioService.GetImageURL(objectName)
	if err := h.Repository.UpdateAcid(id, map[string]interface{}{"img": imageURL}); err != nil {
		logrus.Error("Failed to update acid with image URL:", err)
		h.MinioService.DeleteImage(objectName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acid with image"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Image uploaded successfully"})
}
