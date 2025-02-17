package usecases

import (
	"github.com/yesetoda/kushena/models"
	"github.com/yesetoda/kushena/repositories"
)

type UsecaseImplemented struct {
	Repo repositories.RepositoryInterface
}

func NewUsecase(repo repositories.RepositoryInterface) UsecaseInterface {
	return &UsecaseImplemented{
		Repo: repo,
	}
}

func (usecase *UsecaseImplemented) CreateOrder(order models.Order) error {
	err := usecase.Repo.CreateOrder(order)
	return err
}

func (usecase *UsecaseImplemented) UpdateOrder(order *models.Order) error {
	err := usecase.Repo.UpdateOrder(order)
	return err

}
func (usecase *UsecaseImplemented) DeleteOrder(id string) error {
	err := usecase.Repo.DeleteOrder(id)
	return err

}
func (usecase *UsecaseImplemented) GetOrderById(id string) (*models.Order, error) {
	var order *models.Order
	order, err := usecase.Repo.GetOrderById(id)
	return order, err

}
func (usecase *UsecaseImplemented) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	orders, err := usecase.Repo.GetAllOrders()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (usecase *UsecaseImplemented) GetAllMyOrders(id string) ([]models.Order, error) {
	var orders []models.Order
	orders, err := usecase.Repo.GetAllMyOrders(id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (usecase *UsecaseImplemented) CreateFood(food models.Food) error {
	err := usecase.Repo.CreateFood(food)
	return err

}
func (usecase *UsecaseImplemented) UpdateFood(food *models.Food) error {
	err := usecase.Repo.UpdateFood(food)
	return err

}
func (usecase *UsecaseImplemented) DeleteFood(id string) error {
	err := usecase.Repo.DeleteFood(id)
	return err

}
func (usecase *UsecaseImplemented) GetFoodById(id string) (*models.Food, error) {
	var food *models.Food
	food, err := usecase.Repo.GetFoodById(id)
	return food, err

}
func (usecase *UsecaseImplemented) GetAllFoods() ([]models.Food, error) {
	var foods []models.Food
	foods, err := usecase.Repo.GetAllFoods()
	if err != nil {
		return nil, err
	}

	return foods, nil

}

func (usecase *UsecaseImplemented) CreateDrink(drink *models.Drink) error {
	err := usecase.Repo.CreateDrink(drink)
	return err

}
func (usecase *UsecaseImplemented) UpdateDrink(drink *models.Drink) error {
	err := usecase.Repo.UpdateDrink(drink)
	return err

}
func (usecase *UsecaseImplemented) DeleteDrink(id string) error {
	err := usecase.Repo.DeleteDrink(id)
	return err

}
func (usecase *UsecaseImplemented) GetDrinkById(id string) (*models.Drink, error) {
	var drink *models.Drink
	drink, err := usecase.Repo.GetDrinkById(id)
	return drink, err

}
func (usecase *UsecaseImplemented) GetAllDrinks() ([]models.Drink, error) {
	var drinks []models.Drink
	drinks, err := usecase.Repo.GetAllDrinks()
	if err != nil {
		return nil, err
	}
	return drinks, nil

}
