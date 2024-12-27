package api

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//recommendation
	recommendation := router.Group("/recommendation")
	{
		//recommendation/{userId}
		userId := recommendation.Group("/:userId")
		{
			userId.GET("", h.getUserRecommendations)
		}
	}

	return router
}
