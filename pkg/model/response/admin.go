package response

//model structs used to output data

type AdminDataOutput struct {
	ID           uint
	UserName     string
	Email        string
	Phone        string
	IsSuperAdmin bool
}

type AdminDashboard struct {
	CompletedOrders int     `json:"completed_orders"`
	PendingOrders   int     `json:"pending_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
	TotalOrders     int     `json:"total_orders"`
	TotalOrderItems int     `json:"total_order_items"`
	OrderValue      float64 `json:"order_value"`
	CreditedAmount  float64 `json:"credited_amount"`
	PendingAmount   float64 `json:"pending_amount"`
	TotalUsers      int     `json:"total_users"`
	VerifiedUsers   int     `json:"verified_users"`
	OrderedUsers    int     `json:"ordered_users"`
}
