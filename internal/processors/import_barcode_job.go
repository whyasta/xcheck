package processors

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/repositories"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocraft/work"
)

func ImportBarcodeJob(job *work.Job) error {
	csvFile := job.ArgString("csv_file")
	table := job.ArgString("table")
	importID := job.ArgInt64("import_id")
	headers := job.ArgString("headers")
	withAssign := job.ArgBool("with_assign")
	eventID := job.ArgInt64("event_id")
	ticketTypeID := job.ArgInt64("ticket_type_id")
	sessions := job.ArgString("sessions")
	gates := job.ArgString("gates")

	if err := job.ArgError(); err != nil {
		db, _ := config.ConnectToDB()
		importRepo := repositories.NewImportRepository(db)
		_, err := importRepo.Update(importID, &map[string]interface{}{"status": constant.ImportStatusFailed, "status_message": "Failed"})
		defer func() {
			dbInstance, _ := db.DB()
			_ = dbInstance.Close()
		}()
		fmt.Println("=> Import barcode job error", err.Error())
		return err
	}

	fmt.Println("=> Import barcode job ID", job.ID)
	fmt.Println("=> Import barcode job args", []interface{}{csvFile, table, importID, headers, withAssign, eventID, ticketTypeID, sessions, gates})
	db, _ := config.ConnectToDB()

	importJob := NewImport(db, importID, table, csvFile, strings.Split(headers, ","), false)

	fmt.Println("=> Importing data...")
	importJob.ImportData()

	importRepo := repositories.NewImportRepository(db)

	row, err := importRepo.Update(importID, &map[string]interface{}{"status": constant.ImportStatusCompleted, "status_message": "Completed"})
	if err == nil {
		if err = os.Remove(row.FileName); err != nil {
			return err
		}
		// dispatch other job (if any)

		if withAssign {
			sessionSlice := make([]int64, len(strings.Split(sessions, ",")))
			for i, str := range strings.Split(sessions, ",") {
				num, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
				}
				sessionSlice[i] = num
			}

			gateSlice := make([]int64, len(strings.Split(gates, ",")))
			for i, str := range strings.Split(gates, ",") {
				num, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
				}
				gateSlice[i] = num
			}

			barcodeRepo := repositories.NewBarcodeRepository(db)
			_, _, _, _, _, err := barcodeRepo.AssignBarcodesWithEvent(importID, eventID, ticketTypeID, sessionSlice, gateSlice)
			if err != nil {
				fmt.Println("Error assigning barcodes:", err)
				return err
			}
		}
	}
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	fmt.Println("Importing done")

	return nil
}
