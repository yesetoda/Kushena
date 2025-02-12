package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	Addresses   []string           `json:"addresses" bson:"addresses"`
	Status      string             `json:"status" bson:"status"`
}
