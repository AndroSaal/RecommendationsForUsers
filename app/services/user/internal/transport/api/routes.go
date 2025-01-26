package api

import "github.com/gin-gonic/gin"

func (h *UserHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	user := router.Group("/user")

	singUp := user.Group("sign-up")
	{
		// POST user/sing-up
		singUp.POST("", h.signUpUser)

		// GET user/sing-up/userId
		userId := singUp.Group("/userId")
		{
			userId.GET("", h.getUserById)
		}

		// GET user/sing-up/email
		email := singUp.Group("/email")
		{
			email.GET("", h.getUserByEmail)
		}

		// user/sing-up/{userId}
		userIdInPath := singUp.Group("/:userId")
		{
			// PUT user/sing-up/{userId}/verify-email
			userIdInPath.PUT("/verify-email", h.verifyEmail)

			// PATCH user/sing-up/{userId}/edit
			userIdInPath.PATCH("/edit", h.editUser)
		}
	}

	return router
}
