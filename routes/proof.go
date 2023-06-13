package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func ProofRoutes(c *gin.Engine) {
	proofGroup := c.Group("/proof")
	{
		//proofGroup.Use(utils.AuthMiddleware)
		proofGroup.POST("/:id", controllers.CreateProof)
		proofGroup.GET("/", controllers.GetAllProofs)
		proofGroup.GET("/:id", controllers.GetProofByID)
	}
}
