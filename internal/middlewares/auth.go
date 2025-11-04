package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(jwtSecret string, exists func(int) (bool, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := getClaimsFromHeader(c, jwtSecret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		userId, ok := claims["userId"].(float64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		ok, err = exists(int(userId))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userId", int(userId))
		c.Next()
	}
}

func jwtKeyFunc(jwtSecret string) func(token *jwt.Token) (interface{}, error) {
    return func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(jwtSecret), nil
    }
}

func getClaimsFromHeader(c *gin.Context, jwtSecret string) (jwt.MapClaims, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return jwt.MapClaims{}, errors.New("Authorization header is required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		return jwt.MapClaims{}, errors.New("Bearer token is required")
	}

	token, err := jwt.Parse(tokenString, jwtKeyFunc(jwtSecret))

	if err != nil || !token.Valid {
		return jwt.MapClaims{}, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || claims["type"] != "auth" {
		return jwt.MapClaims{}, errors.New("Invalid token")
	}

	return claims, nil
}
