package infra

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var errMissingSecret = errors.New("JWT_SECRET não configurado")

// ponytail: single-purpose claims, no extra fields until needed
type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func jwtSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errMissingSecret
	}
	return []byte(secret), nil
}

func GenerateToken(userID string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
	})
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}
	return claims.UserID, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente"})
			c.Abort()
			return
		}

		userID, err := ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
