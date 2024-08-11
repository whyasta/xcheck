package models

type Gate struct {
	ID       int64  `gorm:"column:id; primary_key; not null" json:"id" mapstructure:"id"`
	GateName string `gorm:"column:gate_name"  mapstructure:"gate_name" json:"gate_name" validate:"required,min=3,max=20"`
	EventID  int64  `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id,omitempty" validate:"required"`
	Users    []User `gorm:"many2many:gate_users;" json:"users,omitempty"`
	CommonModel
}

type GateUser struct {
	GateID int64 `gorm:"column:gate_id;" json:"gate_id" mapstructure:"gate_id"`
	UserID int64 `gorm:"column:user_id;" json:"user_id" mapstructure:"user_id"`
}
