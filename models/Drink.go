package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Drink struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
}
