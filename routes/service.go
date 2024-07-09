package routes

import (
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/repository"
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceRoutes(r *fiber.App, mongodb *mongo.Client) {
	serviceRepo := repository.NewServiceRepository(mongodb)
	userService := repository.NewUserService(mongodb)
	serviceService := service.NewService(serviceRepo, userService)
	serviceHandler := handlers.NewServiceHandler(serviceService)

	router := r.Group("/reg-service")
	router.Post("/", middleware.AuthMiddleware("superadmin"), serviceHandler.CreateService)
}

func Gateway(r *fiber.App, mongoDB *mongo.Client) {
	serviceRepo := repository.NewServiceRepository(mongoDB)
	userService := repository.NewUserService(mongoDB)
	service := service.NewService(serviceRepo, userService)
	serviceHandler := handlers.NewServiceHandler(service)

	r.All("/:service/:path", serviceHandler.GetService)
}
