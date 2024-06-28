package helpers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var statusCode int

	switch err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound

		if notFoundErr, ok := err.(*NotFoundError); ok {
			log.Println("--------Error Handler-------------")
			log.Printf("NotFoundError, Message is %s and Developer %s", err.Error(), notFoundErr.MessageDev)
			log.Println("---------------------")
		}

	case *BadRequestError:
		statusCode = http.StatusBadRequest

		if badRequestErr, ok := err.(*BadRequestError); ok {
			fmt.Println("--------Error Handler-------------")
			log.Printf("BadRequestError, Message is %s and Developer %s", err.Error(), badRequestErr.MessageDev)
			fmt.Println("---------------------")
		}

	case *InternalServerError:
		statusCode = http.StatusInternalServerError

		if internalServerErr, ok := err.(*InternalServerError); ok {
			log.Println("--------Error Handler-------------")
			log.Printf("InternalServerError, Message is %s and Developer %s", err.Error(), internalServerErr.MessageDev)
			log.Println("---------------------")
		}

	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized

		if unauthorizedErr, ok := err.(*UnauthorizedError); ok {
			log.Println("--------Error Handler-------------")
			log.Printf("UnauthorizedError, Message is %s and Developer %s", err.Error(), unauthorizedErr.MessageDev)
			log.Println("---------------------")
		}

	default:
		statusCode = http.StatusInternalServerError
		log.Println("--------Error Handler-------------")
		log.Printf("Unexpected Error, Message is %s", err.Error())
		log.Println("---------------------")
	}

	response := Response(ResponseParams{StatusCode: statusCode, Message: err.Error()})
	return c.Status(statusCode).JSON(response)
}
