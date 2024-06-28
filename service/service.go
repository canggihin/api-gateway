package service

import (
	"api-gateway/helpers"
	"api-gateway/models"
	"api-gateway/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Register(ctx context.Context, data models.Service) error
	GetService(ctx context.Context, service string) (*models.Service, error)
}

type service struct {
	serviceRepo repository.Service
}

func NewService(r repository.Service) *service {
	return &service{serviceRepo: r}
}

func (s *service) Register(ctx context.Context, data models.Service) error {
	serviceData, _ := s.serviceRepo.GetService(ctx, data.Name)

	checkUrl, _ := s.serviceRepo.CheckServiceUrl(ctx, data.URL)

	if serviceData != nil || checkUrl != nil {
		return &helpers.BadRequestError{Message: "Service already registered"}
	}

	if err := s.serviceRepo.RegisterService(ctx, data); err != nil {
		return &helpers.InternalServerError{Message: "Failed to register service, some error happened", MessageDev: err.Error()}
	}

	return nil
}

func (s *service) GetService(ctx context.Context, service string) (*models.Service, error) {
	result, err := s.serviceRepo.GetService(ctx, service)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &helpers.NotFoundError{Message: "Service not found", MessageDev: err.Error()}
		}
		return nil, &helpers.InternalServerError{Message: "There has been internal server error, please try again later", MessageDev: err.Error()}
	}

	return result, nil
}
