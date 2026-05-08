package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"kantin-mardira-api/internal/config"
	"kantin-mardira-api/internal/routes"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println(".env file not found, using system environment")
    }

    db := config.ConnectDatabase()

    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "success": true,
            "message": "Kantin Mardira API Running",
        })
    })
    routes.SetupRoutes(r, db)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}