package services

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/gocraft/work"
	"github.com/minio/minio-go/v7"
)

type ImportService struct {
	r repositories.ImportRepository
}

func NewImportService(r repositories.ImportRepository) *ImportService {
	return &ImportService{r}
}

func (s *ImportService) CreateImport(data *models.Import) (models.Import, error) {
	return s.r.Save(data)
}

func (s *ImportService) UploadToMinio(ctx context.Context, bucketName string, file *multipart.FileHeader, fileName string) (string, error) {
	// Get Buffer from file
	buffer, err := file.Open()

	if err != nil {
		return "", err
	}
	defer buffer.Close()

	// Create minio connection.
	minioClient, err := config.MinioConnection()
	if err != nil {
		return "", err
	}

	objectName := fileName
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	// Upload the zip file with PutObject
	log.Printf("Uploading %s of size %d to %s\n", objectName, fileSize, bucketName)
	info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return "", err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return objectName, nil
}

func (s *ImportService) RemoveFromMinio(ctx context.Context, bucketName string, objectName string) error {
	// Create minio connection.
	minioClient, err := config.MinioConnection()
	if err != nil {
		return err
	}

	err = minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	log.Printf("Successfully removed %s\n", objectName)
	return nil
}

func (s *ImportService) UpdateStatusImport(id int64, status constant.ImportStatus, errorMessage string) (models.Import, error) {
	return s.r.Update(id, &map[string]interface{}{"status": status, "status_message": errorMessage})
}

func (s *ImportService) DoImportBarcodeJob(id int64) (models.Import, error) {
	fmt.Println("DoImportJob")
	row, err := s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusProcessing, "status_message": "Processing file"})
	if err != nil {
		return models.Import{}, err
	}

	config.GetEnqueuer().Enqueue("import_barcode", work.Q{
		"import_id":      id,
		"headers":        "barcode",
		"csv_file":       row.FileName,
		"table":          "raw_barcodes",
		"with_assign":    false,
		"ticket_type_id": 0,
		"sessions":       "",
		"gates":          "",
	})

	// TODO: Implement import job
	// importJob := utils.NewImport(s.r.GetDB(), id, "raw_barcodes", row.FileName, []string{"barcode"})
	// fmt.Println("Importing data...")
	// importJob.ImportData()
	// fmt.Println("Importing done")
	// // return models.Import{}, nil
	// row, err = s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusCompleted, "status_message": "Completed"})
	// if err == nil {
	// 	if err = os.Remove(row.FileName); err != nil {
	// 		return models.Import{}, err
	// 	}
	// }

	return row, err
}

func (s *ImportService) DoImportJobWithAssign(id int64, eventID int64, ticketTypeID int64, sessions string, gates string) (*work.Job, models.Import, error) {
	fmt.Println("DoImportJob")
	row, err := s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusProcessing, "status_message": "Processing file"})
	if err != nil {
		return &work.Job{}, models.Import{}, err
	}

	job, err := config.GetEnqueuer().Enqueue("import_barcode", work.Q{
		"import_id":      id,
		"headers":        "barcode",
		"csv_file":       row.FileName,
		"table":          "raw_barcodes",
		"with_assign":    true,
		"event_id":       eventID,
		"ticket_type_id": ticketTypeID,
		"sessions":       sessions,
		"gates":          gates,
	})
	if err != nil {
		return &work.Job{}, models.Import{}, err
	}
	return job, row, err
}

func (s *ImportService) DoImportTicketJob(id int64, eventID int64, withHeader bool) (*work.Job, models.Import, error) {
	fmt.Println("DoImportTicketJob")
	row, err := s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusProcessing, "status_message": "Processing file"})
	if err != nil {
		return &work.Job{}, models.Import{}, err
	}

	job, err := config.GetEnqueuer().Enqueue("import_ticket", work.Q{
		"csv_file":    row.FileName,
		"table":       "tickets",
		"import_id":   id,
		"with_header": withHeader,
		"headers":     "order_barcode,order_id,ticket_type_name,name,email,phone_number,identity_id,birth_date,company_name,associate_barcode",
		"event_id":    eventID,
	})
	if err != nil {
		return &work.Job{}, models.Import{}, err
	}
	return job, row, err
}

func (s *ImportService) GetAllImports(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error) {
	return s.r.FindAll(paginate, filters, sorts)
}

func (s *ImportService) GetImportByID(uid int64) (models.Import, error) {
	return s.r.FindByID(uid)
}

func (s *ImportService) CheckValid(importID int64, assignID int64) (bool, error) {
	_, err := s.r.CheckValidImport(importID)
	if err != nil {
		return false, errors.New("invalid import id")
	}

	valid, err := s.r.CheckValidAssign(assignID)
	if err != nil {
		return false, errors.New("invalid event assign id")
	}

	return valid, err
}

func (s *ImportService) DeleteImport(id int64) (models.Import, error) {
	return s.r.Delete(id)
}

func (s *ImportService) GetUploadResult(id int64) (dto.UploadResponse, error) {
	return s.r.GetUploadResult(id)
}
