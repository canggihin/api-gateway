package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDBMiddleware(client *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("mongoClient", client)
		c.Next()
		return nil
	}
}

func InitMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("successfully create connection to MongoDB")
	return client, nil
}
