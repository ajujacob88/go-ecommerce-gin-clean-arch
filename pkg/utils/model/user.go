package model

type BlockUser struct {
	UserID int    `json:"user_id"`
	Reason string `json:"reason"`
}
