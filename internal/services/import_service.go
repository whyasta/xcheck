package services

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"log"
	"os"
	"path/filepath"
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
	log.Println("DoImportJob")
	row, err := s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusProcessing, "error_message": "Processing file"})
	if err != nil {
		return models.Import{}, err
	}

	// TODO: Implement import job
	importJob := utils.NewImport(s.r.GetDB(), id, "import_barcodes", filepath.Join("files", row.FileName), []string{"barcode"})
	log.Println("Importing data...")
	importJob.ImportData()
	log.Println("Importing done")
	// return models.Import{}, nil
	row, err = s.r.Update(id, &map[string]interface{}{"status": constant.ImportStatusCompleted, "error_message": "Completed"})
	if err == nil {
		if err = os.Remove(filepath.Join("files", row.FileName)); err != nil {
			return models.Import{}, err
		}
	}
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
