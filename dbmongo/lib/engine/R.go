package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gin-gonic/gin"

	"github.com/senseyeio/roger"
)

// Algo1 retour de la prédiction de l'algorithme 1
type Algo1 struct {
	Siret      string  `json:"siret" bson:"siret"`
	Prediction float64 `json:"prediction_0_12" bson:"prediction_0_12"`
}

func algo1(c *gin.Context) {
	Rscript, _ := ioutil.ReadFile("R/algo1.R")
	rClient, err := roger.NewRClient("127.0.0.1", 6311)
	rSession, _ := rClient.GetSession()
	defer rSession.Close()

	err = rSession.Assign("periode_train", "2015-01-01 01:00:00")
	err = rSession.Assign("periode_test", "2016-01-01 01:00:00")
	err = rSession.Assign("periode_actual", "2018-02-01 01:00:00")

	if err != nil {
		fmt.Println("Failed")
		return
	}

	value, err := rSession.Eval(string(Rscript))
	if err != nil {
		fmt.Println("Command failed: " + err.Error())
	}

	var result []Algo1

	err = json.Unmarshal([]byte(value.(string)), &result)

	c.JSON(200, result)
}

// R démarre un processus R
func r() {
	cmd := exec.Command("Rscript", "R/rserve.R")
	cmd.Run()
}
