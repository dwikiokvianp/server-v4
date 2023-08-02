package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/middleware"
)

func Routes(c *gin.Engine) {
	AuthRoutes(c)
	c.Use(middleware.AuthMiddleware())
	DashboardRoutes(c)
	UserRoutes(c)
	StatusRoutes(c)
	EmployeeRoutes(c)
	CustomerRoutes(c)
	CityRoutes(c)
	DeliveryOrderRoutes(c)
	HandoverRoutes(c)
	WarehouseRoutes(c)
	ProvinceRoutes(c)
	CompanyRoutes(c)
	HistoryRoutes(c)
	OfficerRoutes(c)
	TransactionDetail(c)
	UpdateTransactionRoutes(c)
	OilRoutes(c)
	SummaryRoutes(c)
	TravelDeliveryRoutes(c)
	TransactionRoutes(c)
	VehicleRoutes(c)
	ProofRoutes(c)
}
