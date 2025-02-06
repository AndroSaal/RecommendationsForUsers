package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/service"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/transport/kafka/producer"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	log     *slog.Logger
	kafka   kafka.KafkaProducer
}

func NewHandler(service service.Service, log *slog.Logger, kafka kafka.KafkaProducer) *Handler {
	return &Handler{
		service: service,
		log:     log,
		kafka:   kafka,
	}
}

func (h *Handler) addNewProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.addNewProduct"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	if err := c.BindJSON(&prdInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := prdInfo.ValidateProductInfo(); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateProduct(ctx, &prdInfo)
	//500
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//Отправка сообщения в кафку
	action := "add"
	prdInfo.ProductId = id
	if err := h.kafka.SendMessage(prdInfo, action); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"productId": id,
	})
}

func (h *Handler) updateProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.updateProduct"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	prdIdStr := c.Param("productId")
	if prdIdStr == "" {
		logMassage(fi, h.log, "productId parametr does not exist in path", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "productId parametr does not exist in path")
		return
	}
	//400
	prdId, err := strconv.Atoi(prdIdStr)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateProductId(prdId); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := c.BindJSON(&prdInfo); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//404 и 500
	prdInfo.ProductId = prdId
	if err := h.service.UpdateProduct(ctx, prdId, &prdInfo); errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	action := "update"
	if errk := h.kafka.SendMessage(prdInfo, action); errk != nil {
		logMassage(fi, h.log, errk.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, errk.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, "OK")
}

func (h *Handler) deleteProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.deleteProduct"
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	//400
	prdIdStr := c.Param("productId")
	if prdIdStr == "" {
		logMassage(fi, h.log, "productId parametr does not exist in path", http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, "productId parametr does not exist in path")
		return
	}
	//400
	prdId, err := strconv.Atoi(prdIdStr)
	if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateProductId(prdId); err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusBadRequest)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//404 и 500
	prdInfo.ProductId = prdId
	if err := h.service.DeleteProduct(ctx, prdId); errors.Is(err, repository.ErrNotFound) {
		logMassage(fi, h.log, err.Error(), http.StatusNotFound)
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		logMassage(fi, h.log, err.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//подключчеить кафку - сообщение откатить
	action := "delete"
	if errk := h.kafka.SendMessage(prdInfo, action); errk != nil {
		logMassage(fi, h.log, errk.Error(), http.StatusInternalServerError)
		newErrorResponse(c, http.StatusInternalServerError, errk.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, "OK")
}

func logMassage(fi string, log *slog.Logger, msg string, code int) {
	log.Error("Transport Level Error: " + fi + ": " + msg + "   Code : " + strconv.Itoa(code))
}
