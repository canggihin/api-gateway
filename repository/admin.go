package repository

import (
	"api-gateway/models"
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepo interface {
	RegisterAdmin(ctx context.Context, data models.UserAdminRegister) error
	CheckUniqueUsername(ctx context.Context, username string, email string, phoneNumber string) bool
	LoginClassic(ctx context.Context, data models.Login) (models.UserAdminRegister, error)
	UpdateRefreshToken(ctx context.Context, username string, refreshToken string) error
}

type adminRepo struct {
	adminColl *mongo.Collection
}

func NewAdminRepository(client *mongo.Client) *adminRepo {
	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("ADMIN_COLLECTION"))
	return &adminRepo{adminColl: collection}
}

func (r *adminRepo) RegisterAdmin(ctx context.Context, data models.UserAdminRegister) error {

	var insert *models.UserAdminRegister
	if data.Role == "admin" || data.Role == "superadmin" {
		insert = &models.UserAdminRegister{
			IDUser:          primitive.NewObjectID(),
			Username:        data.Username,
			Password:        data.Password,
			FullName:        data.FullName,
			PhoneNumber:     data.PhoneNumber,
			Email:           data.Email,
			Status:          "active",
			Role:            data.Role,
			RefreshToken:    "",
			ExpRefreshToken: data.ExpRefreshToken,
		}
	} else {
		return errors.New("role must be admin or superadmin")
	}

	_, err := r.adminColl.InsertOne(ctx, insert)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) CheckUniqueUsername(ctx context.Context, username string, email string, phoneNumber string) bool {
	var result models.UserAdminRegister

	filter := bson.M{
		"username":     username,
		"email":        email,
		"phone_number": phoneNumber,
	}

	if err := r.adminColl.FindOne(ctx, filter).Decode(&result); err == nil {
		return true
	}
	return false
}

func (r *adminRepo) LoginClassic(ctx context.Context, data models.Login) (models.UserAdminRegister, error) {
	var result models.UserAdminRegister

	filter := bson.M{
		"username": data.Username,
	}

	if err := r.adminColl.FindOne(ctx, filter).Decode(&result); err != nil {
		log.Println(err)
		return models.UserAdminRegister{}, err
	}

	return result, nil
}

func (r *adminRepo) UpdateRefreshToken(ctx context.Context, username string, refreshToken string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken, "exp_refresh_token": time.Now().UTC().Add(time.Hour)}}

	_, err := r.adminColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
