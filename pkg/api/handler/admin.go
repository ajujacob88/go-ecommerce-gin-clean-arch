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

//Create Admin - SuperAdmin can create a new admin from admin panel
// @Summary Create a new admin from admin panel
// @ID create-admin
// @Description Super admin can create a new admin from admin panel
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin_details body model.NewAdminInfo true "New Admin Details"

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

func (cr *AdminHandler) AdminLogin(c *gin.Context) {

}
