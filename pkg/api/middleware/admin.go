package middleware

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("AdminAuth")

	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "failed to login", err.Error(), nil))
		return
	}

	adminID, err := ValidateToken2(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "failed to login", err.Error(), nil))
		return
	}
	c.Set("adminID", adminID)
	c.Next()

}
