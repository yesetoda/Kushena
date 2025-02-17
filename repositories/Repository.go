package repositories

import (

	"github.com/yesetoda/kushena/models"
)

type RepositoryInterface interface {
	CreateEmployee(Employee *models.Employee) error
	Login(email, password string) (*models.Employee, error)
	GetEmployeeById(id string) (*models.Employee, error)
	UpdateEmployee(Employee *models.Employee) error
	DeleteEmployee(id string) error
	GetAllEmployees() ([]models.Employee, error)

	CheckIn(id string) error
	CheckOut(id string) error
	Attendance(id string) ([]models.Attendance,error)
	CheckStatus(id string) (models.Attendance, error)
	TodaysWorkingTime(id string) (float64, error)

	Report(interval string) ( []byte, error)
	DailyReport()( []byte, error)
	WeeklyReport()( []byte, error)
	MonthlyReport()( []byte, error)
	YearlyReport()( []byte, error)

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
	UpdateDrink(Drink *models.Drink) error
	DeleteDrink(id string) error
	GetDrinkById(id string) (*models.Drink, error)
	GetAllDrinks() ([]models.Drink, error)
}
