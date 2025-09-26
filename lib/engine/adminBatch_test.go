package engine

import (
	"log/slog"
	"strings"
	"testing"

	"opensignauxfaibles/lib/base"

	"github.com/stretchr/testify/assert"
)

func Test_IsBatchID(t *testing.T) {
	if !base.IsBatchID("1801") {
		t.Error("1801 devrait être un ID de batch")
	}

	if base.IsBatchID("") {
		t.Error("'' ne devrait pas être considéré comme un ID de batch")
	}

	if base.IsBatchID("190193039") {
		t.Error("'190193039' ne devrait pas être considéré comme un ID de batch")
	}
	if !base.IsBatchID("1901_93039") {
		t.Error("'1901_93039'  devrait être considéré comme un ID de batch")
	}

	if base.IsBatchID("abcd") {
		t.Error("'abcd' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'abcd' est bien rejeté: ")
	}
}

func Test_CheckBatchPaths(t *testing.T) {
	testCases := []struct {
		Filepath      base.BatchFile
		ErrorExpected bool
	}{
		{base.NewBatchFile("test_data/empty_file"), false},
		{base.NewBatchFile("test_data/missing_file"), true},
	}
	for _, tc := range testCases {
		mockbatch := base.MockBatch("debit", []base.BatchFile{tc.Filepath})
		err := CheckBatchPaths(&mockbatch)
		if (err == nil && tc.ErrorExpected) ||
			(err != nil && !tc.ErrorExpected) {
			// t.Log(err.Error()) // delete_me
			t.Error("Validity of path " + tc.Filepath.Path() + " is wrongly checked")
		}
	}
}

type TestSinkFactory struct{}

func (TestSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	return &DiscardDataSink{}, nil
}

func Test_ImportBatch(t *testing.T) {

	batchProvider := base.BasicBatchProvider{Batch: base.AdminBatch{}}
	err := ImportBatch(batchProvider, []base.ParserType{}, false, TestSinkFactory{}, DiscardReportSink{})
	if err == nil {
		t.Error("ImportBatch devrait nous empêcher d'importer sans filtre")
	}
}

func Test_ImportBatchWithUnreadableFilter(t *testing.T) {
	batchProvider := base.BasicBatchProvider{
		Batch: base.MockBatch("filter", []base.BatchFile{base.NewBatchFile("this_file_does_not_exist")}),
	}
	err := ImportBatch(batchProvider, []base.ParserType{}, false, TestSinkFactory{}, DiscardReportSink{})
	if err == nil {
		t.Error("ImportBatch devrait échouer en tentant d'ouvrir un fichier filtre illisible")
	}
}

func Test_ImportBatchWithSinkFailure(t *testing.T) {
	batch := base.AdminBatch{
		Files: base.BatchFiles{
			"apdemande": {base.NewBatchFile("../..", "lib/apdemande/testData/apdemandeTestData.csv")},
		},
	}
	batchProvider := base.BasicBatchProvider{Batch: batch}
	err := ImportBatch(batchProvider, []base.ParserType{base.Apdemande}, true, FailSinkFactory{}, DiscardReportSink{})
	if err == nil {
		t.Error("ImportBatch devrait échouer si le sink échoue")
	}
}

func Test_ImportBatchDryRun(t *testing.T) {
	// Set up import
	adminBatch := base.MockBatch(base.Apdemande, []base.BatchFile{base.NewBatchFile("..", "apdemande/testData/apdemandeTestData.csv")})
	batchProvider := base.BasicBatchProvider{Batch: adminBatch}

	noFilter := true

	dataSinkFactory := &DiscardSinkFactory{}
	reportSink := &StdoutReportSink{}

	// Capture logs
	logs := new(strings.Builder)

	logger := slog.New(slog.NewTextHandler(logs, nil))

	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)

	slog.SetDefault(logger)

	// Run import
	err := ImportBatch(batchProvider, []base.ParserType{}, noFilter, dataSinkFactory, reportSink)
	assert.NoError(t, err)

	// Check that the import summary is part of the logs
	assert.Contains(
		t,
		strings.ReplaceAll(logs.String(), "\\", ""),
		`"summary": "../apdemande/testData/apdemandeTestData.csv: intégration terminée, 3 lignes traitées, 0 erreurs fatales, 0 lignes rejetées, 0 lignes filtrées, 3 lignes valides"`,
	)
}
