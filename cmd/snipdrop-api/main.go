package main

import (
	"log"
	"net/http"
	"os"
	"snipdrop-rest-api/internal/app/snipdrop-api/controller"
	"snipdrop-rest-api/internal/app/snipdrop-api/repository"
	"snipdrop-rest-api/internal/app/snipdrop-api/service"
	database "snipdrop-rest-api/internal/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync() // Flushes buffer, if any

	// Use zap logger to replace global loggers (like log.Println)
	zap.ReplaceGlobals(logger)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		zap.L().Info("No .env file found")
	}

	// Connect to the database
	database.ConnectDatabase() // Ensure this sets up the global `Db` variable in your package

	// Set up the Gin router
	router := gin.Default()

	// Setup repository, service, and controller
	snippetRepo := &repository.SnippetRepository{DB: database.Db} // Use the global Db variable
	snippetService := &service.SnippetService{Repo: snippetRepo}
	snippetController := &controller.SnippetController{Service: snippetService}

	// Define routes
	setupRoutes(router, snippetController)

	// Get the port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	if err := router.Run(":" + port); err != nil {
		zap.L().Fatal("Failed to run server", zap.Error(err))
	}
}

func setupRoutes(router *gin.Engine, snippetController *controller.SnippetController) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Snippet routes
	router.POST("/snippets", snippetController.CreateSnippet)
	router.GET("/snippets", snippetController.ListSnippets)
	router.GET("/snippets/:id", snippetController.GetSnippet)
	router.DELETE("/snippets/:id", snippetController.DeleteSnippet)
}
