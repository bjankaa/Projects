package main

import (
	"time"

	"exmaple.com/ulti-restapi/database"
	"exmaple.com/ulti-restapi/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDatabase()

	// setting up server with the help of gin
	server := gin.Default()

	// Configure CORS
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.AvaibleRoutes(server)

	// listening to incoming requets on port 3000
	server.Run(":3000")
}
