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
	var (
		userId entities.UserId
	)

	//400
	if err := c.BindQuery(&userId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := userId.ValidateUserId(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//401 и 500
	usr, err := h.service.GetUserById(userId)
	if err == repository.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, usr)

}

func (h *Handler) getUserByEmail(c *gin.Context) {
	var (
		email entities.Email
	)

	//400
	if err := c.BindQuery(&email); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//401 и 500
	usr, err := h.service.GetUserByEmail(email)
	if err == repository.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, usr)
}

func (h *Handler) editUser(c *gin.Context) {}

func (h *Handler) verifyEmail(c *gin.Context) {}
