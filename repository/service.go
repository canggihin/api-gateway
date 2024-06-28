package repository

import (
	"api-gateway/models"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	RegisterService(ctx context.Context, data models.Service) error
	GetService(ctx context.Context, service string) (*models.Service, error)
	CheckServiceUrl(ctx context.Context, url string) (*models.Service, error)
}

type serviceRepo struct {
	MongoColl *mongo.Collection
}

func NewServiceRepository(client *mongo.Client) *serviceRepo {
	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("SERVICE_COLLECTION"))
	return &serviceRepo{MongoColl: collection}
}

func (r *serviceRepo) RegisterService(ctx context.Context, data models.Service) error {
	_, err := r.MongoColl.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *serviceRepo) GetService(ctx context.Context, service string) (*models.Service, error) {
	var result models.Service
	err := r.MongoColl.FindOne(ctx, bson.M{"name": service}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *serviceRepo) CheckServiceUrl(ctx context.Context, url string) (*models.Service, error) {
	var result models.Service
	err := r.MongoColl.FindOne(ctx, bson.M{"url": url}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
