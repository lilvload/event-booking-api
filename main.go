package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lilvload/event-booking-api/db"
	"github.com/lilvload/event-booking-api/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080") // localhost:8080
}
