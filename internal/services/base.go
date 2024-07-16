package services

import (
	"bigmind/xcheck-be/internal/repositories"

	"gorm.io/gorm"
)

type Service struct {
	AuthService       *AuthService
	UserService       *UserService
	RoleService       *RoleService
	EventService      *EventService
	TicketTypeService *TicketTypeService
}

func RegisterServices(db *gorm.DB) *Service {
	return NewService(
		repositories.NewRepository(db),
	)
}

func NewService(
	repositories *repositories.Repository,
) *Service {
	return &Service{
		AuthService:       NewAuthService(repositories.User),
		UserService:       NewUserService(repositories.User),
		RoleService:       NewRoleService(repositories.Role),
		EventService:      NewEventService(repositories.Event),
		TicketTypeService: NewTicketTypeService(repositories.TicketType),
	}
}
