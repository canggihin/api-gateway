package repository

import (
	"api-gateway/models"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	CheckUniqueUsername(ctx context.Context, username string, email string) bool
	RegisterUser(ctx context.Context, data models.UserRegister) error
	CheckDataPending(ctx context.Context, data models.UserRegister) bool
	UpdateExpActivate(ctx context.Context, username string, newExpActivate time.Time) error
	UpdateStatus(ctx context.Context, email string, status string) error
	LoginClassic(ctx context.Context, data models.Login) (models.UserRegister, error)
	LoginPhoneNUmber(ctx context.Context, data models.LoginPhoneNumber) bool
	UpdateRefreshToken(ctx context.Context, username string, refreshToken string) error
	UserInformation(ctx context.Context, username string) (models.UserRegister, error)
}

func NewUserService(client *mongo.Client) *serviceRepo {
	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USER_COLLECTION"))
	return &serviceRepo{MongoColl: collection}
}

func (r *serviceRepo) RegisterUser(ctx context.Context, data models.UserRegister) error {

	insert := &models.UserRegister{
		IDUser:       primitive.NewObjectID(),
		Username:     data.Username,
		Password:     data.Password,
		FullName:     data.FullName,
		PhoneNumber:  data.PhoneNumber,
		Email:        data.Email,
		Status:       "pending",
		Role:         data.Role,
		Subscription: false,
		ExpActivate:  data.ExpActivate,
	}
	_, err := r.MongoColl.InsertOne(ctx, insert)
	if err != nil {
		return err
	}
	return nil
}

func (r *serviceRepo) CheckUniqueUsername(ctx context.Context, username string, email string) bool {
	var result models.UserRegister

	log.Println(time.Now().UTC())
	// Buat filter untuk memeriksa username, status, dan ExpActivate
	filter := bson.M{
		"$and": []bson.M{
			{"username": username},
			{"email": email},
			{"$or": []bson.M{
				{"status": "active"},
				{"$and": []bson.M{
					{"status": "pending"},
					{"exp_activate": bson.M{"$gt": time.Now().UTC()}},
				}},
			}},
		},
	}

	err := r.MongoColl.FindOne(ctx, filter).Decode(&result)

	return err == nil
}

func (r *serviceRepo) UpdateExpActivate(ctx context.Context, username string, newExpActivate time.Time) error {
	filter := bson.M{"username": username}

	update := bson.M{
		"$set": bson.M{
			"exp_activate": newExpActivate,
		},
	}

	_, err := r.MongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *serviceRepo) CheckDataPending(ctx context.Context, data models.UserRegister) bool {
	var result models.UserRegister
	log.Println(time.Now())
	filter := bson.M{"$and": []bson.M{
		{"username": data.Username},
		{"email": data.Email},
		{"status": "pending"},
		{"exp_activate": bson.M{"$lt": time.Now().UTC()}},
	}}

	err := r.MongoColl.FindOne(ctx, filter).Decode(&result)

	return err == nil
}

func (r *serviceRepo) UpdateStatus(ctx context.Context, email string, status string) error {
	filter := bson.M{"email": email}

	newExpSubs := time.Now().UTC().AddDate(0, 0, 30)
	update := bson.M{
		"$set": bson.M{
			"status":   status,
			"exp_subs": newExpSubs,
		},
	}

	_, err := r.MongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *serviceRepo) LoginClassic(ctx context.Context, data models.Login) (models.UserRegister, error) {
	var result models.UserRegister

	filter := bson.M{
		"$and": []bson.M{
			{"username": data.Username},
			{"status": "active"},
		},
	}

	if err := r.MongoColl.FindOne(ctx, filter).Decode(&result); err != nil {
		log.Println(err)
		return models.UserRegister{}, err
	}

	return result, nil
}

func (r *serviceRepo) UserInformation(ctx context.Context, username string) (models.UserRegister, error) {
	var result models.UserRegister

	filter := bson.M{
		"$and": []bson.M{
			{"username": username},
			{"status": "active"},
		},
	}

	if err := r.MongoColl.FindOne(ctx, filter).Decode(&result); err != nil {
		log.Println(err)
		return models.UserRegister{}, err
	}

	return result, nil
}

func (r *serviceRepo) LoginPhoneNUmber(ctx context.Context, data models.LoginPhoneNumber) bool {
	var result models.LoginResponse

	filter := bson.M{
		"$and": []bson.M{
			{"phone_number": data.PhoneNumber},
			{"status": "active"},
		},
	}

	err := r.MongoColl.FindOne(ctx, filter).Decode(&result)

	return err == nil
}

func (r *serviceRepo) UpdateRefreshToken(ctx context.Context, username string, refreshToken string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken, "exp_refresh_token": time.Now().UTC().Add(time.Hour)}}

	_, err := r.MongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
