package main

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		viper.Set("DB_DIAL", "mongodb://localhost:27016")
		viper.Set("DB_DIAL", "signauxfaibles")
		os.Args = []string{"sfdata", "etablissements"}
		main()
		wg.Done()
	}()
	wg.Wait()
}
