package middleware

import (
	"fmt"
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("AdminAuth")
	fmt.Println("check1")
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "failed to login1", err.Error(), nil))
		return
	}
	fmt.Println("check2")
	adminID, err := ValidateToken2(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "failed to login2", err.Error(), nil))
		return
	}
	c.Set("adminID", adminID)
	c.Next()

}
