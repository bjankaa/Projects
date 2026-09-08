package routes

import (
	"exmaple.com/ulti-restapi/websockets"
	"github.com/gin-gonic/gin"
)

func AvaibleRoutes(server *gin.Engine) {

	manager := websockets.NewManager()

	server.POST("/auth", authenticationMethod)
	server.POST("/logout", logout)
	server.GET("/profile", getProfile)
	server.PUT("/change-email", changeEmail)
	server.PUT("/change-password", changePassword)
	server.GET("/games", getGames)

	//websocket
	server.GET("/game", manager.WSHandler)

	//test
	server.GET("/frontpage", getUsers)

	// Health check
	server.GET("/health", func(c *gin.Context) { c.Status(200) })
}
