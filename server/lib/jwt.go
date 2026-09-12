package lib

import (
	"time"

	"alilacream/ecom/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

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
