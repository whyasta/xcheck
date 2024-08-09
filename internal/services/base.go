package services

import (
	"bigmind/xcheck-be/internal/repositories"

	"gorm.io/gorm"
)

type Service struct {
	AuthService           *AuthService
	UserService           *UserService
	RoleService           *RoleService
	EventService          *EventService
	TicketTypeService     *TicketTypeService
	GateService           *GateService
	SessionService        *SessionService
	GateAllocationService *GateAllocationService
	ImportService         *ImportService
	BarcodeService        *BarcodeService
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
		AuthService:           NewAuthService(repositories.User),
		UserService:           NewUserService(repositories.User, repositories.Base),
		RoleService:           NewRoleService(repositories.Role),
		EventService:          NewEventService(repositories.Event),
		TicketTypeService:     NewTicketTypeService(repositories.TicketType, repositories.Base),
		GateService:           NewGateService(repositories.Gate),
		SessionService:        NewSessionService(repositories.Session),
		ImportService:         NewImportService(repositories.Import),
		BarcodeService:        NewBarcodeService(repositories.Barcode, repositories.GateAllocation),
		GateAllocationService: NewGateAllocationService(repositories.GateAllocation),
	}
}
