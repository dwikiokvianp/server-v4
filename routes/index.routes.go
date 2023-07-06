package routes

import (
	"github.com/gin-gonic/gin"
)

func Routes(c *gin.Engine) {
	AuthRoutes(c)
	DriverRoutes(c)
	UserRoutes(c)
	CityRoutes(c)
	DeliveryOrderRoutes(c)
	HandoverRoutes(c)
	WarehouseRoutes(c)
	ProvinceRoutes(c)
	CompanyRoutes(c)
	HistoryRoutes(c)
	OfficerRoutes(c)
	TransactionDetail(c)
	OilRoutes(c)
	SummaryRoutes(c)
	TravelDeliveryRoutes(c)
	TransactionRoutes(c)
	VehicleRoutes(c)
	ProofRoutes(c)
}
