package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SummaryFuelOutData(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Total Delivery",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "truck.svg",
		},
		{
			"id":         2,
			"title":      "Jetty Fuel Out",
			"quantity":   "160.000",
			"percentage": "25.47",
			"imgSrc":     "jetty.svg",
		},
		{
			"id":         3,
			"title":      "SPOB Solar Out",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "spob.svg",
		},
		{
			"id":         4,
			"title":      "SPOB MFO Out",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "spob.svg",
		},
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}

func TransactionSummaryData(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Total Delivery",
			"quantity":   "30",
			"percentage": "25.36",
			"imgSrc":     "cart.svg",
		},
		{
			"id":         2,
			"title":      "Total Sales Order",
			"quantity":   "16.000",
			"percentage": "25.47",
			"imgSrc":     "speedtest.svg",
		},
		{
			"id":         3,
			"title":      "Outstanding Sales",
			"quantity":   "15",
			"percentage": "25.36",
			"imgSrc":     "notes.svg",
		},
		{
			"id":         4,
			"title":      "Total Sales",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "calendar.svg",
		},
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, data)
}

func SummaryJettyData(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Stock In Hand",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "oil_barrel.svg",
		},
		{
			"id":         2,
			"title":      "Fuel In",
			"quantity":   "160.000",
			"percentage": "25.47",
			"imgSrc":     "arrow-up.svg",
		},
		{
			"id":         3,
			"title":      "Fuel Out",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "arrow-down.svg",
		},
		{
			"id":         4,
			"title":      "Tomorrow Demand",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "callender-jetty.svg",
		},
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, data)
}

func SummarySPOBData(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Solar Stock In Hand",
			"quantity":   "640.000",
			"percentage": "25.36",
			"imgSrc":     "spob-blue.svg",
		},
		{
			"id":         2,
			"title":      "Solar Fuel Out",
			"quantity":   "160.000",
			"percentage": "25.47",
			"imgSrc":     "arrow-down-blue.svg",
		},
		{
			"id":         3,
			"title":      "MFO Stock In Hand",
			"quantity":   "500.000",
			"percentage": "25.36",
			"imgSrc":     "spob-blue.svg",
		},
		{
			"id":         4,
			"title":      "MFO Fuel Out",
			"quantity":   "160.000",
			"percentage": "25.36",
			"imgSrc":     "arrow-down-blue.svg",
		},
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, data)
}

func SummaryDeliveryData(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Total Scheduled Delivery",
			"quantity":   "10",
			"percentage": "25.36",
			"imgSrc":     "delivery-car.svg",
		},
		{
			"id":         2,
			"title":      "On Pickup",
			"quantity":   "20",
			"percentage": "25.36",
			"imgSrc":     "delivery-car.svg",
		},
		{
			"id":         3,
			"title":      "On Delivery",
			"quantity":   "50",
			"percentage": "25.36",
			"imgSrc":     "delivery-car.svg",
		},
		{
			"id":         4,
			"title":      "Success",
			"quantity":   "28",
			"percentage": "25.36",
			"imgSrc":     "delivery-car.svg",
		},
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}
