package urssaf

import (
	"errors"
	"strconv"
	"time"
)

// UrssafToDate converts the URSSAF date format to Date type.
// Les dates urssaf sont au format YYYMMJJ tels que YYY = YYYY - 1900 (e.g: 118 signifie
// 2018)
func UrssafToDate(urssaf string) (time.Time, error) {

	intUrsaff, err := strconv.Atoi(urssaf)
	if err != nil {
		return time.Time{}, errors.New("invalid value for date conversion: " + urssaf)
	}
	strDate := strconv.Itoa(intUrsaff + 19000000)
	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return time.Time{}, errors.New("invalid value for date conversion: " + urssaf)
	}

	return date, nil
}

// UrssafToPeriod converts the URSSAF period format to Period type. On trouve ces
// périodes formatées en 4 ou 6 caractère (YYQM ou YYYYQM).
// si YY < 50 alors YYYY = 20YY sinon YYYY = 19YY.
// si QM == 62 alors période annuelle sur YYYY.
// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY.
// si 0 < M < 4 alors mois M du trimestre Q.
func UrssafToPeriod(urssaf string) (start time.Time, end time.Time, err error) {

	if len(urssaf) == 4 {
		if urssaf[0:2] < "50" {
			urssaf = "20" + urssaf
		} else {
			urssaf = "19" + urssaf
		}
	}

	if len(urssaf) != 6 {
		return start, end, errors.New("invalid value")
	}

	year, err := strconv.Atoi(urssaf[0:4])
	if err != nil {
		return start, end, errors.New("invalid value")
	}

	if urssaf[4:6] == "62" {
		start = time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		quarter, err := strconv.Atoi(urssaf[4:5])
		if err != nil {
			return start, end, err
		}
		monthOfQuarter, err := strconv.Atoi(urssaf[5:6])
		if err != nil {
			return start, end, err
		}
		if monthOfQuarter == 0 {
			start = time.Date(year, time.Month((quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
			end = time.Date(year, time.Month((quarter-1)*3+4), 1, 0, 0, 0, 0, time.UTC)
		} else {
			start = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter), 1, 0, 0, 0, 0, time.UTC)
			end = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter+1), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	return start, end, nil
}
