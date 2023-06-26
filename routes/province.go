package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func ProvinceRoutes(c *gin.Engine) {

	provinceGroup := c.Group("/province")
	{
		provinceGroup.GET("", controllers.FindProvince)
		provinceGroup.GET("/:id", controllers.FindProvinceById)
	}

}
