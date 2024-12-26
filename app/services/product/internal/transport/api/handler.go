package api

import (
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
	kafka   *kafka.Producer
}

func NewHandler(service service.Service, log *slog.Logger, kafka *kafka.Producer) *Handler {
	return &Handler{
		service: service,
		log:     log,
		kafka:   kafka,
	}
}

func (h *Handler) addNewProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.addNewProduct"
	errCode := 0
	//400
	if err := c.BindJSON(&prdInfo); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := prdInfo.ValidateProductInfo(); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateProduct(&prdInfo)
	//500
	if err != nil {
		errCode = http.StatusInternalServerError
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//200
	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"productId": id,
	})

	//Отправка сообщения в кафку
	prdInfo.ProductId = id
	h.kafka.SendMessage(prdInfo)

	defer func() {
		if err != nil {
			h.log.Debug(
				fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error(),
			)
		}
	}()

}

func (h *Handler) updateProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.updateProduct"
	errCode := 0

	//400
	prdIdStr := c.Param("productId")
	if prdIdStr == "" {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "productId parametr does not exist in path")
		return
	}
	//400
	prdId, err := strconv.Atoi(prdIdStr)
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateProductId(prdId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := c.BindJSON(&prdInfo); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//404 и 500
	prdInfo.ProductId = prdId
	if err := h.service.UpdateProduct(prdId, &prdInfo); err == repository.ErrNotFound {
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

	//подключчеить кафку
	// h.kafka.SendMessage(prdInfo)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()
}

func (h *Handler) deleteProduct(c *gin.Context) {
	var prdInfo entities.ProductInfo
	fi := "api.Handler.deleteProduct"
	errCode := 0

	//400
	prdIdStr := c.Param("productId")
	if prdIdStr == "" {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, "productId parametr does not exist in path")
		return
	}
	//400
	prdId, err := strconv.Atoi(prdIdStr)
	if err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//400
	if err := entities.ValidateProductId(prdId); err != nil {
		errCode = http.StatusBadRequest
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//404 и 500
	prdInfo.ProductId = prdId
	if err := h.service.DeleteProduct(prdId); err == repository.ErrNotFound {
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

	//подключчеить кафку - сообщение откатить
	h.kafka.SendMessage(prdInfo)

	defer func() {
		if err != nil {
			h.log.Debug(fi + "TrasportLevelError Code : " + strconv.Itoa(errCode) + " " + err.Error())
		}
	}()
}
