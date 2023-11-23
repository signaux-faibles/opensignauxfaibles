package altares

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isIncrementalEndOfFile(t *testing.T) {
	recordEOF, err := csv.NewReader(strings.NewReader("Fin du fichier : total 122909 ligne(s);")).Read()
	t.Log("match -> ", END_OF_FILE_REGEXP.MatchString(recordEOF[0]))
	require.NoError(t, err)
	type args struct {
		record []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ok", args{recordEOF}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIncrementalEndOfFile(tt.args.record); got != tt.want {
				t.Errorf("isIncrementalEndOfFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
