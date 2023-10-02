package marshal

import (
	"strconv"
	"time"
)

func TimeToCSV(t *time.Time) string {
	if t != nil {
		return t.Format(time.DateOnly)
	}
	return ""
}

func FloatToCSV(f *float64) string {
	if f != nil {
		return strconv.FormatFloat(*f, 'f', -1, 64)
	}
	return ""
}

func IntToCSV(f *int) string {
	if f != nil {
		return strconv.Itoa(*f)
	}
	return ""
}

func BoolToCSV(b *bool) string {
	if b != nil {
		return strconv.FormatBool(*b)
	}
	return ""
}
