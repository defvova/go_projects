package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mergestat/timediff"
)

const ParseTimeLayout string = "2006-01-02 15:04:05.999999 -0700 MST"
const CSVFileName string = "data.csv"
const TempCSVFileName string = "temp_" + CSVFileName

type CsvRepo struct {
	Path     string
	TempPath string
}

type FnOperation func([]string) ([]string, bool)

func NewCsvRepo() *CsvRepo {
	return &CsvRepo{
		Path:     CSVFileName,
		TempPath: TempCSVFileName,
	}
}

func (csvRepo *CsvRepo) LastID() (int, error) {
	file, err := os.Open(csvRepo.Path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	id := 0
	headerRow := true

	for {
		record, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		if headerRow && csvRepo.IsHeader(record) {
			headerRow = false
			continue
		}
		headerRow = false

		i, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		id = i
	}

	return id, nil
}

func (csvRepo *CsvRepo) EditFile(fn FnOperation) error {
	file, err := os.Open(csvRepo.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tempFile, err := os.OpenFile(csvRepo.TempPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()

	r := csv.NewReader(file)
	w := csv.NewWriter(tempFile)

	for {
		record, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		newRecord, result := fn(record)
		if result {
			if err := w.Write(newRecord); err != nil {
				return fmt.Errorf("failed to write record to temp file: %w", err)
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	file.Close()
	tempFile.Close()

	if err := os.Remove(csvRepo.Path); err != nil {
		return fmt.Errorf("failed to remove original file: %w", err)
	}

	if err := os.Rename(csvRepo.TempPath, csvRepo.Path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (csvRepo *CsvRepo) IsHeader(row []string) bool {
	expected := csvRepo.Headers()
	if len(row) != len(expected) {
		return false
	}
	for i := range row {
		if strings.ToLower(row[i]) != expected[i] {
			return false
		}
	}
	return true
}

func (csvRepo *CsvRepo) Headers() []string {
	return []string{"id", "task", "created_at", "completed_at"}
}

func (csvRepo *CsvRepo) SkipOrAddHeaders(file *os.File, w *csv.Writer) {
	info, _ := file.Stat()
	if info.Size() == 0 {
		w.Write(csvRepo.Headers())
	}
}

func (csvRepo *CsvRepo) RenderRow(record []string) string {
	id := record[0]
	text := record[1]
	createdAtDiff := parseDate(record[2])
	completedAtDiff := parseDate(record[3])
	return fmt.Sprintf("%v\t%v\t%v\t%v", id, text, createdAtDiff, completedAtDiff)
}

func parseDate(val string) string {
	var diff string

	if val != "" {
		ts := strings.Split(val, " m=")[0]
		dateAt, _ := time.Parse(ParseTimeLayout, ts)
		diff = timediff.TimeDiff(dateAt)
	}

	return diff
}
