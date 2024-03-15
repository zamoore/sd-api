package main

import (
	"log"
	"net/http"
	"snipdrop-rest-api/internal/app/snipdrop-api/controller"
	"snipdrop-rest-api/internal/app/snipdrop-api/middleware"
	"snipdrop-rest-api/internal/app/snipdrop-api/repository"
	"snipdrop-rest-api/internal/app/snipdrop-api/service"
	"snipdrop-rest-api/internal/pkg/config"
	database "snipdrop-rest-api/internal/pkg/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync() // Flushes buffer, if any

	// Load the global config
	config := config.LoadConfig()

	// Connect to the database
	database, err := database.ConnectDatabase(config)
	if err != nil {
		zap.L().Fatal("Failed to connect to database: %v", zap.Error(err))
	}
	defer database.Close()

	// Set up the Gin router
	router := gin.Default()

	// Setup repository, service, and controller
	snippetRepo := &repository.SnippetRepository{DB: database, Logger: logger}
	snippetService := &service.SnippetService{Repo: snippetRepo, Logger: logger}
	snippetController := &controller.SnippetController{Service: snippetService, Logger: logger}

	// Define routes
	setupRoutes(router, snippetController)

	// Start the server
	if err := router.Run(":" + config.Port); err != nil {
		zap.L().Fatal("Failed to run server", zap.Error(err))
	}
}

func setupRoutes(router *gin.Engine, snippetController *controller.SnippetController) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Snippet routes
	router.POST("/snippets", middleware.EnsureValidToken(), snippetController.CreateSnippet)
	router.GET("/snippets", snippetController.ListSnippets)
	router.GET("/snippets/:id", snippetController.GetSnippet)
	router.DELETE("/snippets/:id", snippetController.DeleteSnippet)
}
