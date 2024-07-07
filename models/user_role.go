package models

// UserRole represents a role that a user can have
type UserRole struct {
	ID          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}
