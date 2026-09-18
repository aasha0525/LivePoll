package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"livepoll/db"
	"livepoll/handlers"
	"livepoll/middleware"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to MongoDB and Redis
	db.ConnectMongo()
	db.ConnectRedis()

	// Create Gin router
	r := gin.Default()

	// CORS configuration
	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		// Authentication
		api.POST("/signup", handlers.Signup)
		api.POST("/login", handlers.Login)

		// Public poll routes
		api.GET("/polls/:id", handlers.GetPoll)
		api.POST("/polls/:id/vote", handlers.Vote)
		api.GET("/ws/polls/:id", handlers.PollWebSocket)

		// Protected routes
		authorized := api.Group("/")
		authorized.Use(middleware.AuthRequired())
		{
			authorized.POST("/polls", handlers.CreatePoll)
			authorized.GET("/polls", handlers.MyPolls)
		}
	}

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 LivePoll backend running on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
