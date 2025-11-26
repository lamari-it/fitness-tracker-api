package main

import (
	"lamari-fit-api/config"
	"lamari-fit-api/database"
	"lamari-fit-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	//database.InitializeDB()
	//database.SeedDatabase()

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
