package dto

type UserRoleRequest struct {
	RoleName    string `gorm:"column:role_name" json:"role_name"`
	Description string `gorm:"column:description" json:"description"`
}
