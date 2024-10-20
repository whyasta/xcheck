package dto

import (
	"bigmind/xcheck-be/internal/models"

	"gorm.io/datatypes"
)

type UserLoginResponse struct {
	ID        int64            `gorm:"column:id; primary_key; not null" json:"id"`
	Username  string           `gorm:"column:username" json:"username" validate:"required,min=5,max=20"`
	Password  string           `gorm:"column:password" json:"password,omitempty" validate:"required,min=2,max=32"`
	Email     string           `gorm:"column:email" json:"email" validate:"required,email"`
	RoleID    *int64           `gorm:"column:role_id" json:"role_id,omitempty"`
	Role      *models.UserRole `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	AuthUuids *string          `gorm:"column:auth_uuids" json:"-"`
	MenuGroup datatypes.JSON   `json:"menu_group"`
}
