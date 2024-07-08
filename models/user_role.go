package models

// UserRole represents a role that a user can have
type UserRole struct {
	ID          int64  `gorm:"column:id; primary_key; not null" json:"id"`
	RoleName    string `gorm:"column:role_name" json:"role_name"`
	Description string `gorm:"column:description" json:"description"`
	BaseModel
}

type UserRoleRequest struct {
	RoleName    string `gorm:"column:role_name" json:"role_name"`
	Description string `gorm:"column:description" json:"description"`
	BaseModel
}
