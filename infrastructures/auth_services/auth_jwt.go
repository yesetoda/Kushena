package auth_services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yesetoda/kushena/infrastructures/token_services"
	"github.com/yesetoda/kushena/usecases"
)

type AuthController struct {
	Usecases usecases.UsecaseInterface
}

func NewAuthController(usecase usecases.UsecaseInterface) AuthController {
	return AuthController{
		Usecases: usecase}
}

func (ac *AuthController) AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "unexpected error"})
				c.Abort()
			}
		}()
		claims, err := token_services.GetClaims(c)
		if err != nil {
			fmt.Println("error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		_, err = ac.Usecases.GetEmployeeById(claims.ID.Hex())
		if err != nil {
			fmt.Println("error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
func (ac *AuthController) RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error occurred"})
				c.Abort()
			}
		}()

		fmt.Println("middleware", role)

		// Extract claims from token
		claims, err := token_services.GetClaims(c)
		fmt.Println("claims", claims, claims.ID, err)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Fetch employee details
		fmt.Println("trying to find the employee", claims.ID)
		employee, err := ac.Usecases.GetEmployeeById(claims.ID.Hex())
		fmt.Println("employee", employee, err)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}

		// Check if employee is nil
		if employee == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "employee record is missing"})
			return
		}

		// Check role authorization
		if employee.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
