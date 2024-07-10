package models

// UserRole represents a role that a user can have
type UserRole struct {
	ID          int64  `gorm:"column:id; primary_key; not null" json:"id"`
	RoleName    string `gorm:"column:role_name" json:"role_name"`
	Description string `gorm:"column:description" json:"description"`
	CommonModel
}

type UserRoleRequest struct {
	RoleName    string `gorm:"column:role_name" json:"role_name"`
	Description string `gorm:"column:description" json:"description"`
}

func (g UserRole) ToEntity() UserRole {
	return UserRole{
		ID:          g.ID,
		RoleName:    g.RoleName,
		Description: g.Description,
	}
}

func (g UserRole) FromEntity(role UserRole) interface{} {
	return UserRole{
		ID:          g.ID,
		RoleName:    g.RoleName,
		Description: g.Description,
	}
}
