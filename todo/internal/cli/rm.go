package cli

import (
	"strings"
	"todo/internal/storage"
)

func Rm(args []string) {
	storage := storage.NewCsvRepo()
	id := strings.Join(args, " ")
	var fn = func(record []string) ([]string, bool) {
		return record, record[0] != id
	}
	storage.EditFile(fn)
}
