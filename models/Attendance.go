package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Attendance struct {
	Id         primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	EmployeeID primitive.ObjectID    `bson:"employee_id" json:"employee_id"`
	Time       time.Time `bson:"time" json:"time"`
	Type       string    `bson:"type" json:"type"`
}
