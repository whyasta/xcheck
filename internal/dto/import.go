package dto

type ImportDto struct {
	ID             int64               `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	FileName       string              `gorm:"column:file_name" mapstructure:"file_name" json:"file_name" validate:"required,min=3,max=100"`
	UploadFileName string              `gorm:"column:upload_file_name" mapstructure:"upload_file_name" json:"upload_file_name" validate:"required,min=3,max=100"`
	ImportedAt     string              `gorm:"column:imported_at" mapstructure:"imported_at" json:"imported_at" validate:"required,min=3,max=20"`
	Status         string              `gorm:"column:status" mapstructure:"status" json:"status" validate:"required,min=3,max=20"`
	StatusMessage  string              `gorm:"column:status_message" mapstructure:"status_message" json:"status_message" validate:"required,min=3,max=20"`
	EventID        *int64              `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id"`
	BarcodeList    []*ImportBarcodeDto `gorm:"foreignKey:import_id;references:id" json:"barcode_list"`
}

type ImportBarcodeDto struct {
	ID           int64  `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	ImportID     int64  `gorm:"column:import_id" mapstructure:"import_id" json:"import_id" validate:"required,min=3,max=20"`
	Barcode      string `gorm:"column:barcode" mapstructure:"barcode" json:"barcode" validate:"required,min=3,max=20"`
	AssignStatus string `gorm:"column:assign_status,default:0" mapstructure:"assign_status" json:"assign_status" validate:"required,min=3,max=20"`
}

type UploadResponse struct {
	SucccessCount   int64  `gorm:"column:success_count" json:"success_count"`
	FailedCount     int64  `gorm:"column:failed_count" json:"failed_count"`
	DuplicateCount  int64  `gorm:"column:duplicate_count" json:"duplicate_count"`
	FailedValues    string `gorm:"column:failed_values" json:"failed_values"`
	DuplicateValues string `gorm:"column:duplicate_values" json:"duplicate_values"`
}
