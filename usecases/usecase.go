package usecases

import (
	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/models"

)

type UsecaseInterface interface {
	CreateEmployee(Employee *models.Employee) error
	Login(c *gin.Context, email, password, secret string) (string, error)
	GetEmployeeById(id string) (*models.Employee, error)
	UpdateEmployee(Employee *models.Employee) error
	DeleteEmployee(id string) error
	GetAllEmployees() ([]models.Employee, error)

	CheckIn(id string) error
	CheckOut(id string) error
	Attendance(id string) ([]models.Attendance,error)
	CheckStatus(id string) (models.Attendance, error)

	DailyReport() ([]byte, error)
	WeeklyReport() ([]byte, error)
	MonthlyReport() ([]byte, error)
	YearlyReport() ([]byte, error)

	CreateOrder(order models.Order) error
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
	GetOrderById(id string) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
	GetAllMyOrders(id string) ([]models.Order, error)

	CreateFood(food models.Food) error
	UpdateFood(food *models.Food) error
	DeleteFood(id string) error
	GetFoodById(id string) (*models.Food, error)
	GetAllFoods() ([]models.Food, error)

	CreateDrink(drink *models.Drink) error
	UpdateDrink(drink *models.Drink) error
	DeleteDrink(id string) error
	GetDrinkById(id string) (*models.Drink, error)
	GetAllDrinks() ([]models.Drink, error)
}
