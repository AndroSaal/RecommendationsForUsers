package api

import (
	"context"
	"errors"
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
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	// ошибка 400 - отсутсвие параметра userId
	userID := c.Param("userId")
	if userID == "" {
		logMassage(fi, h.log, "userId parameter is empty in path", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId parameter is empty in path")
		return
	}

	//ошибка 400 - Некорректный параметр userID
	userId, err := strconv.Atoi(userID)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId parameter incorrect in path")
		return
	} else {
		h.log.Info("%s: Get response for User with Id %d", fi, userId)
	}

	// ошибка 400 - ошибка валидации userId
	if err := entities.ValidateUserId(userId); err != nil {
		logMassage(fi, h.log, "userId validation failed"+": "+err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId validation failed"+": "+err.Error())
		return
	}

	//404 и 500 - рекомендаций нет или внутряняя ошибка сервера
	productIds, err := h.service.GetRecommendations(ctx, userId)
	if errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if productIds == nil {
		logMassage(fi, h.log, "recommendations for this user not found", http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, "recommendations for this user not found")
		return
	}

	// успешное завершение 200
	c.AbortWithStatusJSON(http.StatusOK, productIds)

}

func logMassage(fi string, log *slog.Logger, msg string, code int) {
	log.Error("Transport Level Error: " + fi + ": " + msg + "   Code : " + strconv.Itoa(code))
}
