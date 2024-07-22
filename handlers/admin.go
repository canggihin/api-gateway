package handlers

import (
	"api-gateway/helpers"
	"api-gateway/models"
	"api-gateway/service"
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type adminHandler struct {
	admin service.AdminService
}

func NewAdminHandler(r service.AdminService) *adminHandler {
	return &adminHandler{admin: r}
}

func (h *adminHandler) RegisterAdmin(ctx *fiber.Ctx) error {
	var request models.UserAdminRegister

	if err := ctx.BodyParser(&request); err != nil {
		return helpers.ErrorHandler(ctx, &helpers.BadRequestError{Message: "Failed to parse request body", MessageDev: err.Error()})
	}

	if err := h.admin.RegisterAdmin(context.Background(), request); err != nil {
		return helpers.ErrorHandler(ctx, err)
	}

	res := helpers.Response(helpers.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
		Data: map[string]interface{}{
			"username":  request.Username,
			"full_name": request.FullName,
			"email":     request.Email,
			"status":    "active",
		},
	})

	return ctx.Status(http.StatusCreated).JSON(res)
}

func (h *adminHandler) LoginClassic(ctx *fiber.Ctx) error {
	var request models.Login

	if err := ctx.BodyParser(&request); err != nil {
		return &helpers.BadRequestError{Message: "Invalid Username or Password", MessageDev: err.Error()}
	}

	user, err := h.admin.LoginClassic(context.Background(), request)
	if err != nil {
		log.Println("error in query :", err)
		return helpers.ErrorHandler(ctx, err)
	}
	result := helpers.Response(helpers.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Login Success",
		Data:       user,
	})

	return ctx.Status(http.StatusOK).JSON(result)
}
