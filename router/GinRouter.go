package router

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/controllers"
	"github.com/yesetoda/kushena/infrastructures/auth_services"

)

type GinRoute struct {
	Controller controllers.ContollerInterface
	Auth       auth_services.AuthController
}

func NewGinRoute(controller controllers.ContollerInterface, auth auth_services.AuthController) *GinRoute {
	return &GinRoute{
		Controller: controller,
		Auth:       auth,
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}


func (r *GinRoute) Run() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/", r.Controller.Help)
	router.POST("/employee/login", r.Controller.Login)

	router.POST("/checkin", r.Auth.AuthenticationMiddleware(), r.Controller.CheckIn)
	router.POST("/checkout", r.Auth.AuthenticationMiddleware(), r.Controller.CheckOut)

	report := router.Group("/report")
	report.Use(r.Auth.RoleMiddleware("Manager"))
	{
		report.GET("/daily", r.Controller.DailyReport)
		report.GET("/weekly", r.Controller.WeeklyReport)
		report.GET("/monthly", r.Controller.MonthlyReport)
		report.GET("/yearly", r.Controller.YearlyReport)
	}
	manager := router.Group("/manage")
	manager.Use(r.Auth.RoleMiddleware("Manager"))
	{
		manager.POST("/employee", r.Controller.CreateEmployee)
		manager.GET("/employee/:id", r.Controller.GetEmployeeById)
		manager.PATCH("/employee", r.Controller.UpdateEmployee)
		manager.DELETE("/employee/:id", r.Controller.DeleteEmployee)
		manager.GET("/employees", r.Controller.GetAllEmployees)
	}
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
	port := ":" + os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set", port)
		port = ":8080"
	}
	router.Run(port)
	return nil
}
