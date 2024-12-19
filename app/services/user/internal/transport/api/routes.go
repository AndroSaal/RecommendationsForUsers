package api

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	singUp := router.Group("/sign-up")
	{
		singUp.POST("", h.singUpUser)
	}

	userId := singUp.Group("/:id")
	{
		userId.GET("", h.getUserById)
		userId.PATCH("/edit", h.editUser)
		userId.PUT("/verify-email", h.verifyEmail)
	}

	email := singUp.Group("/:email")
	{
		email.GET("/", h.getUserByEmail)
	}

	return router
}
