package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegister struct {
	IDUser          primitive.ObjectID `json:"id_user" bson:"_id"`
	Username        string             `json:"username" bson:"username" binding:"required"`
	Password        string             `json:"password" bson:"password" binding:"required"`
	FullName        string             `json:"full_name" bson:"full_name" binding:"required"`
	PhoneNumber     string             `json:"phone_number" bson:"phone_number" binding:"required"`
	Email           string             `json:"email" bson:"email" binding:"required"`
	Status          string             `json:"status" bson:"status"`
	Subscription    bool               `json:"subscription" bson:"subscription"`
	ExpSubs         time.Time          `json:"exp_subs" bson:"exp_subs"`
	ExpActivate     time.Time          `json:"exp_activate" bson:"exp_activate"`
	RefreshToken    string             `json:"refresh_token" bson:"refresh_token"`
	ExpRefreshToken time.Time          `json:"exp_refresh_token" bson:"exp_refresh_token"`
}

type Login struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type LoginResponse struct {
	Status          string    `json:"status" bson:"status"`
	Username        string    `json:"username" bson:"username"`
	Subscription    bool      `json:"subscription" bson:"subscription"`
	ExpSubs         time.Time `json:"exp_subs" bson:"exp_subs"`
	RefreshToken    string    `json:"refresh_token" bson:"refresh_token"`
	ExpRefreshToken time.Time `json:"exp_refresh_token" bson:"exp_refresh_token"`
	AccessToken     string    `json:"access_token" bson:"access_token"`
}

type LoginPhoneNumber struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number" binding:"required"`
}
