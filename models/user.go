package models

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username" validate:"required,min=5,max=20"`
	Password     string `json:"password,omitempty" validate:"required,min=2,max=32"`
	PasswordHash string `json:"-"`
	Email        string `json:"email" validate:"required,email"`
	Role         string `json:"role" validate:"required"`
	BaseModel
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignedResponse struct {
	Token string `json:"token"`
}