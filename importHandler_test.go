package main

import (
	"opensignauxfaibles/lib/base"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	adminBatch := base.MockBatch(base.Apdemande, []base.BatchFile{base.NewBatchFile("./lib/apdemande/testData/apdemandeTestData.csv")})
	handler := importBatchHandler{
		Enable:      true,
		Path:        "",
		BatchKey:    adminBatch.Key.String(),
		Parsers:     []string{},
		NoFilter:    true,
		BatchConfig: "",
		DryRun:      true,
		adminBatch:  &adminBatch,
	}

	err := handler.Run()
	assert.NoError(t, err)
}
