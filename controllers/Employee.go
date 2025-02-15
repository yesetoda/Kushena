package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/infrastructures/password_services"
	"github.com/yesetoda/kushena/infrastructures/token_services"
	"github.com/yesetoda/kushena/models"
)

func (controller *ControllerImplementation) CreateEmployee(c *gin.Context) {
	var Employee *models.Employee
	s := ""
	if err := c.ShouldBindBodyWithJSON(&Employee); err != nil {

		s = "error in binding data"
		c.JSON(400, gin.H{"error": s})
		return
	}
	pass, err := password_services.HashPassword(Employee.Password)

	if err != nil {
		s = "error in hashing password"
		c.JSON(400, gin.H{"error": s})
		return
	}
	Employee.Password = pass
	Employee.Status = "out"
	Employee.Role = "Employee"
	err = controller.Usecases.CreateEmployee(Employee)
	if err != nil {
		s = "error in Adding employee " + Employee.Name
		c.JSON(400, gin.H{"error": s})
		return
	}
	s = "Employee " + Employee.Name + " Added successfully"
	c.JSON(200, gin.H{"message": s})
}

func (controller *ControllerImplementation) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	secret := os.Getenv("JWT_SECRET")
	s := ""

	token, err := controller.Usecases.Login(ctx, email, password, secret)
	if err != nil {
		s = "error in login"
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": s})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "message": "login successful"})
}

func (controller *ControllerImplementation) GetEmployeeById(c *gin.Context) {
	var employee *models.Employee
	id := c.Param("id")
	employee, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Employee  not found"})
		return
	}
	c.JSON(200, employee)
}

func (controller *ControllerImplementation) UpdateEmployee(c *gin.Context) {
	var Employee *models.Employee
	s := ""
	if err := c.ShouldBindBodyWithJSON(&Employee); err != nil {
		s = "error in binding data"
		c.JSON(400, gin.H{"error": s})
		return
	}
	err := controller.Usecases.UpdateEmployee(Employee)
	if err != nil {
		s = "error in updating employee " + Employee.Name
		c.JSON(400, gin.H{"error": s})
		return
	}
	s = "Employee " + Employee.Name + " updated successfully"
	c.JSON(200, gin.H{"message": s})
}

func (controller *ControllerImplementation) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	s := ""
	err := controller.Usecases.DeleteEmployee(id)
	if err != nil {
		s = "error in deleting employee " + id
		c.JSON(400, gin.H{"error": s})
	}
	s = "Employee " + id + " deleted successfully"
	c.JSON(200, gin.H{"message": s})
}

func (controller *ControllerImplementation) GetAllEmployees(c *gin.Context) {
	var employees []models.Employee
	employees, err := controller.Usecases.GetAllEmployees()
	if err != nil {
		c.JSON(404, gin.H{"error": "could not get all employees"})
		return
	}
	c.JSON(200, employees)
}

func (controller *ControllerImplementation) CheckIn(c *gin.Context) {
	claim, err := token_services.GetClaims(c)
	s := ""
	if err != nil {
		s = claim.Name + " you donot have the required authorization for this task."
		c.JSON(401, gin.H{"error": s})
		return
	}
	id := claim.ID.Hex()
	emp, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		s = claim.Name + " not found among Employees."
		c.JSON(404, gin.H{"error": s})
		return
	}
	if emp.Status == "in" {
		s = claim.Name + " you are already checked in."
		c.JSON(400, gin.H{"error": s})
		return
	}
	err = controller.Usecases.CheckIn(id)
	if err != nil {
		s = claim.Name + "  failed to check in."
		c.JSON(400, gin.H{"error": s})
		return

	}
	emp.Status = "in"
	err = controller.Usecases.UpdateEmployee(emp)
	if err != nil {
		s = claim.Name + "  failed to update your status."
		c.JSON(400, gin.H{"error": s})
		return
	}
	s = fmt.Sprintf("Hello %s", claim.Name)
	c.JSON(200, gin.H{"message": s})
}

func (controller *ControllerImplementation) CheckOut(c *gin.Context) {
	claim, err := token_services.GetClaims(c)
	s := ""
	if err != nil {
		s = claim.Name + " you donot have the required authorization for this task."
		c.JSON(401, gin.H{"error": s})
		return
	}
	id := claim.ID.Hex()
	emp, err := controller.Usecases.GetEmployeeById(id)
	if err != nil {
		s = claim.Name + " not found among Employees."
		c.JSON(404, gin.H{"error": s})
		return
	}
	if emp.Status == "out" {
		s = claim.Name + " you are already checked out."
		c.JSON(400, gin.H{"error": s})
		return
	}
	err = controller.Usecases.CheckOut(id)
	if err != nil {
		s = claim.Name + "  failed to check out."
		c.JSON(400, gin.H{"error": s})
		return
	}
	emp.Status = "out"
	err = controller.Usecases.UpdateEmployee(emp)
	if err != nil {
		s = claim.Name + "  failed to update your status."
		c.JSON(400, gin.H{"error": s})
		return
	}
	s = fmt.Sprintf("Goodbye %s", claim.Name)
	c.JSON(200, gin.H{"message": s})
}

func (controller *ControllerImplementation) Attendance(c *gin.Context) {
	claim, err := token_services.GetClaims(c)
    s := ""
    if err != nil {
        s = claim.Name + " you donot have the required authorization for this task."
        c.JSON(401, gin.H{"error": s})
        return
    }
    id := claim.ID.Hex()
    emp, err := controller.Usecases.GetEmployeeById(id)
    if err != nil {
        s = claim.Name + " not found among Employees."
        c.JSON(404, gin.H{"error": s})
        return
    }
	fmt.Println("employee", emp)
	fmt.Println("claims",claim)
    attendance, err := controller.Usecases.Attendance(id)
	if err != nil {
		s = claim.Name + "  failed to get attendance."
        c.JSON(400, gin.H{"error": s})
        return
	}
	c.JSON(200, attendance)
}
