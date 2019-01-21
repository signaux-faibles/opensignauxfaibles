package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//
// @summary Liste les fichiers disponibles dans le dépot
// @description Tous ces fichiers sont contenu dans APP_DATA (voir config.toml)
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/files [get]
func adminFiles(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := listFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

type fileSummary struct {
	Name string    `json:"name" bson:"name"`
	Size int64     `json:"size" bson:"size"`
	Date time.Time `json:"date" bson:"date"`
}

func listFiles(basePath string) ([]fileSummary, error) {
	var files []fileSummary
	basePathConf := viper.GetString("APP_DATA")
	b := len(basePathConf)

	currentFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []fileSummary{}, err
	}

	for _, file := range currentFiles {
		if file.IsDir() {
			subPath := fmt.Sprintf("%s/%s", basePath, file.Name())
			subFiles, err := listFiles(subPath)
			if err != nil {
				return []fileSummary{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, fileSummary{
				Name: fmt.Sprintf("%s/%s", basePath, file.Name())[b:],
				Size: file.Size(),
				Date: file.ModTime(),
			})
		}
	}
	return files, nil
}

var importFunctions = map[string]func(*AdminBatch) error{
	"apconso":    importAPConso,
	"bdf":        importBDF,
	"delai":      importDelai,
	"apdemande":  importAPDemande,
	"diane":      importDiane,
	"cotisation": importCotisation,
	"dpae":       importDPAE,
	"altares":    importAltares,
	"procol":     importProcol,
	"ccsf":       importCCSF,
	"debit":      importDebit,
	"effectif":   importEffectif,
	"sirene":     importSirene,
}

//
// @summary Purge la collection RawData
// @description Suppression de tous les objets de données brutes contenus dans la collection RawData (irréversible)
// @Tags Traitements
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/purge [get]
func purge(c *gin.Context) {
	db.DB.C("RawData").RemoveAll(nil)
	c.String(200, "Done")
}

//
// @summary Import de fichiers pour un batch
// @description Effectue l'import de tous les fichiers du batch donné en paramètre
// @Tags Traitements
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Param batch query string true "Clé du batch"
// @Success 200 {string} string ""
// @Router /api/data/import/{batch} [get]
func importBatchHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batch")
	batch := AdminBatch{}
	batch.load(batchKey)
	go importBatch(&batch)
}
func importBatch(batch *AdminBatch) {
	if !batch.Readonly {
		for fnName, fn := range importFunctions {
			journal(info, "importMain", "Début de l'import du type: "+fnName+" pour le batch "+batch.ID.Key)
			err := fn(batch)
			if err != nil {
				journal(critical, "importMain", "Erreur à l'importation du type: "+fnName)
			}
			journal(info, "importMain", "Fin de l'import du type: "+fnName+" pour le batch "+batch.ID.Key)
		}
	} else {
		journal(critical, "importMain", "Le lot "+batch.ID.Key+" est fermé, import impossible.")
	}
}
