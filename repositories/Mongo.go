package repositories

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	EmployeeCollection   *mongo.Collection
	OrderCollection      *mongo.Collection
	FoodCollection       *mongo.Collection
	DrinkCollection      *mongo.Collection
	AttendanceCollection *mongo.Collection
}

func NewRepo() RepositoryInterface {
	mongoURI := os.Getenv("KUSHENA_URI")
	opts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("Kushena")
	EmployeeCollection := db.Collection("Employee")
	OrderCollection := db.Collection("Order")
	FoodCollection := db.Collection("Food")
	DrinkCollection := db.Collection("Drink")
	AttendanceCollection := db.Collection("Attendance")

	EmployeeIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1}, // Field to index (1 for ascending order)
			{Key: "phone_number", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = EmployeeCollection.Indexes().CreateOne(context.TODO(), EmployeeIndexModel)
	if err != nil {
		panic(err)
	}

	FoodIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // Field to index (1 for ascending order)
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = FoodCollection.Indexes().CreateOne(context.TODO(), FoodIndexModel)
	if err != nil {
		panic(err)
	}

	DrinkIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // Field to index (1 for ascending order)
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = DrinkCollection.Indexes().CreateOne(context.TODO(), DrinkIndexModel)
	if err != nil {
		panic(err)
	}

	// OrderIndexModel := mongo.IndexModel{
	// 	Keys: bson.M{
	// 		"name": 1, // Field to index (1 for ascending order)
	// 	},
	// 	Options: options.Index().SetUnique(true),
	// }
	// _, err = OrderCollection.Indexes().CreateOne(context.TODO(), OrderIndexModel)
	// if err != nil {
	// 	panic(err)
	// }

	return &MongoRepository{
		EmployeeCollection:   EmployeeCollection,
		OrderCollection:      OrderCollection,
		FoodCollection:       FoodCollection,
		DrinkCollection:      DrinkCollection,
		AttendanceCollection: AttendanceCollection,
	}

}
