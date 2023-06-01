package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/auth"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	userUseCase services.UserUseCase
	otpUseCase  services.OTPUseCase
}

func NewUserHandler(usecase services.UserUseCase, otpusecase services.OTPUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
		otpUseCase:  otpusecase,
	}
}

//var user domain.Users

// @title Ecommerce REST API
// @version 1.0
// @description Ecommerce REST API built using Go, PSQL, REST API following Clean Architecture.

// @contact
// name: Aju Jacob
// url: https://github.com/ajujacob88
// email: ajujacob88@gmail.com

// @license
// name: MIT
// url: https://opensource.org/licenses/MIT

// @host localhost:3000

// @Basepath /
// @Accept json
// @Produce json
// @Router / [get]

// UserSignup
// @Summary api for Signup a new user
// @ID Signup-user
// @Description Create a new user with the specified details.
// @Tags Users Signup
// @Accept json
// @Produce json
// @Param user_details body model.NewUserInfo true "New user Details"
// @Success 200 {object} res.Response
// @Failure 500 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /user/signup [post]
func (cr *UserHandler) UserSignUp(c *gin.Context) {
	//var user domain.Users
	var newUserInfo model.NewUserInfo
	if err := c.BindJSON(&newUserInfo); err != nil {
		response := res.ErrorResponse(422, "unable to read the request body", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userDetails, err := cr.userUseCase.UserSignUp(c.Request.Context(), newUserInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to create user", err.Error(), nil))
	}

	//twilio otp send

	responseID, err := cr.otpUseCase.TwilioSendOtp(c.Request.Context(), "+91"+userDetails.Phone)
	if err != nil {
		response := res.ErrorResponse(500, "failed to generate otp", err.Error(), nil)

		c.JSON(http.StatusInternalServerError, response)
		return

	}
	response := res.SuccessResponse(200, "Success: Enter the otp and the response id", responseID)
	c.JSON(http.StatusOK, response)

}

// SIGN UP OTP VERIFICATION
// SignupOtpVerify
// @Summary signup otp verification
// @ID Signup-otpverify-user
// @Description verify the otp of a user.
// @Tags Users otp verify
// @Accept json
// @Produce json
// @Param otpverify body model.OTPVerify true "OTP verification details"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /user/signup/otp/verify [post]
func (cr *UserHandler) SignupOtpVerify(c *gin.Context) {
	//var user domain.Users
	var otpverify model.OTPVerify
	if err := c.BindJSON(&otpverify); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to read the request body", err.Error(), nil))
		return
	}
	otpsession, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), otpverify)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "Invalid Otp", err.Error(), nil))
		return
	}

	// Call the OTPVerifyStatusManage method to update the verification status
	err = cr.userUseCase.OTPVerifyStatusManage(c.Request.Context(), otpsession)
	if err != nil {
		response := res.ErrorResponse(500, "Failed to update verification status", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "OTP validation Successfull..Account Created Successfully", nil)
	c.JSON(200, response)
}

// UserLogin
// @Summary User Login
// @ID user-login
// @Description user Login
// @Tags user
// @Accept json
// @Produce json
// @Param user_credentials body req.UserLoginEmail true "user Login Credentials"
// @Success 200 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /user/login/email [post]
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
	user, err := cr.userUseCase.LoginWithEmail(c, body)
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

	c.SetCookie("UserAuth", tokenString["accessToken"], 60*60, "", "", false, true)

	//response := res.SuccessResponse(200, "successfully logged in", tokenString["accessToken"])
	response := res.SuccessResponse(200, "successfully logged in", user.FirstName)
	c.JSON(http.StatusOK, response)
}

// Userhome
// @Summary User Home
// @ID user-home
// @Description user home
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /user/home [get]
func (cr *UserHandler) Homehandler(c *gin.Context) {
	email, ok := c.Get(("user-email"))
	fmt.Println("email in homehandler is", email)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user error1",
		})
	}
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), email.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user error2",
		})
		return
	}
	fmt.Println("user is", user)
	c.JSON(http.StatusOK, res.SuccessResponse(200, "user  home", nil))

}

// UserLogout
// @Summary User_Logout
// @ID user-logout
// @Description logout an logged-in user from the site
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/logout [get]
func (cr *UserHandler) LogoutHandler(c *gin.Context) {
	//c.SetCookie("user-token", "", -1, "/", "localhost", false, true)

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") //indicates to the client that it should not cache any response data and should always revalidate it with the server
	c.SetSameSite(http.SameSiteLaxMode)                                           //sets the SameSite cookie attribute to "Lax" for the response. This attribute restricts the scope of cookies and helps prevent cross-site request forgery attacks
	c.SetCookie("UserAuth", "", -1, "", "", false, true)                          //Immediately by setting the maxAge to -1, and marks the cookie as secure and HTTP-only
	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully Logged-Out"))

}

// AddAddress
// @Summary User can add the user address
// @ID add-address
// @Description Add address
// @Tags user
// @Accept json
// @Produce json
// @Param user_address body model.UserAddressInput true "User address"
// @Success 201 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /user/addresses/ [post]
func (cr *UserHandler) AddAddress(c *gin.Context) {
	var userAddressInput model.UserAddressInput
	if err := c.Bind(&userAddressInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to read the request body", err.Error(), nil))
		return
	}

	userID, err := handlerutil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "unable to fetch user id from context", err.Error(), nil))
		return
	}

	address, err := cr.userUseCase.AddAddress(c.Request.Context(), userAddressInput, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to add the address", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Succesfully added the address", address))

}

// UpdateAddress
// @Summary User can edit the user address
// @ID update-address
// @Description Update address
// @Tags user
// @Accept json
// @Produce json
// @Param user_address body model.UserAddressInput true "User address"
// @Param address_id path string true "address id"
// @Success 201 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /user/addresses/edit/{address_id} [patch]
func (cr *UserHandler) UpdateAddress(c *gin.Context) {
	var userAddressInput model.UserAddressInput

	if err := c.Bind(&userAddressInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "unable to fetch the address body", err.Error(), nil))
		return
	}

	//no need since i am going to update based on address id
	// userID, err := handlerutil.GetUserIdFromContext(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, res.ErrorResponse(400, "unable to fetch user id from context", err.Error(), nil))
	// 	return
	// }

	paramsID := c.Param("address_id")
	addressID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the address id", err.Error(), nil))
		return
	}

	updatedAddress, err := cr.userUseCase.UpdateAddress(c.Request.Context(), userAddressInput, addressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to update the address", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Succesfully updated the address", updatedAddress))

}

// DeleteAddress
// @Summary User can delete the user address
// @ID delete-address
// @Description Delete address
// @Tags user
// @Accept json
// @Produce json
// @Param address_id path string true "address id"
// @Success 201 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /user/addresses/{address_id} [delete]
func (cr *UserHandler) DeleteAddress(c *gin.Context) {
	paramsID := c.Param("address_id")
	addressID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.ErrorResponse(422, "failed to parse the address id", err.Error(), nil))
		return
	}

	if err = cr.userUseCase.DeleteAddress(c.Request.Context(), addressID); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to delete the address", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "Succesfully deleted the address", nil))

}
