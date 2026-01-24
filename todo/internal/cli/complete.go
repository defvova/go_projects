package cli

import (
	"strings"
	"time"
	"todo/internal/storage"
)

func Complete(args []string) {
	storage := storage.NewCsvRepo()
	id := strings.Join(args, " ")
	now := time.Now().String()
	var fn = func(record []string) ([]string, bool) {
		if record[0] == id {
			record[3] = now
		}
		return record, true
	}
	storage.EditFile(fn)
}
