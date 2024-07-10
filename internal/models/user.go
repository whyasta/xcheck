package models

type User struct {
	ID       int64    `gorm:"column:id; primary_key; not null" json:"id"`
	Username string   `gorm:"column:username" json:"username" validate:"required,min=5,max=20"`
	Password string   `gorm:"column:password" json:"password,omitempty" validate:"required,min=2,max=32"`
	Email    string   `gorm:"column:email" json:"email" validate:"required,email"`
	RoleID   int64    `gorm:"column:role_id" json:"role_id" validate:"required"`
	Role     UserRole `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	CommonModel
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignedResponse struct {
	Token string `json:"token"`
}

type UserRequest struct {
	Username string `gorm:"column:username" json:"username" validate:"required,min=5,max=20"`
	Password string `gorm:"column:password" json:"password" validate:"required,min=2,max=32"`
	Email    string `gorm:"column:email" json:"email" validate:"required,email"`
	RoleID   int64  `gorm:"column:role_id" json:"role_id" validate:"required"`
}
