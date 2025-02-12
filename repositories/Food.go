package repositories

import (
	"context"

	"github.com/yesetoda/Kushena/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *MongoRepository) CreateFood(food models.Food) error {
	food.Id = primitive.NewObjectID()
	_, err := repo.FoodCollection.InsertOne(context.Background(), food)
	return err

}
func (repo *MongoRepository) UpdateFood(food *models.Food) error {
	_, err := repo.FoodCollection.UpdateOne(context.Background(), bson.M{"_id": food.Id}, bson.M{"$set": food})
	return err

}
func (repo *MongoRepository) DeleteFood(id string) error {
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.FoodCollection.DeleteOne(context.Background(), bson.M{"_id": fid})
	return err

}
func (repo *MongoRepository) GetFoodById(id string) (*models.Food, error) {
	var food models.Food
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = repo.FoodCollection.FindOne(context.Background(), bson.M{"_id": fid}).Decode(&food)
	return &food, err

}
func (repo *MongoRepository) GetAllFoods() ([]models.Food, error) {
	var foods []models.Food
	cursor, err := repo.FoodCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var food models.Food
		err := cursor.Decode(&food)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil

}
