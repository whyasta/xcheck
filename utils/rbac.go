package utils

import "gorm.io/gorm/utils"

type Menu struct {
	Label      string   `json:"label"`
	Slug       string   `json:"slug"`
	Permission []string `json:"-"`
	// Child []Menu `json:"child"`
}

type MenuID int
type RoleID int64

const (
	ADM RoleID = iota + 1
	CHECKER
	REDEMPTION
)

const (
	DASHBOARD MenuID = iota + 1
	USER_MANAGEMENT
	EVENT_SYNC
	USER_MANAGEMENT_USER
	USER_MANAGEMENT_ROLE
	EVENT_DETAIL
	EVENT_DETAIL_SESSION
	EVENT_DETAIL_GATE
	EVENT_DETAIL_TICKET_TYPE
	EVENT_DETAIL_BARCODE
	EVENT_DETAIL_CHECKIN
	EVENT_DETAIL_CHECKOUT
	EVENT_DETAIL_REPORT_TRAFFIC
	EVENT_DETAIL_REPORT_UNIQUE
	EVENT_DETAIL_REPORT_REALTIME
	EVENT_DETAIL_REDEMPTION_TICKET_FILE
	EVENT_DETAIL_REDEMPTION_TICKET_DATA
	EVENT_DETAIL_REDEMPTION_TICKET_REDEMPTION
	EVENT_DETAIL_REDEMPTION_REPORT
)

func GetAllMenu() []Menu {
	return []Menu{
		Menu{
			Label:      "Dashboard",
			Slug:       "/dashboard",
			Permission: []string{"ADMIN", "CHECKER", "REDEMPTION"},
		},
		Menu{
			Label:      "User Management",
			Slug:       "/user-management",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Event Sync",
			Slug:       "/event-sync",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "User",
			Slug:       "/user-management/user",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Role",
			Slug:       "/user-management/role",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Event Detail",
			Slug:       "/event/${slug}",
			Permission: []string{"ADMIN", "CHECKER", "REDEMPTION"},
		},
		Menu{
			Label:      "Session",
			Slug:       "/event/${slug}/session",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Gate",
			Slug:       "/event/${slug}/gate",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Ticket Type",
			Slug:       "/event/${slug}/ticket-type",
			Permission: []string{"ADMIN"},
		},
		Menu{
			Label:      "Associate Barcode",
			Slug:       "/event/${slug}/barcode",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Check In",
			Slug:       "/event/${slug}/checker/check-in",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Check Out",
			Slug:       "/event/${slug}/checker/check-out",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Traffic Visitor",
			Slug:       "/event/${slug}/report/traffic-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Unique Visitor",
			Slug:       "/event/${slug}/report/unique-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Real Time Visitor",
			Slug:       "/event/${slug}/report/real-time-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Ticket File",
			Slug:       "/event/${slug}/redemption/ticket-file",
			Permission: []string{"ADMIN", "REDEMPTION"},
		},
		Menu{
			Label:      "Ticket Data",
			Slug:       "/event/${slug}/redemption/ticket-data",
			Permission: []string{"ADMIN", "REDEMPTION"},
		},
		Menu{
			Label:      "Ticket Redemption",
			Slug:       "/event/${slug}/redemption/ticket-redemption",
			Permission: []string{"ADMIN", "REDEMPTION"},
		},
		Menu{
			Label:      "Report",
			Slug:       "/event/${slug}/redemption/report",
			Permission: []string{"ADMIN", "REDEMPTION"},
		},
	}
}
func (r MenuID) GetMenu() Menu {
	return [...]Menu{
		Menu{
			Label:      "Dashboard",
			Slug:       "/dashboard",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "User Management",
			Slug:       "/user-management",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Event Sync",
			Slug:       "/event-sync",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "User",
			Slug:       "/user-management/user",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Role",
			Slug:       "/user-management/role",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Event Detail",
			Slug:       "/event/${slug}",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Session",
			Slug:       "/event/${slug}/session",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Gate",
			Slug:       "/event/${slug}/gate",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Ticket Type",
			Slug:       "/event/${slug}/ticket-type",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Associate Barcode",
			Slug:       "/event/${slug}/barcode",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Check In",
			Slug:       "/event/${slug}/checker/check-in",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Check Out",
			Slug:       "/event/${slug}/checker/check-out",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Traffic Visitor",
			Slug:       "/event/${slug}/report/traffic-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Unique Visitor",
			Slug:       "/event/${slug}/report/unique-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Real Time Visitor",
			Slug:       "/event/${slug}/report/real-time-visitor",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Ticket File",
			Slug:       "/event/${slug}/redemption/ticket-file",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Ticket Data",
			Slug:       "/event/${slug}/redemption/ticket-data",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Ticket Redemption",
			Slug:       "/event/${slug}/redemption/ticket-redemption",
			Permission: []string{"ADMIN", "CHECKER"},
		},
		Menu{
			Label:      "Report",
			Slug:       "/event/${slug}/redemption/report",
			Permission: []string{"ADMIN", "CHECKER"},
		},
	}[r-1]
}

func (r RoleID) GetRole() string {
	return [...]string{
		"ADMIN",
		"CHECKER",
		"REDEMPTION",
	}[r-1]
}

func (r RoleID) GetMenu() []Menu {
	var result []Menu

	for _, v := range GetAllMenu() {
		if utils.Contains(v.Permission, r.GetRole()) {
			result = append(result, v)
		}
	}

	return result
}
