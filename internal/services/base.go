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
	SyncService           *SyncService
	ReportService         *ReportService
	RedeemService         *RedeemService
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
		SyncService:           NewSyncService(repositories.Base, repositories.Event, repositories.TicketType, repositories.Gate, repositories.Session),
		TicketTypeService:     NewTicketTypeService(repositories.TicketType, repositories.Base),
		GateService:           NewGateService(repositories.Gate),
		SessionService:        NewSessionService(repositories.Session),
		ImportService:         NewImportService(repositories.Import),
		BarcodeService:        NewBarcodeService(repositories.Barcode, repositories.Gate, repositories.Session),
		GateAllocationService: NewGateAllocationService(repositories.GateAllocation),
		ReportService:         NewReportService(repositories.Base, repositories.Report),
		RedeemService:         NewRedeemService(repositories.Redeem),
	}
}
