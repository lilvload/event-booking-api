package routes

import (
	_ "github.com/lilvload/event-booking-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/lilvload/event-booking-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	v1 := server.Group("/api/v1")

	v1.GET("/events", getEvents)
	v1.GET("/events/:id", getEvent)

	authenticated := v1.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	v1.POST("/signup", signup)
	v1.POST("/login", login)

	// API documentation: http://localhost:8080/swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
