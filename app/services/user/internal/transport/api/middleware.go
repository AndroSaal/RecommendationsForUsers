package api

import (
	"fmt"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	//возвращение ошибки внутри логгера (чтобы мы увидели)
	slog.Error(fmt.Sprintf("error at newErrorResponse (%d, %s)", statusCode, message))
	//возварщение ошибки в качестве ответа (чтобы увидел клиент)
	c.AbortWithStatusJSON(statusCode, entities.ErrorResponse{
		Reason: message,
	})
}
