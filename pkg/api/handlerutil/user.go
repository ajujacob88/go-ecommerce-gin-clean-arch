package handlerutil

import (
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) (int, error) {
	ID := c.Value("userID") //the value of userID is taken from the middleware (look the middleware, there the userID is retrieved from the jwttoken)
	fmt.Println("in handlerutil user id is", ID)
	userID, err := strconv.Atoi(fmt.Sprintf("%v", ID))
	return userID, err
}
