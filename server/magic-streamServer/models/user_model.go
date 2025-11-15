package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)



type User struct{
	Id bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" `
	UserId string `bson:"user_id" json:"user_id"`
	FirstName string `bson:"frist_name" json:"first_name" validate:"required,min=3,max=100"`
	LastName string `bson:"last_name" json:"last_name" validate:"required,min=3,max=100"`
	Email string   `bson:"email" json:"email" validate:"required" `
	Password string  `bson:"password" json:"password" validate:"required, min=6"`
	Role string   `bson:"role" json:"role" validate:"one of=ADMIN USER"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	Token string `bson:"token" json:"token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	FavouriteGenre[] Genre `bson:"favourite_genre" json:"favourite_genre" validate:"required,dive"`


}
