package utils

import (
	"log"

	"gorm.io/gorm"
)

type Paginate struct {
	limit int
	page  int
}

func NewPaginate(limit int, page int) *Paginate {
	return &Paginate{limit: limit, page: page}
}

func (p *Paginate) PaginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.limit
	log.Println("offset: ", offset)

	return db.Offset(offset).
		Limit(p.limit)
}

func (p *Paginate) Limit(count int64) int {
	return p.limit
}

func (p *Paginate) Page(count int64) int {
	return p.page
}
