package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/api", "http://localhost:8080", false, true)
	c.JSON(http.StatusOK, map[string]string{
		"message": "User logged out",
	})

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
