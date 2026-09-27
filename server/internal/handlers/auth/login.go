package auth

import (
	"alilacream/ecom/internal/models"
	"alilacream/ecom/lib"
	"alilacream/ecom/store"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(s store.CustomerRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var customer models.Customer
		c.Header("Content-Type", "application/json")
		if err := c.BindJSON(&customer); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		existantCustomer, err := s.VerifyCrendiatials(c.Request.Context(), &customer)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		cookie, err := lib.GenerateJwt(*existantCustomer)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.SetCookie("jwt", cookie, 84600, "/api", "http://localhost:8080", false, true)
		c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Welcome back %s", existantCustomer.Username),
		})
	}
}
