package api

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() {
	router := gin.New()

	singUp := router.Group("/sign-up")
	{
		singUp.POST("", h.singUpUser)
		singUp.GET("", h.getUser)
	}

	userId := singUp.Group("/:id")
	{
		userId.PATCH("/edit", h.editUser)
		userId.PUT("/verify-email", h.verifyEmail)
	}
}
