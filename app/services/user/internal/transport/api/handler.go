package api

import (
	"net/http"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) singUpUser(c *gin.Context) {
	var usrInfo entities.UserInfo

	if err := c.BindJSON(&usrInfo); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := usrInfo.ValidateUserInfo(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(usrInfo)

	if err == repository.ErrAlreadyExists {
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"userId": id,
	})

}

func (h *Handler) getUserById(c *gin.Context) {
	// var (
	// 	userId entities.UserId
	// )

	// if err := c.BindQuery(&email); err != nil {
	// 	if err == c.Err
	// }
}

func (h *Handler) getUserByEmail(c *gin.Context) {
	// var (
	// 	userId entities.UserId
	// )

	// if err := c.BindQuery(&email); err != nil {
	// 	if err == c.Err
	// }
}

func (h *Handler) editUser(c *gin.Context) {}

func (h *Handler) verifyEmail(c *gin.Context) {}
