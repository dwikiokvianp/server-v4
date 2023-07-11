package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"server-v2/config"
	"server-v2/routes"
)

func main() {
	failedLoadEnv := godotenv.Load(".env")
	if failedLoadEnv != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl, port := os.Getenv("DB_URL"), os.Getenv("PORT")
	config.InitDatabase(databaseUrl)

	server := gin.New()
	myCorsConfig := cors.DefaultConfig()
	myCorsConfig.AllowAllOrigins = true
	myCorsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	server.Use(cors.New(myCorsConfig))
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello mantap mantap from github actions",
		})
	})
	routes.Routes(server)
	err := server.Run(port)

	if err != nil {
		log.Fatal("Error running server")
	}
}
