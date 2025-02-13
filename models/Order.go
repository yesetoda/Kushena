package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodOrder struct {
	FoodId   primitive.ObjectID `json:"food_id" bson:"food_id "`
	Quantity float64            `json:"quantity" bson:"quantity"`
}

type DrinkOrder struct {
	DrinkId  primitive.ObjectID `json:"drink_id" bson:"drink_id"`
	Quantity float64            `json:"quantity" bson:"quantity"`
}
type Order struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	EmployeeId  primitive.ObjectID `json:"employee_id" bson:"employee_id"`
	TableNumber int                `json:"table_number" bson:"table_number"`
	Foods       []FoodOrder        `json:"foods" bson:"foods"`
	Drinks      []DrinkOrder       `json:"drinks" bson:"drinks"`

	TotalPrice float64   `json:"total_price" bson:"total_price"`
	Status     string    `json:"status" bson:"status"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}
