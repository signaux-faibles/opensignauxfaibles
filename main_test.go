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
		viper.Set("Verbose", true) // cf https://github.com/spf13/viper#setting-overrides
		os.Args = []string{"sfdata", "--help"}
		main()
		wg.Done()
	}()
	wg.Wait()
}
