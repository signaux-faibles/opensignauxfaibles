package main

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_convertAndConcat(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "output_*.csv")
	require.NoError(t, err)
	convertAndConcat(
		[]string{"resources/SF_DATA_20230706.txt", "resources/S_202011095834-3_202310020319.csv", "resources/S_202011095834-3_202311010315.csv"},
		output,
	)
}

func Test_readArgs(t *testing.T) {
	_, here, _, _ := runtime.Caller(0)
	tests := []struct {
		name       string
		args       []string
		wantInputs []string
		wantOutput string
	}{
		{
			name:       "tout va bien",
			args:       []string{here, "stock", "increment1", "increment2", "output"},
			wantInputs: []string{"stock", "increment1", "increment2"},
			wantOutput: "output",
		},
	}
	for _, tt := range tests {
		os.Args = tt.args
		t.Run(tt.name, func(t *testing.T) {
			gotInputs, gotOutput := readArgs()
			if !reflect.DeepEqual(gotInputs, tt.wantInputs) {
				t.Errorf("readArgs() gotInputs = %v, want %v", gotInputs, tt.wantInputs)
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("readArgs() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
