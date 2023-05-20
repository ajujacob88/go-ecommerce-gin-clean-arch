package handler

import (
	"fmt"
	"net/http"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/model"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}

}

// Create Admin - SuperAdmin can create a new admin from admin panel
// @Summary Create a new admin from admin panel
// @ID create-admin
// @Description Super admin can create a new admin from admin panel
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin_details body model.NewAdminInfo true "New Admin Details"
// @Success 201 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /admin/admins [post]
func (cr *AdminHandler) CreateAdmin(c *gin.Context) {
	var newAdminInfo model.NewAdminInfo
	if err := c.Bind(&newAdminInfo); err != nil {
		//The 422 status code is often used in API scenarios where clients submit data that fails validation, such as missing required fields, invalid data formats, or conflicting information.
		response := res.ErrorResponse(422, "unable to read the request body", err.Error(), nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	//finding out the admin id of the admin who is trying to create the new user., if the admin is super admin, then only he can able to create a new admin.
	adminID, err := handlerutil.GetAdminIdFromContext(c)
	fmt.Println("Admin ID is(for superuser check)", adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to fetch the admin ID", err.Error(), nil))
		return
	}
	//Now call the create admin method from admin usecase. The admin data will be saved to domain.admin after the succesful execution of the function
	newAdminOutput, err := cr.adminUseCase.CreateAdmin(c.Request.Context(), newAdminInfo, adminID)

	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to create the admin", err.Error(), nil))
		return
	}

	//if no error, then  201 status as new admin is created succesfully
	c.JSON(http.StatusCreated, res.SuccessResponse(201, "admin created successfully", newAdminOutput))

}

// AdminLogin
// @Summary Admin Login
// @ID admin-login
// @Description Admin Login
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin_credentials body model.AdminLoginInfo true "Admin Login Credentials"
// @Success 200 {object} res.Response
// @Failure 422 {object} res.Response
// @Failure 400 {object} res.Response
// @Router /admin/login [post]
func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	//receive the data from request body
	var body model.AdminLoginInfo
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, res.Response{StatusCode: 422, Message: "unable to process the request", Errors: err.Error(), Data: nil})
		return
	}
	//call the adminlogin method of the adminusecase to login as an admin
	tokenString, adminDataInModel, err := cr.adminUseCase.AdminLogin(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse(400, "failed to login", err.Error(), nil))
		return
	}
	c.SetSameSite(http.SameSiteLaxMode) //sets the SameSite attribute of the cookie to "Lax" mode. It is a security measure that helps protect against certain types of cross-site request forgery (CSRF) attacks.
	c.SetCookie("AdminAuth", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, res.SuccessResponse(200, "Succesfully Logged in", adminDataInModel))
}
