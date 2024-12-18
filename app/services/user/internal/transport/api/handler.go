package api

import (
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) singUpUser(c *gin.Context) {}

func (h *Handler) getUser(c *gin.Context) {}

func (h *Handler) editUser(c *gin.Context) {}

func (h *Handler) verifyEmail(c *gin.Context) {}
