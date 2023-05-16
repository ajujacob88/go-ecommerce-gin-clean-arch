package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email string
	jwt.RegisteredClaims
}

// The middleware verifies the presence and validity of a token stored in a cookie and sets the user's email in the Gin context if the authorization is successful.
func AuthorizationMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(role + "-auth") //Inside the middleware, the function first tries to retrieve the JWT token from the cookie named role + "-token".
		//fmt.Println("token string is", tokenString)
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Needs to login",
			})
			return
		}
		claims, err1 := ValidateToken(tokenString)
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err1,
			})
			return
		}
		c.Set(role+"-email", claims.Email)
		c.Next()
	}
}

func ValidateToken(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTCofig()), nil
		},
	)
	if err != nil || !token.Valid {
		return claims, errors.New("not valid token")
	}
	//checking the expiry of the token
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return claims, errors.New("token expired re-login")
	}
	return claims, nil
}

/*
func LoginHandler(c *gin.Context) {
	// implement login logic here
	// user := c.PostForm("user")
	// pass := c.PostForm("pass")

	// // Throws Unauthorized error
	// if user != "john" || pass != "lark" {
	// 	return c.AbortWithStatus(http.StatusUnauthorized)
	// }

	// Create the Claims
	// claims := jwt.MapClaims{
	// 	"name":  "John Lark",
	// 	"admin": true,
	// 	"exp":   time.Now().Add(time.Hour * 72).Unix(),
	// }

	// Create token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	})

	ss, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}

*/
