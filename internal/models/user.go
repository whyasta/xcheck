package models

// type AuthUUID struct {
// 	AuthUUID json.RawMessage `gorm:"column:auth_uuid" json:"auth_uuid"`
// }

type User struct {
	ID        int64     `gorm:"column:id; primary_key; not null" json:"id"`
	Username  string    `gorm:"column:username" json:"username" validate:"required,min=5,max=20"`
	Password  string    `gorm:"column:password" json:"password,omitempty" validate:"required,min=2,max=32"`
	Email     string    `gorm:"column:email" json:"email" validate:"required,email"`
	RoleID    *int64    `gorm:"column:role_id" json:"role_id,omitempty"`
	Role      *UserRole `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	AuthUuids *string   `gorm:"column:auth_uuids" json:"-"`
	CommonModel
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignedResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRequest struct {
	Username string `gorm:"column:username" mapstructure:"username" json:"username" validate:"required,min=5,max=20"`
	Password string `gorm:"column:password" mapstructure:"password" json:"password" validate:"required,min=2,max=32"`
	Email    string `gorm:"column:email" mapstructure:"email" json:"email" validate:"required,email"`
	RoleID   int64  `gorm:"column:role_id" mapstructure:"role_id" json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Username string  `gorm:"column:username" mapstructure:"username" json:"username,omitempty" validate:"omitempty"`
	Password string  `gorm:"column:password" mapstructure:"password" json:"password,omitempty" validate:"omitempty"`
	Email    *string `gorm:"column:email" mapstructure:"email" json:"email,omitempty" validate:"omitempty,email"`
	RoleID   int64   `gorm:"column:role_id" mapstructure:"role_id" json:"role_id,omitempty"`
}

type UserID struct {
	// In: path
	ID int `json:"id"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserLoginBodyParams struct {
	// required: true
	// in: body
	UserLogin *UserLogin `json:"UserLogin"`
}

type UserRefreshTokenBodyParams struct {
	// required: true
	// in: body
	UserLogin *RefreshToken `json:"RefreshToken"`
}

type UserCreateBodyParams struct {
	// required: true
	// in: body
	UserRequest *UserRequest `json:"UserRequest"`
}
