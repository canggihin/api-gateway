package handlers

import (
	"api-gateway/helpers"
	"api-gateway/middleware"
	"api-gateway/models"
	"api-gateway/service"
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type serviceHandler struct {
	services service.Service
}

func NewServiceHandler(s service.Service) *serviceHandler {
	return &serviceHandler{services: s}
}

func (h *serviceHandler) CreateService(ctx *fiber.Ctx) error {
	var request models.Service

	if err := ctx.BodyParser(&request); err != nil {
		return helpers.ErrorHandler(ctx, &helpers.BadRequestError{Message: "Failed to parse request body", MessageDev: err.Error()})

	}

	if err := h.services.Register(context.Background(), request); err != nil {
		return helpers.ErrorHandler(ctx, err)
	}

	res := helpers.Response(helpers.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Service created successfully",
	})
	return ctx.Status(http.StatusCreated).JSON(res)
}

func (h *serviceHandler) GetService(ctx *fiber.Ctx) error {
	serviceName := ctx.Params("service")
	path := ctx.Params("path")
	queryParams := ctx.Queries()

	var username string
	if queryParams["username"] != "" {
		username = queryParams["username"]
	}

	checkUsername, err := h.services.CheckUsername(username)
	if err != nil {
		return helpers.ErrorHandler(ctx, err)
	}
	if checkUsername == (models.UserRegister{}) {
		return helpers.ErrorHandler(ctx, &helpers.BadRequestError{Message: "Username not found"})
	}

	token := ctx.Get("x-authorization")
	resultToken, err := middleware.DecodeToken(token)
	if err != nil {
		return helpers.ErrorHandler(ctx, err)
	}

	if checkUsername.Username != resultToken["data"].(map[string]interface{})["username"] {
		return helpers.ErrorHandler(ctx, &helpers.BadRequestError{Message: "Username not match"})
	}

	data, err := h.services.GetService(context.Background(), serviceName)
	if err != nil {
		return helpers.ErrorHandler(ctx, err)
	}

	for _, header := range data.Headers {
		ctx.Request().Header.Set(header, ctx.Get(header))
	}

	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	ctx.Request().URI().SetHost(data.URL)
	ctx.Request().URI().SetPath(serviceName + "/" + path)
	ctx.Request().URI().SetScheme(data.Schema)

	fullPath := ctx.Request().URI().String()
	log.Println("Full Path:", fullPath)

	if err := proxy.Do(ctx, fullPath); err != nil {
		return helpers.ErrorHandler(ctx, &helpers.InternalServerError{Message: "Failed to proxy request", MessageDev: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
