package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yesetoda/Kushena/controllers"
	"github.com/yesetoda/Kushena/infrastructures/auth_services"

)

type GinRoute struct {
	Controller controllers.ContollerInterface
	Auth       auth_services.AuthController
}

func NewGinRoute(controller controllers.ContollerInterface,auth auth_services.AuthController) *GinRoute {
	return &GinRoute{
		Controller: controller,
		Auth: auth,
	}
}

func (r *GinRoute) Run() error {
	router := gin.Default()
	router.GET("/", r.Controller.Help)

	// self actions
	router.POST("/checkin",r.Auth.AuthenticationMiddleware(), r.Controller.CheckIn)
	router.POST("/checkout",r.Auth.AuthenticationMiddleware(), r.Controller.CheckOut)

	router.GET("/report", r.Controller.Report)

	router.POST("/employee/login", r.Controller.Login)

	router.POST("/employee",r.Auth.RoleMiddleware("Manager"), r.Controller.CreateEmployee)
	router.GET("/employee/:id", r.Auth.RoleMiddleware("Manager"), r.Controller.GetEmployeeById)
	router.PATCH("/employee",r.Auth.RoleMiddleware("Manager"), r.Controller.UpdateEmployee)
	router.DELETE("/employee/:id",r.Auth.RoleMiddleware("Manager"), r.Controller.DeleteEmployee)
	router.GET("/employees",r.Auth.RoleMiddleware("Manager"), r.Controller.GetAllEmployees)

	actions := router.Group("/action")
	actions.Use(r.Auth.AuthenticationMiddleware())
	{
		actions.POST("/order", r.Controller.CreateOrder)
		actions.PATCH("/order", r.Controller.UpdateOrder)
		actions.DELETE("/order/:id", r.Controller.DeleteOrder)
		actions.GET("/order/:id", r.Controller.GetOrderById)
		actions.GET("/orders", r.Controller.GetAllOrders)
	
		actions.POST("/food", r.Controller.CreateFood)
		actions.PATCH("/food", r.Controller.UpdateFood)
		actions.DELETE("/food/:id", r.Controller.DeleteFood)
		actions.GET("/food/:id", r.Controller.GetFoodById)
		actions.GET("/foods", r.Controller.GetAllFoods)
	
		actions.POST("/drink", r.Controller.CreateDrink)
		actions.PATCH("/drink", r.Controller.UpdateDrink)
		actions.DELETE("/drink/:id", r.Controller.DeleteDrink)
		actions.GET("/drink/:id", r.Controller.GetDrinkById)
		actions.GET("/drinks", r.Controller.GetAllDrinks)
	}

	router.NoMethod(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Method not allowed"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route not found"})
	})

	router.Run(":8080")
	return nil
}
