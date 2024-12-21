package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// service service.Service
}

func NewHandler( /*service service.Service*/ ) *Handler {
	return &Handler{
		// service: service,
	}
}

func (h *Handler) singUpUser(c *gin.Context) {
	// var usrInfo entities.UserInfo

	// if err := c.BindJSON(&usrInfo); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// if err := usrInfo.ValidateUserInfo(); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// id, err := h.service.CreateUser(usrInfo)

	// if err == repository.ErrAlreadyExists {
	// 	newErrorResponse(c, http.StatusConflict, err.Error())
	// 	return
	// } else if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
		"userId": 50,
	})

	fmt.Printf("singUpUser")

}

func (h *Handler) getUserById(c *gin.Context) {

	// //400
	// userIdstr, ok := c.GetQuery("userId")
	// if !ok {
	// 	newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path")
	// 	return

	// }

	// userId, err := strconv.Atoi(userIdstr)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// if err := entities.ValidateUserId(userId); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// //401 и 500
	// usr, err := h.service.GetUserById(userId)
	// if err == repository.ErrNotFound {
	// 	newErrorResponse(c, http.StatusNotFound, err.Error())
	// 	return
	// } else if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// //200
	c.AbortWithStatusJSON(http.StatusOK, "OK")

	fmt.Printf("getUserById")

}

func (h *Handler) getUserByEmail(c *gin.Context) {

	// //400
	// email, ok := c.GetQuery("email")
	// if !ok {
	// 	newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path")
	// 	return

	// }

	// if err := entities.ValidateEmail(email); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, "email parametr does not exist in path"+err.Error())
	// 	return
	// }

	// //401 и 500
	// usr, err := h.service.GetUserByEmail(email)
	// if err == repository.ErrNotFound {
	// 	newErrorResponse(c, http.StatusNotFound, err.Error())
	// 	return
	// } else if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// //200
	// c.AbortWithStatusJSON(http.StatusOK, usr)
	fmt.Printf("getUserByEmail")
}

func (h *Handler) editUser(c *gin.Context) {
	// var usrInfo entities.UserInfo

	// //400
	// userIdStr := c.Param("userId")
	// if userIdStr == "" {
	// 	newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
	// 	return
	// }

	// userId, err := strconv.Atoi(userIdStr)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// if err := entities.ValidateUserId(userId); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// if err := c.BindJSON(&usrInfo); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// usrInfo.UsrId = userId
	// if err := h.service.UpdateUser(usrInfo); err == repository.ErrNotFound {
	// 	newErrorResponse(c, http.StatusNotFound, err.Error())
	// 	return
	// } else if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	fmt.Printf("editUser")
}

func (h *Handler) verifyEmail(c *gin.Context) {
	// userIdstr := c.Param("userId")

	// if userIdstr == "" {
	// 	newErrorResponse(c, http.StatusBadRequest, "userId parametr does not exist in path")
	// 	return
	// }

	// userId, err := strconv.Atoi(userIdstr)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return

	// }

	// if err := entities.ValidateUserId(userId); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// code := c.Query("code")
	// if code == "" {
	// 	newErrorResponse(c, http.StatusBadRequest, "code parametr does not exist in query")
	// 	return
	// }

	// if err := entities.ValidateCode(code); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return

	// }

	// verified, err := h.service.VerifyCode(userId, code)
	// if err == repository.ErrNotFound {
	// 	newErrorResponse(c, http.StatusNotFound, err.Error())
	// 	return
	// }

	// c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
	// 	"verified": verified,
	// })

	fmt.Printf("verifyEmail")
}
