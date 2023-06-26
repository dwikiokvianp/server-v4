package routes

import (
	"github.com/gin-gonic/gin"
)

func Routes(c *gin.Engine) {
	AuthRoutes(c)
	UserRoutes(c)
	CityRoutes(c)
	CompanyRoutes(c)
	HistoryRoutes(c)
	OfficerRoutes(c)
	OilRoutes(c)
	SummaryRoutes(c)
	TransactionRoutes(c)
	VehicleRoutes(c)
	ProofRoutes(c)
}
