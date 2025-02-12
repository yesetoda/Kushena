package models

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Claims struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}
