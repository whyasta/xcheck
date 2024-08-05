package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Filter struct {
	Property  string   `json:"prop,omitempty"`
	Operation string   `json:"opr"`
	Collation *string  `json:"coll,omitempty"`
	Value     string   `json:"val,omitempty"`
	Items     []Filter `json:"items,omitempty"`
}

func NewFilter(property string, operation string, collation *string, value string, items []Filter) *Filter {
	return &Filter{
		Property:  property,
		Operation: operation,
		Collation: collation,
		Value:     value,
		Items:     items,
	}
}

func NewFilters(items []Filter) *[]Filter {
	return &items
}

func (f *Filter) FilterResult(operator string, db *gorm.DB) *gorm.DB {
	// log.Println("budal")
	if f.Items != nil {
		for i := 0; i < len(f.Items); i++ {
			val := f.Items[i]
			if i == 0 {
				db = val.FilterResult("", db)
			} else {
				db = val.FilterResult(f.Operation, db)
			}
		}
		return db
	}

	if strings.ToLower(f.Operation) == "like" {
		f.Value = "%" + f.Value + "%"
		if f.Value == "%%" { // hack for empty value
			return db
		}
	}

	query := fmt.Sprintf("%s %s ?", f.Property, f.Operation)
	// log.Println(query, f.Value)
	if strings.ToLower(operator) == "or" {
		return db.Or(query, f.Value)
	}
	return db.Where(query, f.Value)
}

// func FilterResult(db *gorm.DB, filters []Filter) *gorm.DB {
// 	for _, filter := range filters {
// 		fmt.Printf("%s %s ?", filter.Property, filter.Operation)
// 		//tx = tx.Where(fmt.Sprintf("%s %s ?", filter.Property, filter.Operation), filter.Value)
// 	}
// 	return db
// }

// 	return db
// }
