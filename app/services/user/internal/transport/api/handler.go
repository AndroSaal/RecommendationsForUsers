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
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400 - ошибка десериализации данных
	if err := c.BindJSON(&usrInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400 - ошибка валидации данных
	if err := usrInfo.ValidateUserInfo(); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(ctx, &usrInfo)

	//409 и 500 - уже существует и внутренняя ошибка сервера
	if errors.Is(err, repository.ErrAlreadyExists) {
		logMassage(fi, h.log, err.Error(), http.StatusConflict)
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	usrInfo.UsrId = id
	if err := h.kafka.SendMessage(usrInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200 - успешное завершение
	c.AbortWithStatusJSON(http.StatusOK, map[string]int{
		"userId": id,
	})

}

func (h *UserHandler) getUserById(c *gin.Context) {
	fi := "api.Handler.getUserById"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400 - ошибка получения userId из Query
	userIdstr, ok := c.GetQuery("userId")
	if !ok {
		logMassage(fi, h.log, "missing userId param", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return

	}
	//400 - ошибка - некорректный userId (не число)
	userId, err := strconv.Atoi(userIdstr)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400 - ошибка валидации userId
	if err := entities.ValidateUserId(userId); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//404 и 500 - ошибка NotFound и InternalServerError
	usr, err := h.service.GetUserById(ctx, userId)
	if errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200 - успешное завершение
	c.AbortWithStatusJSON(http.StatusOK, usr)

}

func (h *UserHandler) getUserByEmail(c *gin.Context) {
	fi := "api.Handler.getUserByEmail"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400 - ошибка получения email из Query
	email, ok := c.GetQuery("email")
	if !ok {
		logMassage(fi, h.log, "missing email param", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path")
		return

	}
	//400 - ошибка валидации email
	if err := entities.ValidateEmail(email); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//404 и 500 - ошибка - пользователь не найден или внутренняя ошибка сервера
	usr, err := h.service.GetUserByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200 - успешное завершение
	c.AbortWithStatusJSON(http.StatusOK, usr)
}

func (h *UserHandler) editUser(c *gin.Context) {
	var usrInfo entities.UserInfo
	fi := "api.Handler.editUser"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400 - нет параметра
	userIdStr := c.Param("userId")
	if userIdStr == "" {
		logMassage(fi, h.log, "missing userId Param", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return
	}
	//400 - некорректный параметр (не число)
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400 - параметр не прошел валидацию
	if err := entities.ValidateUserId(userId); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400 - ошибка десериализация данных в структуру
	if err := c.BindJSON(&usrInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//400 - валидация данных
	if err := usrInfo.ValidateUserInfo(); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	usrInfo.UsrId = userId
	//404 и 500 - ошибки NotFound и InternalServerError
	if err := h.service.UpdateUser(ctx, userId, &usrInfo); errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.kafka.SendMessage(usrInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200 - успешное завершение
	c.AbortWithStatusJSON(http.StatusOK, "OK")
}

func (h *UserHandler) verifyEmail(c *gin.Context) {
	userIdstr := c.Param("userId")
	fi := "verifyEmail"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400 - нет параметра UserId
	if userIdstr == "" {
		logMassage(fi, h.log, "missing userId Param", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
		return
	}
	//400 - некоректный параметр
	userId, err := strconv.Atoi(userIdstr)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return

	}
	//400 - ошибка валидации параметра
	if err := entities.ValidateUserId(userId); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400 - ошибка - код не не найден в Query
	code := c.Query("code")
	if code == "" {
		logMassage(fi, h.log, "missing email param", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "code parametr does not exist in query")
		return
	}
	//400 - ошибка валидации кода
	if err := entities.ValidateCode(code); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return

	}
	//404 и 500 - ошибка - пользователь с таким userID не найдет или внутренняя ошибка сервера
	verified, err := h.service.VerifyCode(ctx, userId, code)
	if errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//200 - успешное завершение
	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"verified": verified,
	})
}

func logMassage(fi string, log *slog.Logger, msg string, code int) {
	log.Error("Transport Level Error: " + fi + ": " + msg + "   Code : " + strconv.Itoa(code))
}
