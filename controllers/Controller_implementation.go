package controllers

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/Kushena/infrastructures/token_services"
	"github.com/yesetoda/Kushena/models"
	"github.com/yesetoda/Kushena/usecases"

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

func (controller *ControllerImplementation) GetFoodPrice(id string) (float64, error) {
	food, err := controller.Usecases.GetFoodById(id)
	if err != nil {
		return 0, err
	}
	return food.Price, nil

}
func (controller *ControllerImplementation) GetDrinkPrice(id string) (float64, error) {
	drink, err := controller.Usecases.GetDrinkById(id)
	if err != nil {
		return 0, err
	}
	return drink.Price, nil

}

func (controller *ControllerImplementation) CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	claim,err := token_services.GetClaims(c)
	if err != nil {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }
	order.EmployeeId = claim.ID
	order.CreatedAt = time.Now().UTC()

	var wg sync.WaitGroup
	var mu sync.Mutex // To protect price updates
	price := 0.0
	errChan := make(chan error, len(order.Foods)+len(order.Drinks)) // Collect errors

	// Fetch food prices in parallel
	for i := range order.Foods {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cost, err := controller.GetFoodPrice(order.Foods[i].FoodId.Hex())
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			price += cost * order.Foods[i].Quantity
			mu.Unlock()
		}(i)
	}

	// Fetch drink prices in parallel
	for i := range order.Drinks {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cost, err := controller.GetDrinkPrice(order.Drinks[i].DrinkId.Hex())
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			price += cost * order.Drinks[i].Quantity
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}
	}

	order.TotalPrice = price
	order.Status = "pending"

	err = controller.Usecases.CreateOrder(order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Order created successfully"})
}

func (controller *ControllerImplementation) UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindBodyWithJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // To protect price updates
	price := 0.0
	errChan := make(chan error, len(order.Foods)+len(order.Drinks)) // Collect errors

	// Fetch food prices in parallel
	for i := range order.Foods {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cost, err := controller.GetFoodPrice(order.Foods[i].FoodId.Hex())
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			price += cost * order.Foods[i].Quantity
			mu.Unlock()
		}(i)
	}

	// Fetch drink prices in parallel
	for i := range order.Drinks {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cost, err := controller.GetDrinkPrice(order.Drinks[i].DrinkId.Hex())
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			price += cost * order.Drinks[i].Quantity
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}
	}

	order.TotalPrice = price

	err := controller.Usecases.UpdateOrder(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Order updated successfully"})
}

func (controller *ControllerImplementation) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	err := controller.Usecases.DeleteOrder(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Order deleted successfully"})
}

func (controller *ControllerImplementation) GetOrderById(c *gin.Context) {
	var order *models.Order
	id := c.Param("id")
	order, err := controller.Usecases.GetOrderById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(200, order)
}

func (controller *ControllerImplementation) GetAllOrders(c *gin.Context) {
	var orders []models.Order
	orders, err := controller.Usecases.GetAllOrders()
	if err != nil {
		c.JSON(404, gin.H{"error": "Orders not found"})
		return
	}
	c.JSON(200, orders)
}

func (controller *ControllerImplementation) CreateFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindBodyWithJSON(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := controller.Usecases.CreateFood(food)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Food created successfully"})

}
func (controller *ControllerImplementation) UpdateFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindBodyWithJSON(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := controller.Usecases.UpdateFood(&food)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Food updated successfully"})

}
func (controller *ControllerImplementation) DeleteFood(c *gin.Context) {
	id := c.Param("id")
	err := controller.Usecases.DeleteFood(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Food deleted successfully"})

}
func (controller *ControllerImplementation) GetFoodById(c *gin.Context) {
	var food *models.Food
	id := c.Param("id")
	food, err := controller.Usecases.GetFoodById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Food not found"})
		return
	}
	c.JSON(200, food)

}
func (controller *ControllerImplementation) GetAllFoods(c *gin.Context) {
	var foods []models.Food
	foods, err := controller.Usecases.GetAllFoods()
	if err != nil {
		c.JSON(404, gin.H{"error": "Foods not found"})
		return
	}
	c.JSON(200, foods)

}

func (controller *ControllerImplementation) CreateDrink(c *gin.Context) {
	var drink models.Drink
	if err := c.ShouldBindBodyWithJSON(&drink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := controller.Usecases.CreateDrink(&drink)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Drink created successfully"})

}
func (controller *ControllerImplementation) UpdateDrink(c *gin.Context) {
	var drink models.Drink
	if err := c.ShouldBindBodyWithJSON(&drink); err != nil {
		// resp := models.Response{
		// 	HttpStatusCode: 400,
		// 	Message: "Error binding",
		// 	Data: nil,
		// }
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := controller.Usecases.UpdateDrink(&drink)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Drink updated successfully"})

}
func (controller *ControllerImplementation) DeleteDrink(c *gin.Context) {
	id := c.Param("id")
	err := controller.Usecases.DeleteDrink(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Drink deleted successfully"})

}
func (controller *ControllerImplementation) GetDrinkById(c *gin.Context) {
	var drink *models.Drink
	id := c.Param("id")

	drink, err := controller.Usecases.GetDrinkById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Drink not found"})
		return
	}
	c.JSON(200, drink)

}
func (controller *ControllerImplementation) GetAllDrinks(c *gin.Context) {
	var drinks []models.Drink
	drinks, err := controller.Usecases.GetAllDrinks()
	if err != nil {
		c.JSON(404, gin.H{"error": "Drinks not found"})
		return
	}
	c.JSON(200, drinks)

}
