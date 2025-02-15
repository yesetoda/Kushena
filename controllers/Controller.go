package controllers

import (
	"github.com/gin-gonic/gin"
)

type ContollerInterface interface {
	Help(ctx *gin.Context)

	CreateEmployee(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetEmployeeById(ctx *gin.Context)
	UpdateEmployee(ctx *gin.Context)
	DeleteEmployee(ctx *gin.Context)
	GetAllEmployees(ctx *gin.Context)

	CheckIn(ctx *gin.Context)
	CheckOut(ctx *gin.Context)
	Attendance(ctx *gin.Context)

	DailyReport(c *gin.Context)
	MonthlyReport(c *gin.Context)
	WeeklyReport(c *gin.Context)
	YearlyReport(c *gin.Context)

	CreateOrder(ctx *gin.Context)
	UpdateOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)
	GetOrderById(ctx *gin.Context)
	GetAllOrders(ctx *gin.Context)

	CreateFood(ctx *gin.Context)
	UpdateFood(ctx *gin.Context)
	DeleteFood(ctx *gin.Context)
	GetFoodById(ctx *gin.Context)
	GetAllFoods(ctx *gin.Context)

	CreateDrink(ctx *gin.Context)
	UpdateDrink(ctx *gin.Context)
	DeleteDrink(ctx *gin.Context)
	GetDrinkById(ctx *gin.Context)
	GetAllDrinks(ctx *gin.Context)
}
