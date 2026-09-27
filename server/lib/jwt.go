package lib

import (
	"errors"
	"fmt"
	"time"

	"alilacream/ecom/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	UserID string
	jwt.RegisteredClaims
}

func GenerateJwt(client models.Customer) (string, error) {
	claims := jwt.MapClaims{
		"ssd":   client.ID,
		"email": client.Email,
		"name":  "jwt",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign the token with you're secret
	return token.SignedString("secret")
}

func GetUserFromToken(tokenString string, secretkey string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secretkey), nil
	})
	if err != nil {
		return "", err
	}
	if claim, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claim.UserID, nil
	}
	return "", errors.New("Not a valid Token")
}
