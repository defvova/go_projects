package cli

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
	"todo/internal/storage"
)

func Add(args []string) error {
	storage := storage.NewCsvRepo()
	id, _ := storage.LastID()
	id++
	now := time.Now().String()

	file, err := os.OpenFile(storage.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	storage.SkipOrAddHeaders(file, w)

	record := []string{strconv.Itoa(id), strings.Join(args, " "), now, ""}
	if err := w.Write(record); err != nil {
		return err
	}

	return w.Error()
}
