package utils

import (
	"fmt"

	"gorm.io/gorm"
)

type Sort struct {
	Property  string `json:"prop,omitempty"`
	Direction string `json:"dir"`
}

func NewSort(property string, direction string) *Sort {
	return &Sort{
		Property:  property,
		Direction: direction,
	}
}

func (s *Sort) SortResult(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf("%s %s", s.Property, s.Direction))
}
