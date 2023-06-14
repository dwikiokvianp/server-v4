package routes

import (
	"github.com/gin-gonic/gin"
)

func Routes(c *gin.Engine) {
	AuthRoutes(c)
	UserRoutes(c)
	OfficerRoutes(c)
	OilRoutes(c)
	TransactionRoutes(c)
	VehicleRoutes(c)
	ProofRoutes(c)
}
