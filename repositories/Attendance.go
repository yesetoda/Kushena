package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yesetoda/Kushena/models"
)

func (repo *MongoRepository) Attendance(attendance_type, id string) error {
	fmt.Println(id, "checking ", attendance_type)
	eid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	attendance := models.Attendance{
		Id:         primitive.NewObjectID(),
		EmployeeID: eid,
		Time:       time.Now().UTC(),
		Type:       attendance_type,
	}
	_, err = repo.AttendanceCollection.InsertOne(context.TODO(), attendance)
	return err
}

func (repo *MongoRepository) CheckIn(id string) error {
	return repo.Attendance("in", id)
}

func (repo *MongoRepository) CheckOut(id string) error {
	return repo.Attendance("out", id)
}
