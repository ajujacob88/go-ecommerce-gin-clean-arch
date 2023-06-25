package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/common"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	interfaces "github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/repository/interface"

	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}
}

func (c *adminDatabase) IsSuperAdmin(ctx context.Context, adminID int) (bool, error) {
	var isSuperAdmin bool
	superAdminCheckingQuery := `SELECT is_super_admin
								FROM admins
								WHERE id = $1` //In the SQL query string, the placeholder $1 indicates the position of the first parameter that will be passed when executing the query. In this case, the value of adminID is passed as the first parameter to the Raw method.
	err := c.DB.Raw(superAdminCheckingQuery, adminID).Scan(&isSuperAdmin).Error //This line executes the SQL query using the DB.Raw method provided by the c.DB database connection. It scans the result into the isSuperAdmin variable using the &isSuperAdmin syntax. Any error that occurs during the execution is assigned to the err variable.
	return isSuperAdmin, err
}

func (c *adminDatabase) CreateAdmin(ctx context.Context, newAdminInfo request.NewAdminInfo) (domain.Admin, error) {
	var newAdmin domain.Admin
	createAdminQuery := `	INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
						 	VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *;`
	err := c.DB.Raw(createAdminQuery, newAdminInfo.UserName, newAdminInfo.Email, newAdminInfo.Phone, newAdminInfo.Password).Scan(&newAdmin).Error
	newAdmin.Password = "" //By setting it to an empty string before returning, the function ensures that the password is not accessible outside of the function scope.
	return newAdmin, err
}

func (c *adminDatabase) FindAdmin(ctx context.Context, email string) (domain.Admin, error) {
	var adminData domain.Admin
	findAdminQuery := `	SELECT *
						FROM admins
						WHERE email = $1;`

	err := c.DB.Raw(findAdminQuery, email).Scan(&adminData).Error
	return adminData, err
}

func (c *adminDatabase) ListAllUsers(ctx context.Context, queryParams common.QueryParams) ([]domain.Users, bool, error) {
	findQuery := "SELECT * FROM users"
	params := []interface{}{}

	fmt.Println("queryparams is", queryParams)

	if queryParams.Query != "" && queryParams.Filter != "" {
		findQuery = fmt.Sprintf("%s WHERE LOWER(%s) LIKE $%d", findQuery, queryParams.Filter, len(params)+1)
		params = append(params, "%"+strings.ToLower(queryParams.Query)+"%")
		fmt.Println("params is ", params)
	}
	if queryParams.SortBy != "" {
		findQuery = fmt.Sprintf("%s ORDER BY %s %s", findQuery, queryParams.SortBy, orderByDirection(queryParams.SortDesc))
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		findQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", findQuery, len(params)+1, len(params)+2)
		params = append(params, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}

	var users []domain.Users
	err := c.DB.Raw(findQuery, params...).Scan(&users).Error
	if err != nil {
		return nil, false, err
	}

	return users, len(users) > 0, nil
}

func orderByDirection(sortDesc bool) string {
	if sortDesc {
		return "DESC"
	}
	return "ASC"
}

/* Another method for list all users, but above method is good and standard one

func (c *userDatabase) ListAllUsers(ctx context.Context, queryParams model.QueryParams) ([]domain.Users, error) {

	findQuery := "SELECT * FROM users"
	if queryParams.Query != "" && queryParams.Filter != "" {
		findQuery = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", findQuery, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			findQuery = fmt.Sprintf("%s ORDER BY %s DESC", findQuery, queryParams.SortBy)
		} else {
			findQuery = fmt.Sprintf("%s ORDER BY %s ASC", findQuery, queryParams.SortBy)
		}
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		findQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", findQuery, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		findQuery = fmt.Sprintf("%s LIMIT 10 OFFSET 0", findQuery)
	}
	var users []domain.Users

	err := c.DB.Raw(findQuery).Scan(&users).Error
	if len(users) == 0 {
		return users, fmt.Errorf("no users found")
	}
	return users, err
}

*/

func (c *adminDatabase) FindUserByID(ctx context.Context, userID int) (domain.Users, error) {
	var user domain.Users
	findUser := `SELECT * FROM users WHERE id = $1;`
	err := c.DB.Raw(findUser, userID).Scan(&user).Error
	if user.ID == 0 {
		return domain.Users{}, fmt.Errorf("no user is found with that id")
	}
	return user, err
}

func (c *adminDatabase) BlockUser(ctx context.Context, blockInfo request.BlockUser, adminID int) (domain.UserInfo, error) {
	var userInfo domain.UserInfo
	blockQuery := `UPDATE user_infos SET is_blocked = 'true', blocked_at = NOW(), blocked_by = $1, reason_for_blocking = $2 WHERE users_id = $3 RETURNING *;`
	err := c.DB.Raw(blockQuery, adminID, blockInfo.Reason, blockInfo.UserID).Scan(&userInfo).Error

	if userInfo.UsersID == 0 {
		return domain.UserInfo{}, fmt.Errorf("User not found")
	}
	return userInfo, err
}

func (c *adminDatabase) UnblockUser(ctx context.Context, userID int) (domain.UserInfo, error) {
	var userInfo domain.UserInfo
	unblockQuery := `UPDATE user_infos SET is_blocked = 'false', reason_for_blocking = '' WHERE users_id = $1 RETURNING *;`
	err := c.DB.Raw(unblockQuery, userID).Scan(&userInfo).Error
	if userInfo.UsersID == 0 {
		return domain.UserInfo{}, fmt.Errorf("no user found")
	}
	return userInfo, err
}

func (c *adminDatabase) FetchOrdersSummaryData(ctx context.Context) (response.AdminDashboard, error) {
	var adminDashboard response.AdminDashboard
	orderSummaryFetchQuery := ` 	SELECT 
									COUNT(CASE WHEN order_status_id = 11 THEN id END) AS completed_orders,
									COUNT(CASE WHEN order_status_id = 1 OR order_status_id = 2 OR order_status_id = 3 OR order_status_id = 4 OR order_status_id = 5 THEN id END) AS pending_orders, 
									COUNT(CASE WHEN order_status_id = 7 OR order_status_id = 8 THEN id END) AS cancelled_orders,
									COUNT(id) AS total_orders,
									SUM (CASE WHEN order_status_id != 7 AND order_status_id != 8 THEN order_total_price ELSE 0 END) AS order_value,
									COUNT(DISTINCT user_id) AS ordered_users
									FROM orders;`

	err := c.DB.Raw(orderSummaryFetchQuery).Scan(&adminDashboard).Error
	if err != nil {
		return response.AdminDashboard{}, err
	}

	return adminDashboard, nil
}

func (c *adminDatabase) FetchTotalOrderedItems(ctx context.Context) (int, error) {
	var totalOrderedItems int
	totalOrderedItemsQuery := `	SELECT
								COUNT(id) AS total_order_items
								FROM order_lines;`
	err := c.DB.Raw(totalOrderedItemsQuery).Scan(&totalOrderedItems).Error
	if err != nil {
		return 0, err
	}
	return totalOrderedItems, nil
}

func (c *adminDatabase) FetchTotalCreditedAmount(ctx context.Context) (float64, error) {
	var totalCreditedAmount float64
	creditedAmountQuery := `SELECT
							sum(order_total_price) AS credited_amount
							FROM payment_details WHERE payment_status_id = 2;`

	err := c.DB.Raw(creditedAmountQuery).Scan(&totalCreditedAmount).Error
	if err != nil {
		return 0, err
	}

	return totalCreditedAmount, nil
}

func (c *adminDatabase) FetchUsersCount(ctx context.Context, adminDashboardData response.AdminDashboard) (response.AdminDashboard, error) {
	userCountQuery := `	SELECT 
						COUNT(*) AS total_users, 
						COUNT(CASE WHEN is_verified = true THEN 1 END) AS verified_users
  						FROM user_infos;`

	err := c.DB.Raw(userCountQuery).Scan(&adminDashboardData).Error
	if err != nil {
		return response.AdminDashboard{}, err
	}

	return adminDashboardData, nil
}

func (c *adminDatabase) FetchFullSalesReport(ctx context.Context, reqReportRange common.SalesReportDateRange) ([]response.SalesReport, error) {
	var fullSalesReport []response.SalesReport

	limit := reqReportRange.Pagination.Count
	offset := (reqReportRange.Pagination.PageNumber - 1) * limit

	reportQuery := ` 	SELECT o.user_id,u.first_name, u.email, o.id AS order_id, o.order_total_price, 
						c.coupon_code AS applied_coupon_code,o.applied_coupon_discount, 
						os.status AS order_status, ds.status AS delivery_status, pm.payment_type, o.order_date  
						FROM orders o
						INNER JOIN order_statuses os ON o.order_status_id = os.id 
						INNER JOIN payment_method_infos pm ON o.payment_method_info_id = pm.id 
						INNER JOIN users u ON o.user_id = u.id 
 						LEFT JOIN coupons c ON o.applied_coupon_id = c.id
						LEFT JOIN delivery_statuses ds ON o.delivery_status_id = ds.id
						WHERE order_date >= $1 AND order_date <= $2
						ORDER BY o.order_date 
						LIMIT  $3 OFFSET $4`

	if c.DB.Raw(reportQuery, reqReportRange.StartDate, reqReportRange.EndDate, limit, offset).Scan(&fullSalesReport).Error != nil {
		return []response.SalesReport{}, errors.New("failed to collect the data to create the sales report")
	}

	return fullSalesReport, nil
}
