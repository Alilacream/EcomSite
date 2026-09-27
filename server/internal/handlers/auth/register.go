package auth

import (
	"net/http"

	"alilacream/ecom/internal/models"
	"alilacream/ecom/lib"
	"alilacream/ecom/store"

	"github.com/gin-gonic/gin"
)

func Register(r store.CustomerRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newCustomer models.Customer
		c.Header("Content-Type", "application/json")
		if err := c.BindJSON(&newCustomer); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		registeredCustomer, err := r.Create(c.Request.Context(), &newCustomer)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if registeredCustomer.Username == newCustomer.Username {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		cookie, err := lib.GenerateJwt(*registeredCustomer)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// setting up the jwt cookie
		c.SetCookie("jwt", cookie, 84600, "/api", "http://localhost:8080", false, true)

		c.JSON(http.StatusOK, map[string]string{
			"response": "New user has registered!",
		})
	}
}
