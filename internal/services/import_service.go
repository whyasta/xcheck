package services

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"

	"github.com/gocraft/work"
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

func (s *ImportService) UpdateStatusImport(id int64, status string, errorMessage string) (models.Import, error) {
	return s.r.Update(id, &map[string]interface{}{"status": status, "error_message": errorMessage})
}

func (s *ImportService) DoImportJob(id int64) (models.Import, error) {
	fmt.Println("DoImportJob")
	row, err := s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusProcessing, "error_message": "Processing file"})
	if err != nil {
		return models.Import{}, err
	}

	config.GetEnqueuer().Enqueue("import_barcode", work.Q{
		"import_id": id,
		"headers":   "barcode",
		"csv_file":  row.FileName,
		"table":     "import_barcodes",
	})

	// TODO: Implement import job
	// importJob := utils.NewImport(s.r.GetDB(), id, "import_barcodes", row.FileName, []string{"barcode"})
	// fmt.Println("Importing data...")
	// importJob.ImportData()
	// fmt.Println("Importing done")
	// // return models.Import{}, nil
	// row, err = s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusCompleted, "error_message": "Completed"})
	// if err == nil {
	// 	if err = os.Remove(row.FileName); err != nil {
	// 		return models.Import{}, err
	// 	}
	// }

	return row, err
}

func (s *ImportService) GetAllImports(paginate *utils.Paginate, filters []utils.Filter) ([]models.Import, int64, error) {
	return s.r.FindAll(paginate, filters)
}

func (s *ImportService) GetImportByID(uid int64) (models.Import, error) {
	return s.r.FindByID(uid)
}

func (s *ImportService) CheckValid(importId int64, assignId int64) (bool, error) {
	_, err := s.r.CheckValidImport(importId)
	if err != nil {
		return false, errors.New("invalid import id")
	}

	valid, err := s.r.CheckValidAssign(assignId)
	if err != nil {
		return false, errors.New("invalid event assign id")
	}

	return valid, err
}
