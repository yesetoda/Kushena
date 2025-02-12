package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/Kushena/infrastructures/password_services"
	"github.com/yesetoda/Kushena/infrastructures/token_services"
	"github.com/yesetoda/Kushena/models"
)

func (controller *ControllerImplementation) CreateEmployee(c *gin.Context) {
	var Employee *models.Employee
	if err := c.ShouldBindBodyWithJSON(&Employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	pass, err := password_services.HashPassword(Employee.Password)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	Employee.Password = pass
	Employee.Status = "out"
	err = controller.Usecases.CreateEmployee(Employee)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Employee created successfully"})
}

func (controller *ControllerImplementation) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	secret := os.Getenv("JWT_SECRET")

	token, err := controller.Usecases.Login(ctx, email, password, secret)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "could not login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "message": "login successful"})
}

func (controller *ControllerImplementation) GetEmployeeById(c *gin.Context) {
	var employee *models.Employee
	id := c.Param("id")
	employee, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(200, employee)
}

func (controller *ControllerImplementation) UpdateEmployee(c *gin.Context) {
	var Employee *models.Employee
	if err := c.ShouldBindBodyWithJSON(&Employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := controller.Usecases.UpdateEmployee(Employee)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func (controller *ControllerImplementation) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	err := controller.Usecases.DeleteEmployee(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Employee deleted successfully"})
}

func (controller *ControllerImplementation) GetAllEmployees(c *gin.Context) {
	var employees []models.Employee
	employees, err := controller.Usecases.GetAllEmployees()
	if err != nil {
		c.JSON(404, gin.H{"error": "Employees not found"})
		return
	}
	c.JSON(200, employees)
}

func (controller *ControllerImplementation) CheckIn(c *gin.Context) {
	claim, err := token_services.GetClaims(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	id := claim.ID.Hex()
	emp, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}
	if emp.Status == "in" {
		c.JSON(400, gin.H{"error": "Employee is already checked in"})
		return
	}
	err = controller.Usecases.CheckIn(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	emp.Status = "in"
	err = controller.Usecases.UpdateEmployee(emp)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Employee Checked In"})
}

func (controller *ControllerImplementation) CheckOut(c *gin.Context) {
	claim, err := token_services.GetClaims(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	id := claim.ID.Hex()
	emp, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}
	if emp.Status == "out" {
		c.JSON(400, gin.H{"error": "Employee is already checked out"})
		return
	}
	err = controller.Usecases.CheckOut(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	emp.Status = "out"
	err = controller.Usecases.UpdateEmployee(emp)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Employee Checked Out"})
}
