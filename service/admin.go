package service

import (
	"api-gateway/helpers"
	"api-gateway/models"
	"api-gateway/repository"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type AdminService interface {
	RegisterAdmin(ctx context.Context, data models.UserAdminRegister) error
	LoginClassic(ctx context.Context, data models.Login) (models.LoginResponse, error)
}

type adminService struct {
	adminRepo repository.AdminRepo
}

func NewAdminService(r repository.AdminRepo) *adminService {
	return &adminService{adminRepo: r}
}

func (s *adminService) RegisterAdmin(ctx context.Context, data models.UserAdminRegister) error {
	if !strings.Contains(data.Email, "@") || (!strings.HasSuffix(data.Email, "@gmail.com") && !strings.HasSuffix(data.Email, "@yahoo.com")) {
		return &helpers.BadRequestError{Message: "Invalid Email"}
	}

	if !strings.HasPrefix(data.PhoneNumber, "62") {
		return &helpers.BadRequestError{Message: "Phone number must start with 62"}
	}

	if strings.Contains(data.Username, " ") {
		return &helpers.BadRequestError{Message: "Username can't use Space"}
	}

	hashedPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return &helpers.BadRequestError{Message: "Failed to hash password", MessageDev: err.Error()}
	}

	insert := &models.UserAdminRegister{
		Username:    data.Username,
		Password:    hashedPassword,
		FullName:    strings.ToLower(data.FullName),
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
		Role:        strings.ToLower(data.Role),
	}
	if isUniqueUsername := s.adminRepo.CheckUniqueUsername(ctx, data.Username, data.Email, data.PhoneNumber); isUniqueUsername {
		return &helpers.BadRequestError{Message: "Username or Email already registered, Please Activate Now"}
	}

	if err := s.adminRepo.RegisterAdmin(ctx, *insert); err != nil {
		return &helpers.InternalServerError{Message: "Failed to register admin, some error happened", MessageDev: err.Error()}
	}
	return nil
}

func (s *adminService) LoginClassic(ctx context.Context, data models.Login) (models.LoginResponse, error) {
	user, err := s.adminRepo.LoginClassic(ctx, data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.LoginResponse{}, &helpers.NotFoundError{Message: fmt.Sprintf("Unable to Find Username with %s", data.Username), MessageDev: err.Error()}
		}
	}
	if err := helpers.ComparePassword(user.Password, data.Password); err != nil {
		return models.LoginResponse{}, &helpers.UnauthorizedError{Message: "Invalid Password", MessageDev: err.Error()}
	}

	accessToken, refreshToken, err := helpers.EncodeWithStructAdmin(&user)
	if err != nil {
		return models.LoginResponse{}, &helpers.InternalServerError{Message: "Failed to Create Token", MessageDev: err.Error()}
	}

	if err := s.adminRepo.UpdateRefreshToken(ctx, user.Username, refreshToken); err != nil {
		return models.LoginResponse{}, &helpers.InternalServerError{Message: "Failed to Update Refresh Token", MessageDev: err.Error()}
	}
	result := models.LoginResponse{
		Username:     user.Username,
		Status:       user.Status,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return result, nil
}
