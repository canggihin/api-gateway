package main

import (
	"api-gateway/middleware"
	"api-gateway/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	db, err := middleware.InitMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
		return
	}

	configCors := cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-XSRF-TOKEN, X-CSRF-Token, Authorization, X-M2M-Origin, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers, Access-Control-Allow-Credentials, Origin, Accept, X-Requested-With, access-control-allow-origin, access-control-allow-methods, access-control-allow-headers",
		AllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		AllowCredentials: false,
	})

	app := fiber.New()
	app.Use(logger.New())
	app.Use(configCors)

	routes.ServiceRoutes(app, db)
	routes.UserRoutes(app, db)
	routes.Gateway(app, db)
	log.Fatal(app.Listen(":4000"))
}
