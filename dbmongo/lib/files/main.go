package files

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"
)

// FileSummary représente un fichier
type FileSummary struct {
	Name string    `json:"name" bson:"name"`
	Size int64     `json:"size" bson:"size"`
	Date time.Time `json:"date" bson:"date"`
}

// ListFiles liste les fichiers présents dans APP_DATA
func ListFiles(basePath string) ([]FileSummary, error) {
	var files []FileSummary
	basePathConf := viper.GetString("APP_DATA")
	b := len(basePathConf)

	currentFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []FileSummary{}, err
	}

	for _, file := range currentFiles {
		if file.IsDir() {
			subPath := fmt.Sprintf("%s/%s", basePath, file.Name())
			subFiles, err := ListFiles(subPath)
			if err != nil {
				return []FileSummary{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, FileSummary{
				Name: fmt.Sprintf("%s/%s", basePath, file.Name())[b:],
				Size: file.Size(),
				Date: file.ModTime(),
			})
		}
	}
	return files, nil
}
