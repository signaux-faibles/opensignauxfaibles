package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func Test_forgeBrowserToken(t *testing.T) {
	viper.SetConfigFile("testData/config.toml")
	viper.ReadInConfig()

	browser := Browser{
		Name:    "test",
		IP:      "10.10.10.10",
		Created: time.Date(2017, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		Email:   "testeur@domaine.test",
	}

	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkIjoiMjAxNy0wMS0wMVQwMDowMDowMFoiLCJlbWFpbCI6InRlc3RldXJAZG9tYWluZS50ZXN0IiwiaXAiOiIxMC4xMC4xMC4xMCIsIm5hbWUiOiJ0ZXN0In0.3ZcIWjYmGwbfelDMoFzfBPLvcwa0sKDN14iuLX1bfeg"
	browserToken, err := forgeBrowserToken(browser)
	fmt.Println(browserToken)
	fmt.Println(browser)
	if err != nil {
		t.Error("Erreur forgeBrowserToken: " + err.Error())
	} else if browserToken.BrowserToken != testToken {
		t.Error("Le token généré est différent du token attendu")
	} else {
		t.Log("Test forgeBrowserToken ok")
	}

}

func Test_readBrowserToken(t *testing.T) {
	viper.SetConfigFile("testData/config.toml")
	viper.ReadInConfig()

	browser := Browser{
		Name:    "test",
		IP:      "10.10.10.10",
		Created: time.Date(2017, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		Email:   "testeur@domaine.test",
	}

	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkIjoiMjAxNy0wMS0wMVQwMDowMDowMFoiLCJlbWFpbCI6InRlc3RldXJAZG9tYWluZS50ZXN0IiwiaXAiOiIxMC4xMC4xMC4xMCIsIm5hbWUiOiJ0ZXN0In0.3ZcIWjYmGwbfelDMoFzfBPLvcwa0sKDN14iuLX1bfeg"
	testBrowser, err := readBrowserToken(testToken)

	if err != nil {
		t.Error("Erreur forgeBrowserToken: " + err.Error())
	} else if testBrowser != browser {
		t.Error("L'objet généré est différent du token attendu")
	} else {
		t.Log("Test forgeBrowserToken ok")
	}

}
