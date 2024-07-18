package models

type Import struct {
	ID           int64            `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	FileName     string           `gorm:"column:file_name" mapstructure:"file_name" json:"file_name" validate:"required,min=3,max=20"`
	ImportedAt   string           `gorm:"column:imported_at" mapstructure:"imported_at" json:"imported_at" validate:"required,min=3,max=20"`
	Status       string           `gorm:"column:status" mapstructure:"status" json:"status" validate:"required,min=3,max=20"`
	ErrorMessage string           `gorm:"column:error_message" mapstructure:"error_message" json:"error_message" validate:"required,min=3,max=20"`
	BarcodeList  []*ImportBarcode `gorm:"foreignKey:import_id;references:id" json:"barcode_list"`
	CommonModel
}

type ImportBarcode struct {
	ID           int64  `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	ImportID     int64  `gorm:"column:import_id" mapstructure:"import_id" json:"import_id" validate:"required,min=3,max=20"`
	Barcode      string `gorm:"column:barcode" mapstructure:"barcode" json:"barcode" validate:"required,min=3,max=20"`
	AssignStatus string `gorm:"column:assign_status,default:0" mapstructure:"assign_status" json:"assign_status" validate:"required,min=3,max=20"`
	CommonModel
}
