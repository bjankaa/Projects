package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exmaple.com/ulti-restapi/database"
	"exmaple.com/ulti-restapi/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	// setting up server with the help of gin
	server := gin.Default()
	routes.AvaibleRoutes(server)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: server,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	// Close database connections
	if database.Database != nil {
		_ = database.Database.Close()
	}

	log.Println("Server exiting")
}
