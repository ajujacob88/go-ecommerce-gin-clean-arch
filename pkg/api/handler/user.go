package handler

import (
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/verify"
	"github.com/gin-gonic/gin"
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

	response := res.SuccessResponse(200, "OTP validation OK..Account Created Successfully", user)
	c.JSON(200, response)
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
