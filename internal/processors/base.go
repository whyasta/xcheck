package processors

import (
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

	"gorm.io/gorm"
)

type Import struct {
	DB         *gorm.DB
	ImportID   int64
	Table      string
	CsvFile    string
	Headers    []string
	WithHeader bool
}

func NewImport(
	db *gorm.DB,
	importID int64,
	table string,
	csvFile string,
	headers []string,
	withHeader bool,
) *Import {
	fmt.Println("new import", []interface{}{csvFile, table, importID, headers})
	return &Import{
		CsvFile:    csvFile,
		Table:      table,
		DB:         db,
		ImportID:   importID,
		Headers:    headers,
		WithHeader: withHeader,
	}
}

func (i Import) ImportData() error {
	fmt.Println("=> import data")

	start := time.Now()

	csvReader, csvFile, err := i.openCsvFile()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer csvFile.Close()

	jobs := make(chan []interface{})
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
	if i.WithHeader {
		// Read and skip the first row (header)
		if _, err := reader.Read(); err != nil {
			log.Fatal("Error reading the first row:", err)
			return nil, nil, err
		}
	}
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
					*outerError = fmt.Errorf("outer error %v", err)
					importRepo := repositories.NewImportRepository(db)
					importRepo.Update(i.ImportID, &map[string]interface{}{"status": constant.ImportStatusFailed, "status_message": err})
				}
			}()

			query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				i.Table,
				strings.Join(i.Headers, ",")+",import_id",
				strings.Join(generateQuestionsMark(len(i.Headers)), ",")+",?",
			)

			fmt.Printf("=> do insert %s\n", query)

			if len(values) != len(i.Headers) {
				values = values[0:len(i.Headers)]
			}

			importID := []interface{}{i.ImportID}
			values = append(values, importID...)

			err := db.WithContext(context.Background()).Exec(query, values...).Error
			if err != nil {
				importRepo := repositories.NewImportRepository(db)
				importRepo.Update(i.ImportID, &map[string]interface{}{"status": constant.ImportStatusFailed, "status_message": err.Error()})
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

func generateQuestionsMark(n int) []string {
	s := make([]string, 0)
	for i := 0; i < n; i++ {
		s = append(s, "?")
	}
	return s
}