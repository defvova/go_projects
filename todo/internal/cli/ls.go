package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"todo/internal/storage"
)

func Ls(args []string) error {
	storage := storage.NewCsvRepo()
	isAll := slices.Contains(args, "-a")
	file, err := os.Open(storage.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	for {
		record, err := r.Read()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if storage.IsHeader(record) {
			fmt.Fprintln(w, strings.Join(storage.Headers(), "\t"))
		} else {
			if isAll {
				fmt.Fprintln(w, storage.RenderRow(record))
			} else {
				if record[3] == "" {
					fmt.Fprintln(w, storage.RenderRow(record))
				}
			}
		}
	}
}
