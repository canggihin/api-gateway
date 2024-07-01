package handlers

import (
	"api-gateway/helpers"
	"api-gateway/models"
	"api-gateway/service"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	user service.UserService
}

func NewUserHandler(user service.UserService) *userHandler {
	return &userHandler{user: user}
}

func (h *userHandler) LoginClassic(ctx *fiber.Ctx) error {
	var request models.Login

	if err := ctx.BodyParser(&request); err != nil {
		return &helpers.BadRequestError{Message: "Invalid Username or Password", MessageDev: err.Error()}
	}

	user, err := h.user.LoginClassic(context.Background(), request)
	if err != nil {
		return helpers.ErrorHandler(ctx, err)
	}
	result := helpers.Response(helpers.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Login Success",
		Data:       user,
	})

	return ctx.Status(http.StatusOK).JSON(result)
}
func (h *userHandler) Register(ctx *fiber.Ctx) error {
	var request models.UserRegister

	if err := ctx.BodyParser(&request); err != nil {
		return helpers.ErrorHandler(ctx, &helpers.BadRequestError{Message: "Failed to parse request body", MessageDev: err.Error()})
	}

	if err := h.user.RegisterUser(context.Background(), request); err != nil {
		return helpers.ErrorHandler(ctx, err)
	}

	res := helpers.Response(helpers.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
		Data: map[string]interface{}{
			"username":  request.Username,
			"full_name": request.FullName,
			"email":     request.Email,
			"status":    "pending",
		},
	})

	return ctx.Status(http.StatusCreated).JSON(res)
}

func (h *userHandler) UpdateStatus(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	status := ctx.Query("status")
	if err := h.user.UpdateActivate(context.Background(), status, email); err != nil {
		return helpers.ErrorHandler(ctx, err)
	}
	return nil
}
