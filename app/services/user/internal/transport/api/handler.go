package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/service"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/kafka/producer"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.Service
	log     *slog.Logger
	kafka   kafka.Producer
}

func NewHandler(service service.Service, log *slog.Logger, kafka kafka.Producer) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
		kafka:   kafka,
	}
}

func (h *UserHandler) signUpUser(c *gin.Context) {
	var usrInfo entities.UserInfo
	fi := "api.Handler.signUpUser"
	errCode := 0
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	if err := c.BindJSON(&usrInfo); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := usrInfo.ValidateUserInfo(); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(ctx, &usrInfo)
	//409 и 500
	if errors.Is(err, repository.ErrAlreadyExists) {
		errCode = http.StatusConflict
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, map[string]int{
		"userId": id,
	})

	usrInfo.UsrId = id
	h.kafka.SendMessage(usrInfo)

	defer func() {
		if err != nil {
			h.log.Debug(
				fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error(),
			)
		}
	}()

}

func (h *UserHandler) getUserById(c *gin.Context) {
	fi := "api.Handler.getUserById"
	errCode := 0
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	userIdstr, ok := c.GetQuery("userId")
	if !ok {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return

	}
	//400
	userId, err := strconv.Atoi(userIdstr)
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateUserId(userId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//404 и 500
	usr, err := h.service.GetUserById(ctx, userId)
	if errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, usr)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()

}

func (h *UserHandler) getUserByEmail(c *gin.Context) {
	fi := "api.Handler.getUserByEmail"
	errCode := 0
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	email, ok := c.GetQuery("email")
	if !ok {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path")
		return

	}
	//400
	if err := entities.ValidateEmail(email); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path"+err.Error())
		return
	}

	//404 и 500
	usr, err := h.service.GetUserByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, usr)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()
}

func (h *UserHandler) editUser(c *gin.Context) {
	var usrInfo entities.UserInfo
	fi := "api.Handler.editUser"
	errCode := 0
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	userIdStr := c.Param("userId")
	if userIdStr == "" {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return
	}
	//400
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateUserId(userId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := c.BindJSON(&usrInfo); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//404 и 500
	usrInfo.UsrId = userId
	if err := h.service.UpdateUser(ctx, userId, &usrInfo); errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//200
	c.AbortWithStatusJSON(http.StatusOK, "OK")

	h.kafka.SendMessage(usrInfo)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()
}

func (h *UserHandler) verifyEmail(c *gin.Context) {
	userIdstr := c.Param("userId")
	fi := "verifyEmail"
	errCode := 0
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	if userIdstr == "" {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return
	}
	//400
	userId, err := strconv.Atoi(userIdstr)
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return

	}
	//400
	if err := entities.ValidateUserId(userId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	code := c.Query("code")
	if code == "" {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "code parametr does not exist in query")
		return
	}
	//400
	if err := entities.ValidateCode(code); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return

	}
	//404 и 500
	verified, err := h.service.VerifyCode(ctx, userId, code)
	if errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return

	}
	//200
	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"verified": verified,
	})

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()
}
