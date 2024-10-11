package processors

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/repositories"
	"fmt"
	"os"
	"strings"

	"github.com/gocraft/work"
)

func ImportTicketJob(job *work.Job) error {
	csvFile := job.ArgString("csv_file")
	table := job.ArgString("table")
	importID := job.ArgInt64("import_id")
	headers := job.ArgString("headers")
	eventID := job.ArgInt64("event_id")
	withHeader := job.ArgBool("with_header")

	if err := job.ArgError(); err != nil {
		db, _ := config.ConnectToDB()
		importRepo := repositories.NewImportRepository(db)
		_, err := importRepo.Update(importID, &map[string]interface{}{"status": constant.ImportStatusFailed, "status_message": "Failed"})
		defer func() {
			dbInstance, _ := db.DB()
			_ = dbInstance.Close()
		}()
		fmt.Println("=> Import ticket job error", err.Error())
		return err
	}

	fmt.Println("=> Import ticket job ID", job.ID)
	fmt.Println("=> Import ticket job args", []interface{}{csvFile, table, importID, headers, eventID})
	db, _ := config.ConnectToDB()

	importJob := NewImport(db, importID, table, csvFile, strings.Split(headers, ","), withHeader)

	fmt.Println("=> Importing data...")
	importJob.ImportData()

	importRepo := repositories.NewImportRepository(db)

	row, err := importRepo.Update(importID, &map[string]interface{}{"status": constant.ImportStatusCompleted, "status_message": "Completed"})
	if err == nil {
		if err = os.Remove(row.FileName); err != nil {
			return err
		}
		// dispatch other job (if any)
	}
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	fmt.Println("Importing done")

	// validating all records
	ticketRepo := repositories.NewTicketRepository(db)
	ticketRepo.ValidateImport(importID, eventID)

	return nil
}
