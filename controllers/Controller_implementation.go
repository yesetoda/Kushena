package controllers

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/infrastructures/token_services"
	"github.com/yesetoda/kushena/models"
	"github.com/yesetoda/kushena/usecases"

)

type ControllerImplementation struct {
	Usecases usecases.UsecaseInterface
}

func NewController(usecases usecases.UsecaseInterface) ContollerInterface {
	return &ControllerImplementation{
		Usecases: usecases,
	}
}

func (controller *ControllerImplementation) Help(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Help for the Restaurant API"})
}

func (controller *ControllerImplementation) FoodTotalPrice(items []models.FoodOrder, price *float64, wg *sync.WaitGroup, mu *sync.Mutex, errChan chan error) {
	for i := range items {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			items[i].TotalPrice = items[i].Price * items[i].Quantity
			mu.Lock()
			*price += items[i].TotalPrice
			mu.Unlock()
		}(i)
	}
}

func (controller *ControllerImplementation) DrinkTotalPrice(items []models.DrinkOrder, price *float64, wg *sync.WaitGroup, mu *sync.Mutex, errChan chan error) {
	for i := range items {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			items[i].TotalPrice = items[i].Price * items[i].Quantity
			mu.Lock()
			*price += items[i].TotalPrice
			mu.Unlock()
		}(i)
	}
}
func (controller *ControllerImplementation) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claim, err := token_services.GetClaims(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	order.EmployeeId = claim.ID
	order.CreatedAt = time.Now().UTC()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var price float64
	errChan := make(chan error, len(order.Foods)+len(order.Drinks))

	controller.FoodTotalPrice(order.Foods, &price, &wg, &mu, errChan)
	controller.DrinkTotalPrice(order.Drinks, &price, &wg, &mu, errChan)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}
	}

	order.TotalPrice = price
	order.Status = "pending"

	if err := controller.Usecases.CreateOrder(order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Order created successfully"})
}

func (controller *ControllerImplementation) UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var price float64
	errChan := make(chan error, len(order.Foods)+len(order.Drinks))

	controller.FoodTotalPrice(order.Foods, &price, &wg, &mu, errChan)
	controller.DrinkTotalPrice(order.Drinks, &price, &wg, &mu, errChan)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}
	}

	order.TotalPrice = price

	if err := controller.Usecases.UpdateOrder(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Order updated successfully"})
}

func (controller *ControllerImplementation) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := controller.Usecases.DeleteOrder(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Order deleted successfully"})
}

func (controller *ControllerImplementation) GetOrderById(c *gin.Context) {
	id := c.Param("id")
	order, err := controller.Usecases.GetOrderById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(200, order)
}

func (controller *ControllerImplementation) GetAllOrders(c *gin.Context) {
	orders, err := controller.Usecases.GetAllOrders()
	if err != nil {
		c.JSON(404, gin.H{"error": "Orders not found"})
		return
	}
	c.JSON(200, orders)
}

func (controller *ControllerImplementation) CreateFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Usecases.CreateFood(food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Food created successfully"})
}

func (controller *ControllerImplementation) UpdateFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Usecases.UpdateFood(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Food updated successfully"})
}

func (controller *ControllerImplementation) DeleteFood(c *gin.Context) {
	id := c.Param("id")
	if err := controller.Usecases.DeleteFood(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Food deleted successfully"})
}

func (controller *ControllerImplementation) GetFoodById(c *gin.Context) {
	id := c.Param("id")
	food, err := controller.Usecases.GetFoodById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Food not found"})
		return
	}
	c.JSON(200, food)
}

func (controller *ControllerImplementation) GetAllFoods(c *gin.Context) {
	foods, err := controller.Usecases.GetAllFoods()
	if err != nil {
		c.JSON(404, gin.H{"error": "Foods not found"})
		return
	}
	c.JSON(200, foods)
}

func (controller *ControllerImplementation) CreateDrink(c *gin.Context) {
	var drink models.Drink
	if err := c.ShouldBindJSON(&drink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Usecases.CreateDrink(&drink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Drink created successfully"})
}

func (controller *ControllerImplementation) UpdateDrink(c *gin.Context) {
	var drink models.Drink
	if err := c.ShouldBindJSON(&drink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Usecases.UpdateDrink(&drink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Drink updated successfully"})
}

func (controller *ControllerImplementation) DeleteDrink(c *gin.Context) {
	id := c.Param("id")
	if err := controller.Usecases.DeleteDrink(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Drink deleted successfully"})
}

func (controller *ControllerImplementation) GetDrinkById(c *gin.Context) {
	id := c.Param("id")
	drink, err := controller.Usecases.GetDrinkById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Drink not found"})
		return
	}
	c.JSON(200, drink)
}

func (controller *ControllerImplementation) GetAllDrinks(c *gin.Context) {
	drinks, err := controller.Usecases.GetAllDrinks()
	if err != nil {
		c.JSON(404, gin.H{"error": "Drinks not found"})
		return
	}
	c.JSON(200, drinks)
}
