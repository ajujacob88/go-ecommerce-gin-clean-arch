package middleware

import (
	"errors"
	"fmt"
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
		fmt.Println("token string is", tokenString)
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Needs to login",
			})
			return
		}
		claims, err1 := ValidateToken(tokenString)
		fmt.Println("claims is", claims, "err is", err1)
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
			return []byte(config.GetJWTConfig()), nil
		},
	)
	fmt.Println("after parsing, err is ", err, "token is", token)

	if err != nil || !token.Valid {
		return claims, errors.New("not valid token")
	}
	//checking the expiry of the token
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return claims, errors.New("token expired re-login")
	}
	return claims, nil
}

// validate token for admin,, please merge this with the previous ValidateToken(for user) lateron
func ValidateToken2(tokenString string) (int, error) {
	//parses, validates, verifies the signature and returns the parsed token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		//retrieve the secret key which is stored in the env signing the string
		return []byte(config.GetJWTConfig()), nil
	})
	if err != nil {
		return 0, err
	}

	//extract the id claim from the token
	var parsedID interface{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		parsedID = claims["id"]
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, fmt.Errorf("token expired, re-login")
		}
	}
	//type assertion - It attempts to assert that the value stored in the parsedID variable is of type float64.
	value, ok := parsedID.(float64)
	if !ok {
		return 0, fmt.Errorf("expected a float value, but got %T, parsing id failed", parsedID)

	}

	id := int(value)
	return id, nil

}
