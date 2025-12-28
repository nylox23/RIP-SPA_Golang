package handler

import (
	"net/http"
	"spa_InDb/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllAcids(ctx *gin.Context) {
	var acids []ds.Acid
	var err error

	search := ctx.Query("search")
	if search == "" {
		acids, err = h.Repository.GetAllAcids()
	} else {
		acids, err = h.Repository.SearchAcidsByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	carbonateID := h.Repository.GetDraftCarbonateID(1)

	ctx.HTML(http.StatusOK, "acids.page.tmpl", gin.H{
		"data":         acids,
		"cart_count":   h.Repository.GetAcidCount(),
		"search":       search,
		"carbonate_id": carbonateID,
	})
}

func (h *Handler) GetAcidById(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	acid, err := h.Repository.GetAcidByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "acid.page.tmpl", acid)
}

func (h *Handler) AddToCarbonate(ctx *gin.Context) {
	userID := uint(1)
	strId := ctx.PostForm("acid_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	carbonateID, err := h.Repository.GetDraftCarbonate(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.Repository.AddToCarbonate(uint(id), carbonateID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	ctx.Redirect(http.StatusFound, "/acids")
}

func (h *Handler) GetCarbonate(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status, err := h.Repository.GetCarbonateStatus(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if status == "удален" {
		ctx.Redirect(http.StatusMovedPermanently, "/acids")
		return
	}

	data, err := h.Repository.GetCarbonate(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "carbonate.page.tmpl", gin.H{
		"data": data,
		"mass": h.Repository.GetCarbonateAmount(uint(id)),
		"id":   id,
	})
}

func (h *Handler) DeleteCarbonate(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateCarbonateStatus(uint(id), "удален"); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/acids")
}
