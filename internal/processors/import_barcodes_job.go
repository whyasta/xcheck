package processors

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocraft/work"
	"gorm.io/gorm"
)

type Import struct {
	DB       *gorm.DB
	ImportID int64
	Table    string
	CsvFile  string
	Headers  []string
}

func NewImport(
	db *gorm.DB,
	importId int64,
	table string,
	csvFile string,
	headers []string,
) *Import {
	fmt.Println("new import", []interface{}{csvFile, table, importId, headers})
	return &Import{
		CsvFile:  csvFile,
		Table:    table,
		DB:       db,
		ImportID: importId,
		Headers:  headers,
	}
}

func ImportBarcodeJob(job *work.Job) error {
	csvFile := job.ArgString("csv_file")
	table := job.ArgString("table")
	importId := job.ArgInt64("import_id")
	headers := job.ArgString("headers")
	if err := job.ArgError(); err != nil {
		fmt.Println("=> import barcode job error", err.Error())
		return err
	}

	fmt.Println("=> import barcode job", []interface{}{csvFile, table, importId, headers})
	db, _ := config.ConnectToDB()

	importJob := NewImport(db, importId, table, csvFile, strings.Split(headers, ","))

	fmt.Println("Importing data...")
	importJob.ImportData()
	fmt.Println("Importing done")

	importRepo := repositories.NewImportRepository(db)

	row, err := importRepo.Update(importId, &map[string]interface{}{"status": constant.ImportStatusCompleted, "error_message": "Completed"})
	if err == nil {
		if err = os.Remove(row.FileName); err != nil {
			return err
		}
	}
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	return nil
}

func generateQuestionsMark(n int) []string {
	s := make([]string, 0)
	for i := 0; i < n; i++ {
		s = append(s, "?")
	}
	return s
}

func (i Import) ImportData() error {
	fmt.Println("=> import data")

	start := time.Now()

	csvReader, csvFile, err := i.openCsvFile()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvFile.Close()

	jobs := make(chan []interface{}, 0)
	wg := new(sync.WaitGroup)

	go i.dispatchWorkers(jobs, wg)
	i.readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
	return nil
}

func (i Import) dispatchWorkers(jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= 10; workerIndex++ {
		go func(workerIndex int, db *gorm.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				i.doInsertJob(workerIndex, counter, db, job)
				wg.Done()
				counter++
			}
		}(workerIndex, i.DB, jobs, wg)
	}
}

func (i Import) openCsvFile() (*csv.Reader, *os.File, error) {
	fmt.Println("=> open csv file")

	f, err := os.Open(i.CsvFile)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(f)
	return reader, f, nil
}

func (i Import) readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if len(i.Headers) == 0 {
			i.Headers = row
			continue
		}

		rowOrdered := make([]interface{}, 0)
		for _, each := range row {
			rowOrdered = append(rowOrdered, each)
		}

		wg.Add(1)
		jobs <- rowOrdered
	}
	close(jobs)
}

func (i Import) doInsertJob(workerIndex, counter int, db *gorm.DB, values []interface{}) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				i.Table,
				strings.Join(i.Headers, ",")+",import_id",
				strings.Join(generateQuestionsMark(len(i.Headers)), ",")+",?",
			)

			importId := []interface{}{i.ImportID}
			values = append(values, importId...)

			err := db.WithContext(context.Background()).Exec(query, values...).Error
			if err != nil {
				log.Fatal(err.Error())
			}

			if err != nil {
				log.Fatal(err.Error())
			}
		}(&outerError)
		if outerError == nil {
			break
		}
	}

	// if counter%100 == 0 {
	// 	fmt.Println("=> worker", workerIndex, "inserted", counter, "data")
	// }
}