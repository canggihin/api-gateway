package routes

import (
	"api-gateway/handlers"
	"api-gateway/repository"
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *fiber.App, mongoDB *mongo.Client) {
	userRepo := repository.NewUserService(mongoDB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := r.Group("/user")
	router.Post("/", userHandler.Register)
	router.Get("/", userHandler.UpdateStatus)
	router.Post("/login", userHandler.LoginClassic)
}
