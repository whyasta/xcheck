package checks

import (
	"database/sql"
	"reflect"
)

type SQLCheck struct {
	SQL *sql.DB
}

func (s SQLCheck) Pass() bool {
	if s.SQL == nil {
		return false
	}

	err := s.SQL.Ping()

	return err == nil
}

func (s SQLCheck) Name() string {
	if s.SQL == nil {
		return "no_driver"
	}

	driver := s.SQL.Driver()
	return reflect.TypeOf(driver).String()
}
