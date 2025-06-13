package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})

	r.GET("/validate", func(c *gin.Context) {
		cookieToken, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		valid, err := verifyToken(cookieToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return
		}

		if valid {
			c.JSON(http.StatusOK, gin.H{})
		}
	})

	addr := "0.0.0.0:8080"

	err := r.Run(addr)
	if err != nil {
		fmt.Printf("failed to start application on %s: %v\n", addr, err)
		os.Exit(1)
	}
}

func verifyToken(tokenString string) (bool, error) {
	jwtSecret := os.Getenv(("JWT_SECRET"))

	if jwtSecret == "" {
		return false, fmt.Errorf("missing JWT_SECRET env variable")
	}

	mySigningKey := []byte(jwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing jwt using signing method")
		}
		return mySigningKey, nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
