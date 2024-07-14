package utils

import (
	"log"

	"gorm.io/gorm"
)

// swagger:parameters getUserList getEventList
type Paginate struct {
	// required: true
	Limit int `json:"limit"`
	// required: true
	Page int `json:"page"`
}

func NewPaginate(limit int, page int) *Paginate {
	return &Paginate{Limit: limit, Page: page}
}

func (p *Paginate) PaginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	log.Println("offset: ", offset)

	return db.Offset(offset).
		Limit(p.Limit)
}

func (p *Paginate) GetLimit(count int64) int {
	return p.Limit
}

func (p *Paginate) GetPage(count int64) int {
	return p.Page
}
