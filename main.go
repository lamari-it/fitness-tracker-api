package main

import (
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fit-flow-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	database.AutoMigrate()
	database.SeedDatabase()

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}