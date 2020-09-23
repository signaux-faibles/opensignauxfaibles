package marshal

import (
	"errors"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// UrssafToDate convertit le format de date urssaf en type Date.
// Les dates urssaf sont au format YYYMMJJ tels que YYY = YYYY - 1900 (e.g: 118 signifie
// 2018)
func urssafToDate(urssaf string) (time.Time, error) {

	intUrsaff, err := strconv.Atoi(urssaf)
	if err != nil {
		return time.Time{}, base.NewCriticError(errors.New("Valeur non autorisée pour une conversion en date: "+urssaf), "fatal")
	}
	strDate := strconv.Itoa(intUrsaff + 19000000)
	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return time.Time{}, base.NewCriticError(errors.New("Valeur non autorisée pour une conversion en date: "+urssaf), "fatal")
	}

	return date, nil
}
