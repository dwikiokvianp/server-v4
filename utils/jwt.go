package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"
)

func GenerateJWT(email, role, id string) (string, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["id"] = id

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		_ = fmt.Errorf("something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	id := token.Claims.(jwt.MapClaims)["id"]
	fmt.Printf("id: %v ini id yaa dari auth", id)
	c.Set("id", id)
	c.Next()
}
