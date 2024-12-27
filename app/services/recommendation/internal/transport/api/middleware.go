package api

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/entities"
	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	//возвращение ошибки внутри логгера (чтобы мы увидели)
	slog.Error(message)
	//возварщение ошибки в качестве ответа (чтобы увидел клиент)
	c.AbortWithStatusJSON(statusCode, entities.ErrorResponse{
		Reason: message,
	})
}
