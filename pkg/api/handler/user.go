package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/auth"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/verify"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

var user domain.Users

func (cr *UserHandler) UserSignUp(c *gin.Context) {
	//var user domain.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), user)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := cr.userUseCase.Signup(c, user); err != nil {
		response := res.ErrorResponse(400, "failed to signup", err.Error(), user)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//twilio otp check

	_, err := verify.TwilioSendOtp("+91" + user.Phone)
	if err != nil {
		response := res.ErrorResponse(400, "failed to generate otp", err.Error(), user)

		c.JSON(http.StatusBadRequest, response)
		return

	}
	response := res.SuccessResponse(200, "Success: Enter the otp", user)
	c.JSON(200, response)

	// response := res.SuccessResponse(200, "Account Created Successfully", user)
	// c.JSON(200, response)
}

// SIGN UP OTP VERIFICATION

func (cr *UserHandler) SignupOtpVerify(c *gin.Context) {
	var otp req.OTPVerify
	if err := c.BindJSON(&otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}
	if err := verify.TwilioVerifyOTP("+91"+user.Phone, otp.OTP); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Otp",
		})
		//fmt.Println("otp print 4")
		return
	}
	//fmt.Println("otp print 5")

	//user.VerifyStatus = true

	// Call the OTPVerifyStatusManage method to update the verification status
	//fmt.Println("user.id is", user.ID, "and user.Email is", user.Email)
	err := cr.userUseCase.OTPVerifyStatusManage(c.Request.Context(), user.Email, true)
	if err != nil {
		response := res.ErrorResponse(500, "Failed to update verification status", err.Error(), user)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "OTP validation OK..Account Created Successfully", user)
	c.JSON(200, response)
}

func (cr *UserHandler) UserLoginByEmail(c *gin.Context) {
	//receive data from request body
	var body req.UserLoginEmail

	if err := c.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "Input is invalid", err.Error(), nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//copy the body values to user
	var user domain.Users
	copier.Copy(&user, &body)

	// get user from database and check password in usecase
	user, err := cr.userUseCase.LoginWithEmail(c, user)
	if err != nil {
		response := res.ErrorResponse(400, "failed to login", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// generate token using jwt in map
	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to login", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.SetCookie("user-auth", tokenString["accessToken"], 60*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", tokenString["accessToken"])
	c.JSON(http.StatusOK, response)
}

// HomeHandler
func (cr *UserHandler) Homehandler(c *gin.Context) {
	email, ok := c.Get(("user-email"))
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
	}
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), email.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

/*

type Response struct {
	ID      uint   `copier:"must"`
	Name    string `copier:"must"`
	Surname string `copier:"must"`
}

// UserSignUp godoc
// @summary api for user to signup
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/users [get]
// @response 200 {object} []Response "OK"

func (cr *UserHandler) FindAll(c *gin.Context) {
	users, err := cr.userUseCase.FindAll(c.Request.Context())

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := []Response{}
		copier.Copy(&response, &users)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	user, err := cr.userUseCase.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) Save(c *gin.Context) {
	var user domain.Users

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := cr.userUseCase.Save(c.Request.Context(), user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) Delete(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse id",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := cr.userUseCase.FindByID(ctx, uint(id))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	if user == (domain.Users{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is not booking yet",
		})
		return
	}

	cr.userUseCase.Delete(ctx, user)

	c.JSON(http.StatusOK, gin.H{"message": "User is deleted successfully"})
}

*/
