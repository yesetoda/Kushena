package repositories

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yesetoda/kushena/infrastructures/password_services"
	"github.com/yesetoda/kushena/models"
)

func (repo *MongoRepository) CreateEmployee(Employee *models.Employee) error {
	Employee.Id = primitive.NewObjectID()
	_, err := repo.EmployeeCollection.InsertOne(context.TODO(), Employee)
	return err
}

func (repo *MongoRepository) Login(email, password string) (*models.Employee, error) {
	var employee models.Employee
	err := repo.EmployeeCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&employee)
	if err != nil {
		return nil, err
	}
	err = password_services.CheckPasswordHash(password, employee.Password)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
func (repo *MongoRepository) GetEmployeeById(id string) (*models.Employee, error) {
	var employee models.Employee
	eid, err := primitive.ObjectIDFromHex(id)
	fmt.Println("employee id", id, eid)
	if err != nil {
		return nil, err
	}
	err = repo.EmployeeCollection.FindOne(context.TODO(), bson.M{"_id": eid}).Decode(&employee)
	return &employee, err
}
func (repo *MongoRepository) UpdateEmployee(Employee *models.Employee) error {
	res, err := repo.EmployeeCollection.UpdateOne(context.TODO(), bson.M{"_id": Employee.Id}, bson.M{"$set": Employee})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil

}
func (repo *MongoRepository) DeleteEmployee(id string) error {
	eid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	var emp models.Employee
	err = repo.EmployeeCollection.FindOne(context.TODO(), bson.M{"_id": eid}).Decode(&emp)
	if err != nil {
		return err
	}
	if emp.Role == "Manager" {
		return fmt.Errorf("manager cannot be deleted")
	}
	res, err := repo.EmployeeCollection.DeleteOne(context.TODO(), bson.M{"_id": eid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil

}
func (repo *MongoRepository) GetAllEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	cursor, err := repo.EmployeeCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var employee models.Employee
		err := cursor.Decode(&employee)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}
