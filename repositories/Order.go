package repositories

import (
	"context"

	"github.com/yesetoda/Kushena/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *MongoRepository) CreateOrder(order models.Order) error {
	order.Id = primitive.NewObjectID()
	_, err := repo.OrderCollection.InsertOne(context.Background(), order)
	return err
}

func (repo *MongoRepository) UpdateOrder(order *models.Order) error {
	_, err := repo.OrderCollection.UpdateOne(context.Background(), bson.M{"_id": order.Id}, bson.M{"$set": order})
	return err

}
func (repo *MongoRepository) DeleteOrder(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.OrderCollection.DeleteOne(context.Background(), bson.M{"_id": oid})
	return err

}
func (repo *MongoRepository) GetOrderById(id string) (*models.Order, error) {
	var order models.Order
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = repo.OrderCollection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&order)
	return &order, err

}
func (repo *MongoRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	cursor, err := repo.OrderCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var order models.Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
