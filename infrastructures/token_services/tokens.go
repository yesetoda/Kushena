package token_services

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/yesetoda/Kushena/models"
)

func GenerateToken(employee *models.Employee, password, jwtSecret string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)); err != nil {
		return "", errors.New("invalid user name or password")
	}

	accessToken, err := createJWTToken(employee, jwtSecret, 120*time.Minute)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func createJWTToken(employee *models.Employee, jwtSecret string, duration time.Duration) (string, error) {

	expirationTime := time.Now().Add(duration)
	claims := &models.Claims{
		ID:          employee.Id,
		Name:        employee.Name,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		Role:        employee.Role,
		Addresses:   employee.Addresses,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetClaims(c *gin.Context) (*models.Claims, error) {

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return &models.Claims{}, errors.New("missing authorization header")
	}

	TokenString := strings.Split(authHeader, " ")
	if len(TokenString) != 2 || TokenString[0] != "Bearer" {
		return &models.Claims{}, errors.New("invalid token format")
	}
	tokenString := TokenString[1]

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return &models.Claims{}, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, err
	}
	return &models.Claims{}, errors.New("invalid token")
}

func GenerateConfirmationToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	confirmationToken := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range confirmationToken {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		confirmationToken[i] = charset[num.Int64()]
	}

	return string(confirmationToken), nil
}
