package parsing

import (
	"io"
	"strings"
	"time"
)

func IntPtr(v int) *int {
	return &v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func MustParseTime(layout string, value string) time.Time {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return parsedTime
}

// CreateReader simulates a CSV reader with given header and row data.
// Each element of csvRow is one cell.
func CreateReader(header string, separator string, csvRow []string) io.Reader {
	row := strings.Join(csvRow, separator)
	return strings.NewReader(header + "\n" + row)
}
