package services

import (
	"errors"
	"net/http"

	"alilacream/ecom/config"
	"alilacream/ecom/internal/models"
	"alilacream/ecom/lib"
	"alilacream/ecom/store"

	"github.com/gin-gonic/gin"
)

// @CREATE
// user create a cart of the list of items he had
func CreatCart(r store.CacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item models.CartItem
		if err := c.BindJSON(&item); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// getting the user id from the jwt cookie
		userID, err := lib.GetUserFromToken(tokenString, config.GetEnv("SECRET_KEY"))
		// creating his own cart
		if ok, err := r.CreateItem(c.Request.Context(), userID, item); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
			// cart already created
		} else if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Cart already created"))
			return
		}
		c.JSON(http.StatusOK, map[string]string{
			"message": "created cart succesfully",
		})
	}
}
