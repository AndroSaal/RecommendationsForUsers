package api

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	singUp := router.Group("/sign-up")
	{
		// /sing-up
		singUp.POST("", h.singUpUser)

		// /sing-up/userId
		userId := singUp.Group("/userId")
		{
			userId.GET("", h.getUserById)
		}

		// /sing-up/email
		email := singUp.Group("/email")
		{
			email.GET("", h.getUserByEmail)
		}

		// /sing-up/{userId}
		userIdInPath := singUp.Group("/:userId")
		{
			// /sing-up/{userId}/verify-email
			userIdInPath.PUT("/verify-email", h.verifyEmail)

			// /sing-up/{userId}/edit
			userIdInPath.PATCH("/edit", h.editUser)
		}
	}

	return router
}
