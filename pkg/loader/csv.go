package loader

import (
	"encoding/csv"
	"log/slog"
	"strings"

	"github.com/WangYihang/gojob/pkg/utils"
)

func CatCSV(path string, index int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for line := range utils.Cat(path) {
			r := csv.NewReader(strings.NewReader(line))
			record, err := r.Read()
			if err != nil {
				continue
			}
			if index >= 0 && index < len(record) {
				out <- record[index]
			} else {
				slog.Error("invalid csv record", "error", err, "index", index, "line", line)
			}
		}
	}()
	return out
}
