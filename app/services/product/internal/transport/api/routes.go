package api

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//product
	product := router.Group("/product")
	{
		product.POST("", h.addNewProduct)

		//product/{productId}
		productId := product.Group("/:productId")
		{
			productId.PATCH("", h.updateProduct)
			productId.DELETE("", h.deleteProduct)
		}
	}

	return router
}
