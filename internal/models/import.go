package models

type Import struct {
	ID             int64  `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	FileName       string `gorm:"column:file_name" mapstructure:"file_name" json:"file_name" validate:"required"`
	UploadFileName string `gorm:"column:upload_file_name" mapstructure:"upload_file_name" json:"upload_file_name" validate:"required,min=3,max=100"`
	ImportedAt     string `gorm:"column:imported_at" mapstructure:"imported_at" json:"imported_at" validate:"required"`
	Status         string `gorm:"column:status" mapstructure:"status" json:"status" validate:"required"`
	StatusMessage  string `gorm:"column:status_message" mapstructure:"status_message" json:"status_message" validate:"required,min=3,max=100"`
	EventID        *int64 `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id"`
	// TotalBarcode   int64            `gorm:"column:total_barcode"  mapstructure:"total_barcode" json:"total_barcode"`
	BarcodeList []*ImportBarcode `gorm:"foreignKey:import_id;references:id" json:"barcode_list"`
	CommonModel
}

type ImportBarcode struct {
	ID           int64  `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	ImportID     int64  `gorm:"column:import_id" mapstructure:"import_id" json:"import_id" validate:"required"`
	Barcode      string `gorm:"column:barcode" mapstructure:"barcode" json:"barcode" validate:"required,min=3,max=100"`
	AssignStatus string `gorm:"column:assign_status,default:0" mapstructure:"assign_status" json:"assign_status"`
	CommonModel
}
