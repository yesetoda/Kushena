package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodOrder struct {
	FoodId     primitive.ObjectID `json:"food_id" bson:"food_id "`
	FoodName   string             `json:"food_name" bson:"food_name"`
	Price      float64            `json:"price" bson:"price"`
	Quantity   float64            `json:"quantity" bson:"quantity"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
}

type DrinkOrder struct {
	DrinkId    primitive.ObjectID `json:"drink_id" bson:"drink_id"`
	DrinkName  string             `json:"drink_name" bson:"drink_name"`
	Price      float64            `json:"price" bson:"price"`
	Quantity   float64            `json:"quantity" bson:"quantity"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
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
