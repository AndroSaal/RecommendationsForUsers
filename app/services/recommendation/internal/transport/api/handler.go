package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	log     *slog.Logger
}

func NewHandler(service service.Service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) getUserRecommendations(c *gin.Context) {
	fi := "api.Handler.getUserRecommendations"
	errCode := 0
	//400
	// userInterface, ok := c.Get("userId")
	// if userInterface == nil || !ok {
	// 	errCode = http.StatusBadRequest
	// 	newErrorResponse(c, http.StatusBadRequest, "userId parameter does not exist in path")
	// 	return
	// }

	// userId, ok := userInterface.(int)
	// if !ok {
	// 	errCode = http.StatusBadRequest
	// 	newErrorResponse(c, http.StatusBadRequest, "userId parameter incorrect in path")
	// 	return
	// }

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "userId parameter incorrect in path")
		return
	} else {
		h.log.Info("%s: Get response for User with Id %d", fi, userId)
	}

	if err := entities.ValidateUserId(userId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "userId validation failed"+": "+err.Error())
		return
	}

	//404 и 500
	productIds, err := h.service.GetRecommendations(userId)
	if err == repository.ErrNotFound {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if productIds == nil {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, "recommendations for this user not found")
		return

	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, productIds)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()

}
