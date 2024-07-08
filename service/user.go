package service

import (
	"api-gateway/helpers"
	"api-gateway/models"
	"api-gateway/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	RegisterUser(ctx context.Context, data models.UserRegister) error
	UpdateActivate(ctx context.Context, status string, email string) error
	LoginClassic(ctx context.Context, data models.Login) (models.LoginResponse, error)
	UserInformation(ctx context.Context, username string) (models.UserRegister, error)
}

type userRepo struct {
	user repository.UserService
}

func NewUserService(r repository.UserService) *userRepo {
	return &userRepo{user: r}
}

func (s *userRepo) LoginClassic(ctx context.Context, data models.Login) (models.LoginResponse, error) {
	user, err := s.user.LoginClassic(ctx, data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.LoginResponse{}, &helpers.NotFoundError{Message: fmt.Sprintf("Unable to Find Username with %s", data.Username), MessageDev: err.Error()}
		}
	}

	if err := helpers.ComparePassword(user.Password, data.Password); err != nil {
		return models.LoginResponse{}, &helpers.UnauthorizedError{Message: "Invalid Password", MessageDev: err.Error()}
	}

	accessToken, refreshToken, err := helpers.EncodeWithStruct(&user)
	if err != nil {
		return models.LoginResponse{}, &helpers.InternalServerError{Message: "Failed to Create Token", MessageDev: err.Error()}
	}

	if err := s.user.UpdateRefreshToken(ctx, user.Username, refreshToken); err != nil {
		return models.LoginResponse{}, &helpers.InternalServerError{Message: "Failed to Update Refresh Token", MessageDev: err.Error()}
	}

	result := models.LoginResponse{
		Username:     user.Username,
		Status:       user.Status,
		Subscription: user.Subscription,
		ExpSubs:      user.ExpSubs,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return result, nil
}

func (s *userRepo) UserInformation(ctx context.Context, username string) (models.UserRegister, error) {
	user, err := s.user.UserInformation(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.UserRegister{}, &helpers.NotFoundError{Message: fmt.Sprintf("Unable to Find Username with %s", username), MessageDev: err.Error()}
		}
		return models.UserRegister{}, &helpers.InternalServerError{Message: "Failed to Get User Information", MessageDev: err.Error()}
	}

	return user, nil
}
func (s *userRepo) RegisterUser(ctx context.Context, data models.UserRegister) error {

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

	insert := &models.UserRegister{
		Username:    data.Username,
		Password:    hashedPassword,
		FullName:    strings.ToLower(data.FullName),
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
		Role:        strings.ToLower(data.Role),
		ExpActivate: time.Now().UTC().Add(5 * time.Minute),
	}

	if data.Role == "superadmin" {
		insert.Status = "active"
		insert.Subscription = true
	}

	if isUniqueUsername := s.user.CheckUniqueUsername(ctx, data.Username, data.Email); isUniqueUsername {
		return &helpers.BadRequestError{Message: "Username or Email already registered, Please Activate Now"}
	}

	if dataExist := s.user.CheckDataPending(ctx, data); dataExist {
		if err := s.user.UpdateExpActivate(ctx, data.Username, time.Now().UTC().Add(5*time.Minute)); err != nil {
			return &helpers.InternalServerError{Message: "Failed to Update Expired Activate", MessageDev: err.Error()}
		}
		return nil
	}

	if err := s.user.RegisterUser(ctx, *insert); err != nil {
		return &helpers.InternalServerError{Message: "Failed to register user, some error happened", MessageDev: err.Error()}
	}

	return nil
}

func (s *userRepo) UpdateActivate(ctx context.Context, status string, email string) error {
	if err := s.user.UpdateStatus(ctx, email, status); err != nil {
		return &helpers.InternalServerError{Message: "Failed to Activate Your Account, Please try Again Later", MessageDev: err.Error()}
	}
	return nil
}
