package models

// swagger:model
type User struct {
	ID       int64    `gorm:"column:id; primary_key; not null" json:"id"`
	Username string   `gorm:"column:username" json:"username" validate:"required,min=5,max=20"`
	Password string   `gorm:"column:password" json:"password,omitempty" validate:"required,min=2,max=32"`
	Email    string   `gorm:"column:email" json:"email" validate:"required,email"`
	RoleID   int64    `gorm:"column:role_id" json:"role_id" validate:"required"`
	Role     UserRole `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	CommonModel
}

// swagger:model UserLogin
type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// swagger:model
type SignedResponse struct {
	Token string `json:"token"`
}

// swagger:model
type UserRequest struct {
	Username string `gorm:"column:username" mapstructure:"username" json:"username" validate:"required,min=5,max=20"`
	Password string `gorm:"column:password" mapstructure:"password" json:"password" validate:"required,min=2,max=32"`
	Email    string `gorm:"column:email" mapstructure:"email" json:"email" validate:"required,email"`
	RoleID   int64  `gorm:"column:role_id" mapstructure:"role_id" json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Username string `gorm:"column:username" mapstructure:"username" json:"username,omitempty" validate:"omitempty"`
	Password string `gorm:"column:password" mapstructure:"password" json:"password,omitempty" validate:"omitempty"`
	Email    string `gorm:"column:email" mapstructure:"email" json:"email,omitempty" validate:"email,omitempty"`
	RoleID   int64  `gorm:"column:role_id" mapstructure:"role_id" json:"role_id,omitempty"`
}

// swagger:parameters getUser deleteUser
type UserID struct {
	// In: path
	ID int `json:"id"`
}

// swagger:model RefreshToken
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// swagger:parameters authSignin
type UserLoginBodyParams struct {
	// required: true
	// in: body
	UserLogin *UserLogin `json:"UserLogin"`
}

// swagger:parameters authRefreshToken
type UserRefreshTokenBodyParams struct {
	// required: true
	// in: body
	UserLogin *RefreshToken `json:"RefreshToken"`
}

// swagger:parameters createUser
type UserCreateBodyParams struct {
	// required: true
	// in: body
	UserRequest *UserRequest `json:"UserRequest"`
}
