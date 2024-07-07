package models

// UserRole represents a role that a user can have
type UserRole struct {
	ID          int64  `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	BaseModel
}
