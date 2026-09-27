package lib

import (
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	var status bool
	switch status {
	case len(password) < 8:
		return false
	case !strings.ContainsAny(password, "*/-?!_ç'"):
		return false
	case !strings.ContainsAny(password, "1234567890"):
		return false
	}

	return true
}

func IsValidUsername(username string) bool {
	var status bool
	switch status {
	case len(username) < 8:
		return false
	case strings.ContainsAny(username, "</>:?!*"):
		return false
	}
	return true
}
