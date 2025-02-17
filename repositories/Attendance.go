package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yesetoda/kushena/models"
)

func (repo *MongoRepository) TakeAttendance(attendance_type, id string) error {
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
	return repo.TakeAttendance("in", id)
}

func (repo *MongoRepository) CheckOut(id string) error {
	return repo.TakeAttendance("out", id)
}

func (repo *MongoRepository) Attendance(id string) ([]models.Attendance, error) {
	eid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor, err := repo.AttendanceCollection.Find(context.TODO(), bson.M{"employee_id": eid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	attendances := []models.Attendance{}
	for cursor.Next(context.TODO()) {
		var attendance models.Attendance
		err := cursor.Decode(&attendance)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	fmt.Println("attendance for ", id)
	fmt.Println(attendances)

	return attendances, nil
}

func (repo *MongoRepository) CheckStatus( id string) (models.Attendance, error) {
	eid, err := primitive.ObjectIDFromHex(id)
	var attendance models.Attendance 
    if err != nil {
        return attendance, err
    }
	
	resp,err := repo.AttendanceCollection.CountDocuments(context.TODO(), bson.M{"employee_id": eid})
    if err != nil {
		return  attendance,err
    }
	if resp == 0 {
		return attendance, nil
	}
	err = repo.AttendanceCollection.FindOne(context.TODO(),
	bson.M{"employee_id": eid},
	options.FindOne().SetSort(bson.M{"time": -1})).Decode(&attendance)
	return attendance, err

}

func (repo *MongoRepository) TodaysWorkingTime(id string) (float64, error) {
	eid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
        return 0, err
    }
	startTime := time.Now().AddDate(0, 0, -1) // Last 1 day
	endTime := time.Now()

	cursor, err := repo.AttendanceCollection.Aggregate(context.TODO(), bson.A{
		bson.M{"$match": bson.M{
			"employee_id": eid,
			"time":        bson.M{"$gte": startTime, "$lte": endTime},
		}},
		bson.M{"$sort": bson.M{"time": 1}}, // Sort by time ascending
	})
	if err != nil {
		return 0, err
	}
	var records []models.Attendance
	if err = cursor.All(context.TODO(), &records); err != nil {
		return 0, err
	}

	var totalWorkDuration float64
	var lastCheckIn time.Time
	var hasCheckIn bool

	for _, record := range records {
		if record.Type == "in" {
			lastCheckIn = record.Time
			hasCheckIn = true
		} else if record.Type == "out" && hasCheckIn {
			totalWorkDuration += record.Time.Sub(lastCheckIn).Minutes()
			hasCheckIn = false
		}
	}
	if hasCheckIn { // If there's an ongoing shift at the end of the day, add it to the total work duration
        totalWorkDuration += time.Since(lastCheckIn).Minutes()
    }
	return totalWorkDuration, nil
}