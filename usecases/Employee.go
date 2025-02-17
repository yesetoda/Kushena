package usecases

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/infrastructures/email_services"
	"github.com/yesetoda/kushena/infrastructures/token_services"
	"github.com/yesetoda/kushena/models"

)

func (usecase UsecaseImplemented) CreateEmployee(Employee *models.Employee) error {
	err := usecase.Repo.CreateEmployee(Employee)
	email_services.SendEmail(Employee.Email, "Account Created", "Your Kushena Account has been created", "localhost:8080/employee/"+Employee.Id.Hex())
	return err
}

func (usecase UsecaseImplemented) Login(c *gin.Context, email, password, secret string) (string, error) {
	emp, err := usecase.Repo.Login(email, password)
	if err != nil {
		return "", err
	}
	token, err := token_services.GenerateToken(emp, password, secret)
	email_services.SendLoginAlertEmail(c, email, "localhost:8080/employee/"+emp.Id.Hex())

	return token, err
}
func (usecase UsecaseImplemented) GetEmployeeById(id string) (*models.Employee, error) {
	var employee *models.Employee
	fmt.Println("employee id", id)
	employee, err := usecase.Repo.GetEmployeeById(id)
	return employee, err
}
func (usecase *UsecaseImplemented) UpdateEmployee(Employee *models.Employee) error {
	err := usecase.Repo.UpdateEmployee(Employee)
	return err

}
func (usecase *UsecaseImplemented) DeleteEmployee(id string) error {
	err := usecase.Repo.DeleteEmployee(id)
	return err

}
func (usecase *UsecaseImplemented) GetAllEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	employees, err := usecase.Repo.GetAllEmployees()
	if err != nil {
		return nil, err
	}
	return employees, nil
}
func (usecase *UsecaseImplemented) CheckIn(id string) error {
	err := usecase.Repo.CheckIn(id)
	return err
}

func (usecase *UsecaseImplemented) CheckOut(id string) error {
	err := usecase.Repo.CheckOut(id)
	return err
}

func (usecase *UsecaseImplemented) Attendance(id string) ([]models.Attendance, error) {
	var attendance []models.Attendance
    attendance, err := usecase.Repo.Attendance(id)
    if err != nil {
        return nil, err
    }
    return attendance, nil
}
 func (usecase *UsecaseImplemented) CheckStatus(id string) (models.Attendance, error) {
	status, err := usecase.Repo.CheckStatus(id)
	return status, err
    
}

func (usecase *UsecaseImplemented) TodaysWorkingTime(id string) (float64, error) {
	duration, err := usecase.Repo.TodaysWorkingTime(id)
	return duration, err
}