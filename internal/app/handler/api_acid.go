package handler

import (
	"net/http"
	"strconv"
	"web_service_auth/internal/app/ds"
	"web_service_auth/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetAcidsAPI godoc
// @Summary      Получить список кислот
// @Description  Возвращает список кислот, отфильтрованный по названию
// @Tags         acids
// @Produce      json
// @Success      200  {object}  dto.AcidListResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/acids [get]
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

// GetAcidAPI godoc
// @Summary      Получить кислоту по ID
// @Description  Возвращает данные кислоты по ID
// @Tags         acids
// @Produce      json
// @Param        id   path      int  true  "ID кислоты"
// @Success      200  {object}  dto.AcidResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/acids/{id} [get]
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

// CreateAcidAPI godoc
// @Summary      Создать кислоту
// @Description  Создает новую кислоту без изображения
// @Tags         acids
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AcidCreateRequest  true  "Данные создавемой кислоты"
// @Success      201      {object}  dto.AcidResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/acids [post]
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

// UpdateAcidAPI godoc
// @Summary      Обновить кислоту
// @Description  Обновляет данные существующей кислоты
// @Tags         acids
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "ID кислоты"
// @Param        request  body      dto.AcidUpdateRequest true  "Обновляемые данные кислоты"
// @Success      200      {object}  dto.AcidResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/acids/{id} [put]
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

// DeleteAcidAPI godoc
// @Summary      Удалить кислоту
// @Description  Проводит "удаление" кислоты, устанавливая флаг is_delete в true
// @Tags         acids
// @Produce      json
// @Param        id   path      int  true  "ID кислоты"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/acids/{id} [delete]
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

// AddAcidToCarbonateAPI godoc
// @Summary      Добавить кислоту в заявку
// @Description  Добавляет кислоту в текущую заявку-черновик пользователя
// @Tags         acids
// @Produce      json
// @Param        id   path      int  true  "ID кислоты"
// @Success      201  {object}  dto.MessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/acids/{id}/toCarbonate [post]
func (h *Handler) AddAcidToCarbonateAPI(c *gin.Context) {
	acidIDStr := c.Param("id")
	acidID, err := strconv.Atoi(acidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid acid ID"})
		return
	}

	userID := h.GetCurrentUserID(c)

	carbonateID, err := h.Repository.GetDraftCarbonate(userID, true)
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

// AddAcidImageAPI godoc
// @Summary      Добавить изображение кислоты
// @Description  Загружает изображение для кислоты, замещая предыдущее, если есть
// @Tags         acids
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path      int   true  "ID кислоты"
// @Param        image  formData  file  true  "Файл с изображением"
// @Success      200    {object}  dto.MessageResponse
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/acids/{id}/image [post]
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
