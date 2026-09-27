package lib

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), err
}

func CheckPasswordHash(pass, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
