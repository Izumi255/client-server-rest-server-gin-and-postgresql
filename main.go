package main

import (
	"log"

	"incident-tracker/database"
	_ "incident-tracker/docs"
	"incident-tracker/handlers"
	"incident-tracker/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Incident Tracker API
// @version 1.0
// @description REST API сервіс обліку технічних інцидентів та аварій на базі Gin та PostgreSQL.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введіть токен у форматі: Bearer <jwt_token>
func main() {
	database.InitDB()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterHandler)
			auth.POST("/login", handlers.LoginHandler)
		}

		incidents := v1.Group("/incidents")
		{
			incidents.GET("", handlers.GetIncidentsHandler)
			incidents.GET("/:id", handlers.GetIncidentByIDHandler)
		}

		protected := v1.Group("/incidents")
		protected.Use(middleware.AuthMiddleware(handlers.GetJWTSecret()))
		{
			protected.POST("", handlers.CreateIncidentHandler)
			protected.PATCH("/:id/status", handlers.UpdateIncidentStatusHandler)
			protected.DELETE("/:id", handlers.DeleteIncidentHandler)
		}
	}

	log.Println("[СЕРВЕР] Сервіс успішно запущено на http://localhost:8080")
	log.Println("[СЕРВЕР] Swagger UI доступний за адресою: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
