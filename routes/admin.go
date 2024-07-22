package routes

import (
	"api-gateway/handlers"
	"api-gateway/repository"
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AdminRoutes(r *fiber.App, mongoDB *mongo.Client) {
	adminRepo := repository.NewAdminRepository(mongoDB)
	adminService := service.NewAdminService(adminRepo)
	adminHandler := handlers.NewAdminHandler(adminService)

	router := r.Group("/admin")
	router.Post("/", adminHandler.RegisterAdmin)
	router.Post("/login", adminHandler.LoginClassic)
}
