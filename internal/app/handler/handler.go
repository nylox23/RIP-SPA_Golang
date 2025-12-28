package handler

import (
	"net/http"
	"spa_InMemory/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetAcids(ctx *gin.Context) {
	var acids []repository.Acid
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		acids, err = h.Repository.GetAcids()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		acids, err = h.Repository.GetAcidsByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "acid_select.html", gin.H{
		"acids": acids,
		"query": searchQuery,
	})
}

func (h *Handler) GetAcid(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	acid, err := h.Repository.GetAcid(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "acid.html", gin.H{
		"acid": acid,
	})
}

func (h *Handler) GetSelected(ctx *gin.Context) {
	var acids map[repository.Acid]float32
	var err error

	acids, err = h.Repository.GetSelected()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "carbonate.html", gin.H{
		"acids": acids,
	})
}
