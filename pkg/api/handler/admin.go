package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handlerutil"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	services "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/usecase/interface"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/utils"
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
// @Param admin_details body request.NewAdminInfo true "New Admin Details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /admin/admins [post]
func (cr *AdminHandler) CreateAdmin(c *gin.Context) {
	var newAdminInfo request.NewAdminInfo
	if err := c.Bind(&newAdminInfo); err != nil {
		//The 422 status code is often used in API scenarios where clients submit data that fails validation, such as missing required fields, invalid data formats, or conflicting information.
		response := response.ErrorResponse(422, "unable to read the request body", err.Error(), nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//finding out the admin id of the admin who is trying to create the new user., if the admin is super admin, then only he can able to create a new admin.
	adminID, err := handlerutil.GetAdminIdFromContext(c)
	fmt.Println("Admin ID is(for superuser check)", adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to fetch the admin ID", err.Error(), nil))
		return
	}
	//Now call the create admin method from admin usecase. The admin data will be saved to domain.admin after the succesful execution of the function
	newAdminOutput, err := cr.adminUseCase.CreateAdmin(c.Request.Context(), newAdminInfo, adminID)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to create the admin", err.Error(), nil))
		return
	}

	//if no error, then  201 status as new admin is created succesfully
	c.JSON(http.StatusCreated, response.SuccessResponse(201, "admin created successfully", newAdminOutput))

}

// AdminLogin
// @Summary Admin Login
// @ID admin-login
// @Description Admin Login
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin_credentials body request.AdminLoginInfo true "Admin Login Credentials"
// @Success 200 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/login [post]
func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	//receive the data from request body
	var body request.AdminLoginInfo
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{StatusCode: 422, Message: "unable to process the request", Errors: err.Error(), Data: nil})
		return
	}
	//call the adminlogin method of the adminusecase to login as an admin
	tokenString, adminDataInModel, err := cr.adminUseCase.AdminLogin(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "failed to login", err.Error(), nil))
		return
	}
	c.SetSameSite(http.SameSiteLaxMode) //sets the SameSite attribute of the cookie to "Lax" mode. It is a security measure that helps protect against certain types of cross-site request forgery (CSRF) attacks.
	c.SetCookie("AdminAuth", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, response.SuccessResponse(200, "Succesfully Logged in", adminDataInModel))
}

// AdminLogout
// @Summary Admin_Logout
// @ID admin-logout
// @Description logout an logged-in admin from the site
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/logout [get]
func (cr *AdminHandler) AdminLogout(c *gin.Context) {
	// Set the user authentication cookie's expiration to -1 to invalidate it.
	c.Writer.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", "", -1, "", "", false, true)
	//c.Status(http.StatusOK)
	c.JSON(http.StatusOK, response.SuccessResponse(200, "Succesfully Logged-Out"))

}

// ListAllUsers
// @Summary Admin can list out all the registered users
// @ID list-all-users
// @Description The admin can list out all the registered users.
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Enter the page no to display"
// @Param limit query int false "Number of items to retrieve per page"
//
//	query query string false "Search query string"
//	filter query string false "filter criteria for showing the users"
//
// @Param sort_by query string false "sorting criteria for showing the users"
// @Param sort_desc query bool false "sorting in descending order"
// @Success 200 {object} response.Response
// @Success 204 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users [get]
func (cr *AdminHandler) ListAllUsers(c *gin.Context) {
	var viewUserInfo common.QueryParams

	viewUserInfo.Page, _ = strconv.Atoi(c.Query("page"))
	viewUserInfo.Limit, _ = strconv.Atoi(c.Query("limit"))
	viewUserInfo.Query = c.Query("query")
	viewUserInfo.Filter = c.Query("filter")
	viewUserInfo.SortBy = c.Query("sort_by")
	viewUserInfo.SortDesc, _ = strconv.ParseBool(c.Query("sort_desc"))

	users, isAnyUsers, err := cr.adminUseCase.ListAllUsers(c.Request.Context(), viewUserInfo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to fetch users", err.Error(), nil))
		return
	}
	//if isAnyUsers == false, then return status no content bacause user table is empty

	if !isAnyUsers {
		c.JSON(http.StatusNoContent, response.SuccessResponse(204, "No user found", users))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(200, "Succesfully fetched all users", users))

}

// Find User By ID
// @Summary Admin Fetch a specific user details by user id
// @ID find-user-by-id
// @Description Admin Fetch a specific user details by user id
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_id path string true "provide the ID of the user to be fetched"
// @Success 200 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/{user_id} [get]
func (cr *AdminHandler) FindUserByID(c *gin.Context) {
	paramsID := c.Param("id")
	userID, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to parse the user id", err.Error(), nil))
		return
	}
	user, err := cr.adminUseCase.FindUserByID(c.Request.Context(), userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to fetch the user", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "Succesfully fetched the user details", user))

}

// BlockUser
// @Summary Admin can block a registered user
// @ID block-user
// @Description Admin can block a user
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_id body request.BlockUser true "ID of the user to be blocked"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/:id/block [put]
func (cr *AdminHandler) BlockUser(c *gin.Context) {
	var blockUser request.BlockUser
	if err := c.Bind(&blockUser); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to read the request body", err.Error(), nil))
		return
	}

	adminID, err := handlerutil.GetAdminIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(400, "unable to fetch admin id from context", err.Error(), nil))
		return
	}
	blockedUser, err := cr.adminUseCase.BlockUser(c.Request.Context(), blockUser, adminID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to block the user", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "Succesfully blocked the user", blockedUser))
}

// UnblockUser
// @Summary Admin can unlock a blocked user
// @ID unblock-user
// @Description Admin can unblock a blocked user
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_id path string true "ID of the user to be unblocked"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/unblock/{user_id} [put]
func (cr *AdminHandler) UnblockUser(c *gin.Context) {
	paramsID := c.Param("id")
	fmt.Println("paramsid str is", paramsID)
	id, err := strconv.Atoi(paramsID)
	fmt.Println("paramsid is", id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(422, "failed to parse user id", err.Error(), nil))
		return
	}
	unBlockedUser, err := cr.adminUseCase.UnblockUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to unblock user", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "successfully unblocked the user", unBlockedUser))
}

// ---------------ADMIN DASHBOARD--------
// AdminDashboard
// @Summary Admin Dashboard
// @ID admin-dashboard
// @Description Admins dashboard will give summary regarding orders, users, products, etc.
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/dashboard [get]
func (cr *AdminHandler) AdminDashboard(c *gin.Context) {
	adminDashboard, err := cr.adminUseCase.AdminDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to fetch admindashboard", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(200, "successfully fetched the admin dashboard", adminDashboard))

}

// ------------FULL SALES REPORT-----------
// FullSalesReport
// @Summary Admin can download the full sales report
// @ID sales-report
// @Description Admin can download sales report in .csv format
// @Tags Admin
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date - Format - 2023-Mar-01"
// @Param end_date query string false "End Date - Format - 2023-Jun-01"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/sales-report/ [get]
func (cr *AdminHandler) FullSalesReport(c *gin.Context) {
	// time range to fetch details
	startDate, err1 := utils.StringToTime(c.Query("start_date"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "errorRR1", err1.Error(), nil))
		return
	}
	endDate, err2 := utils.StringToTime(c.Query("end_date"))
	if err2 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "errorRR2", err2.Error(), nil))
		return
	}
	//pages

	count, err3 := strconv.Atoi(c.Query("count"))
	if err3 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "errorRR3", err3.Error(), nil))
		return
	}
	pageNo, err4 := strconv.Atoi(c.Query("page_number"))
	if err4 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "error4", err4.Error(), nil))
		return
	}
	//join all errors and check error
	// err := errors.Join(err1, err2, err3, err4)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid inputs", err.Error(), nil))
	// 	return
	// }

	reqReportRange := common.SalesReportDateRange{
		StartDate: startDate,
		EndDate:   endDate,
		Pagination: common.Pagination{
			Count:      count,
			PageNumber: pageNo,
		},
	}

	salesReport, err := cr.adminUseCase.FetchFullSalesReport(c.Request.Context(), reqReportRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to create sales report", err.Error(), nil))
		return
	}
	if salesReport == nil {
		c.JSON(http.StatusOK, response.SuccessResponse(200, "there is no sales report this period", nil))
		return
	}

	//set headers for downloading in browser
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=salesreport.csv")
	csvWriter := csv.NewWriter(c.Writer)
	headers := []string{
		"UserID", "FirstName", "Email",
		"OrderID", "OrderTotalPrice", "AppliedCouponCode",
		"AppliedCouponDiscount", "OrderStatus", "DeliveryStatus",
		"PaymentType", "OrderDate",
	}

	if err := csvWriter.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to generate sales report", err.Error(), nil))
		return
	}

	if err := csvWriter.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to generate sales report", err.Error(), nil))
		return
	}

	for _, sales := range salesReport {
		row := []string{
			fmt.Sprintf("%v", sales.UserID),
			sales.FirstName,
			sales.Email,
			fmt.Sprintf("%v", sales.OrderID),
			fmt.Sprintf("%v", sales.OrderTotalPrice),
			sales.AppliedCouponCode,
			fmt.Sprintf("%v", sales.AppliedCouponDiscount),
			sales.OrderStatus,
			sales.DeliveryStatus,
			sales.PaymentType,
			sales.OrderDate.Format("2006-01-02 15:04:05"),
		}

		if err := csvWriter.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "failed to create sales report", err.Error(), nil))
			return
		}
	}
	// Flush the writer's buffer to ensure all data is written to the client
	csvWriter.Flush()

}
