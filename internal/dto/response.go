package dto

type ResponseDto struct {
	ID       int64    `gorm:"column:id; primary_key; not null" json:"id" mapstructure:"id"`
	GateName string   `gorm:"column:gate_name"  mapstructure:"gate_name" json:"gate_name" validate:"required,min=3,max=100"`
	EventID  int64    `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id" validate:"required"`
	Users    []*int64 `gorm:"column:users"  mapstructure:"users" json:"users"`
}
