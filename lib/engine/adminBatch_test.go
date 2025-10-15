package engine

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckBatchPaths(t *testing.T) {
	testCases := []struct {
		Filepath      BatchFile
		ErrorExpected bool
	}{
		// Really exists but is empty
		{NewBatchFile("test_data/empty_file"), false},

		// Does not exist
		{NewBatchFile("test_data/missing_file"), true},
	}

	for _, tc := range testCases {
		mockbatch := MockBatch("debit", []BatchFile{tc.Filepath})
		err := CheckBatchPaths(&mockbatch)
		if (err == nil && tc.ErrorExpected) ||
			(err != nil && !tc.ErrorExpected) {
			// t.Log(err.Error()) // delete_me
			t.Error("Validity of path " + tc.Filepath.Path() + " is wrongly checked")
		}
	}
}

func Test_ImportBatchDryRun(t *testing.T) {
	// Set up import
	adminBatch := MockBatch(Apdemande, []BatchFile{NewBatchFile("..", "parsing/apdemande/testData/apdemandeTestData.csv")})

	dataSinkFactory := &DiscardSinkFactory{}
	reportSink := &StdoutReportSink{}

	// Capture logs
	logs := new(strings.Builder)

	logger := slog.New(slog.NewTextHandler(logs, nil))

	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)

	slog.SetDefault(logger)

	// Run import
	err := ImportBatch(
		adminBatch,
		[]ParserType{},
		EmptyRegistry{},
		NoFilter,
		dataSinkFactory,
		reportSink,
	)
	assert.NoError(t, err)

	// Check that the import summary is part of the logs
	// assert.Contains(
	// 	t,
	// 	strings.ReplaceAll(logs.String(), "\\", ""),
	// 	`"summary": "../parsing/apdemande/testData/apdemandeTestData.csv: intégration terminée, 3 lignes traitées, 0 erreurs fatales, 0 lignes rejetées, 0 lignes filtrées, 3 lignes valides"`,
	// )
}
