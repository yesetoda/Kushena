package repositories

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yesetoda/Kushena/models"
)

func (repo *MongoRepository) CreateDrink(drink *models.Drink) error {
	drink.Id = primitive.NewObjectID()
	_, err := repo.DrinkCollection.InsertOne(context.Background(), drink)
	return err

}
func (repo *MongoRepository) UpdateDrink(drink *models.Drink) error {
	res, err := repo.DrinkCollection.UpdateOne(context.Background(), bson.M{"_id": drink.Id}, bson.M{"$set": drink})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("drink not found")
	}
	return err

}
func (repo *MongoRepository) DeleteDrink(id string) error {
	did, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := repo.DrinkCollection.DeleteOne(context.Background(), bson.M{"_id": did})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("drink not found")
	}

	return err

}
func (repo *MongoRepository) GetDrinkById(id string) (*models.Drink, error) {
	var drink models.Drink
	did, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = repo.DrinkCollection.FindOne(context.Background(), bson.M{"_id": did}).Decode(&drink)
	return &drink, err

}
func (repo *MongoRepository) GetAllDrinks() ([]models.Drink, error) {
	var drinks []models.Drink
	cursor, err := repo.DrinkCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var drink models.Drink
		err := cursor.Decode(&drink)
		if err != nil {
			return nil, err
		}
		drinks = append(drinks, drink)
	}
	return drinks, nil

}
