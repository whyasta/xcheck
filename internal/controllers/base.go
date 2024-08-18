package controllers

import (
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"reflect"
	"strconv"

	gormUtils "gorm.io/gorm/utils"
)

type Controller struct {
	HealthController         *HealthController
	AuthController           *AuthController
	UserController           *UserController
	RoleController           *RoleController
	EventController          *EventController
	GateController           *GateController
	TicketTypeController     *TicketTypeController
	GateAllocationController *GateAllocationController
	SessionController        *SessionController
	BarcodeController        *BarcodeController
	SyncController           *SyncController
}

func NewController(
	services *services.Service,
) *Controller {
	return &Controller{
		HealthController:         new(HealthController),
		AuthController:           NewAuthController(services.AuthService),
		UserController:           NewUserController(services.UserService),
		RoleController:           NewRoleController(services.RoleService),
		EventController:          NewEventController(services.EventService),
		GateController:           NewGateController(services.GateService),
		TicketTypeController:     NewTicketTypeController(services.TicketTypeService),
		GateAllocationController: NewGateAllocationController(services.GateAllocationService),
		SessionController:        NewSessionController(services.SessionService),
		BarcodeController:        NewBarcodeController(services.ImportService, services.BarcodeService),
		SyncController:           NewSyncController(services.SyncService),
	}
}

func MakeQueryParams(params map[string][]string, allowedParams []string) map[string]interface{} {
	newParams := make(map[string]interface{})
	for key, value := range params {
		if !gormUtils.Contains(allowedParams, key) {
			continue
		}
		newParams[key] = value[0]
	}
	return newParams
}

func MakePaginationQueryParams(params map[string][]string, allowedParams []string) (*utils.Paginate, map[string]interface{}) {
	newParams := make(map[string]interface{})
	pageParams := make(map[string]interface{})

	for key, value := range params {
		if !gormUtils.Contains(allowedParams, key) && !gormUtils.Contains([]string{"page", "limit"}, key) {
			continue
		}
		if gormUtils.Contains([]string{"page", "limit"}, key) {
			pageParams[key] = value[0]
			continue
		}
		newParams[key] = value[0]
	}

	keys := reflect.ValueOf(params).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}

	if !gormUtils.Contains(strkeys, "page") {
		pageParams["page"] = 1
	}
	if !gormUtils.Contains(strkeys, "limit") {
		pageParams["limit"] = 10
	}

	limit, _ := strconv.Atoi(pageParams["limit"].(string))
	page, _ := strconv.Atoi(pageParams["page"].(string))

	paginate := utils.NewPaginate(limit, page)
	return paginate, newParams
}

func MakePageFilterQueryParams(params map[string][]string, allowedParams []string) (*utils.Paginate, []utils.Filter, []utils.Sort) {
	pageParams := make(map[string]interface{})

	for key, value := range params {
		if !gormUtils.Contains(allowedParams, key) && !gormUtils.Contains([]string{"page", "limit"}, key) {
			continue
		}
		if gormUtils.Contains([]string{"page", "limit"}, key) {
			pageParams[key] = value[0]
			continue
		}
	}

	keys := reflect.ValueOf(params).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}

	if !gormUtils.Contains(strkeys, "page") {
		pageParams["page"] = 1
	}
	if !gormUtils.Contains(strkeys, "limit") {
		pageParams["limit"] = 10
	}

	limit, _ := strconv.Atoi(pageParams["limit"].(string))
	page, _ := strconv.Atoi(pageParams["page"].(string))

	paginate := utils.NewPaginate(limit, page)

	filterParams := MakeFilterQueryParams(params)
	sortParams := MakeSortQueryParams(params)

	// log.Println("ini kah?", sortParams)
	return paginate, filterParams, sortParams
}

func MakeFilterQueryParams(params map[string][]string) []utils.Filter {
	var filters []utils.Filter
	//newParams := make([]map[string]interface{}, 0)

	for key, value := range params {
		if key != "filter" {
			continue
		}

		if err := json.Unmarshal([]byte(value[0]), &filters); err != nil {
			panic(err)
		}
		// fmt.Println(filters)
	}

	// for i := 0; i < len(filters); i++ {
	//     if ()
	//     filters[i].Value = utils.QuoteValue(filters[i].Value)
	// }

	return filters
}

func MakeSortQueryParams(params map[string][]string) []utils.Sort {
	var sorts []utils.Sort

	for key, value := range params {
		print(key)
		if key != "sort" {
			continue
		}

		if err := json.Unmarshal([]byte(value[0]), &sorts); err != nil {
			panic(err)
		}
		// fmt.Println(filters)
	}

	return sorts
}
