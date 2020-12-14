package misc

import (
	"errors"
	"flag"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

// func Test_UrssafToDate(t *testing.T) {
// 	a, e := UrssafToDate("1180101")
// 	if a == time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC) && e == nil {
// 		t.Log("UrssafToDate: 1180101 -> 1er janvier 2018: OK")
// 	} else {
// 		t.Error("UrssafToDate: 1180101 -> 1er janvier 2018: Fail")
// 	}

// 	a, e = UrssafToDate("11a0101")
// 	z := time.Time{}
// 	if a == z && e != nil {
// 		t.Log("UrssafToDate: 11a0101 -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToDate: 1180101 -> erreur: Fail")
// 	}

// 	a, e = UrssafToDate("1180151")
// 	if a == z && e != nil {
// 		t.Log("UrssafToDate: 1180151 -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToDate: 1180151 -> erreur: Fail")
// 	}
// }
// func Test_UrssafToPeriod(t *testing.T) {
// 	a, e := UrssafToPeriod("1862")
// 	b := Periode{
// 		Start: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
// 		End:   time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
// 	}
// 	if a == b && e == nil {
// 		t.Log("UrssafToPeriod: 1862 -> l'année 2018: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 1862 -> l'année 2018: Fail")
// 	}

// 	b = Periode{
// 		Start: time.Date(2018, time.Month(4), 1, 0, 0, 0, 0, time.UTC),
// 		End:   time.Date(2018, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
// 	}
// 	a, e = UrssafToPeriod("1820")
// 	if a == b && e == nil {
// 		t.Log("UrssafToPeriod: 1820 -> 2° trimestre 2018: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 1820 -> 2° trimestre 2018: Fail")
// 	}

// 	b = Periode{
// 		Start: time.Date(1963, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
// 		End:   time.Date(1963, time.Month(8), 1, 0, 0, 0, 0, time.UTC),
// 	}
// 	a, e = UrssafToPeriod("6331")
// 	if a == b && e == nil {
// 		t.Log("UrssafToPeriod: 6331 -> Juillet 1963: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 6331 -> Juillet 1963: Fail")
// 	}

// 	b = Periode{
// 		Start: time.Time{},
// 		End:   time.Time{},
// 	}
// 	a, e = UrssafToPeriod("56331")
// 	if a == b && e != nil {
// 		t.Log("UrssafToPeriod: 56331 -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 56331 -> erreur: Fail")
// 	}

// 	a, e = UrssafToPeriod("56a1")
// 	if a == b && e != nil {
// 		t.Log("UrssafToPeriod: 56a1 -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 56a1 -> erreur: Fail")
// 	}

// 	a, e = UrssafToPeriod("5a31")
// 	if a == b && e != nil {
// 		t.Log("UrssafToPeriod: 5a31 -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 5a31 -> erreur: Fail")
// 	}

// 	a, e = UrssafToPeriod("564a")
// 	if a == b && e != nil {
// 		t.Log("UrssafToPeriod: 564a -> erreur: OK")
// 	} else {
// 		t.Error("UrssafToPeriod: 56564aa1 -> erreur: Fail")
// 	}

// }
func Test_ParsePInt(t *testing.T) {
	a, err1 := ParsePInt("9")
	b := 9

	if *a == b && err1 == nil {
		t.Log("ParsePInt: parse 9 -> OK")
	} else {
		t.Error("ParsePInt: parse 9 -> Fail: " + err1.Error())
	}

	c, err2 := ParsePInt("a")
	if *c == 0 && err2 != nil {
		t.Log("ParsePInt: n'est pas un entier -> OK")
	} else {
		t.Error("ParsePInt: n'est pas un entier -> Fail")
	}

	d, err2 := ParsePInt("")
	if d == nil && err2 == nil {
		t.Log("ParsePInt: une valeur vide -> OK")
	} else {
		t.Error("ParsePInt: une valeur vide -> Fail")
	}
}

func Test_ParsePFloat(t *testing.T) {
	a, err1 := ParsePFloat("0.62")
	b := float64(0.62)

	if *a == b && err1 == nil {
		t.Log("ParsePFloat: parse 0.62 -> OK")
	} else {
		t.Error("ParsePFloat: parse 0.62 -> Fail: " + err1.Error())
	}

	c, err2 := ParsePFloat("abcd")
	if *c == 0.0 && err2 != nil {
		t.Log("ParsePFloat: n'est pas un float -> OK")
	} else {
		t.Error("ParsePFloat: n'est pas un float -> Fail")
	}

	d, err2 := ParsePFloat("")
	if d == nil && err2 == nil {
		t.Log("ParsePFloat: une valeur vide -> OK")
	} else {
		t.Error("ParsePFloat: une valeur vide -> Fail")
	}
}

func Test_Max(t *testing.T) {
	if Max(-4, 2) == 2 && Max(2, -4) == 2 {
		t.Log("Max: OK")
	} else {
		t.Error("Max: Fail")
	}
}

func Test_AllErrors(t *testing.T) {
	e := errors.New("erreur d'exemple")
	a := []error{e, e, e}
	b := []error{nil, nil, e}
	c := []error{nil, nil, nil}
	if AllErrors(a, e) {
		t.Log("AllErrors: tous les éléments du tableau sont des erreurs spécifiques: OK")
	} else {
		t.Error("AllErrors: tous les éléments du tableau sont des erreurs spécifiques: Fail")
	}

	if AllErrors(b, e) {
		t.Error("AllErrors: un seul élément du tableau est une erreur: Fail")
	} else {
		t.Log("AllErrors: un seul élément du tableau est une erreur: OK")
	}

	if AllErrors(b, nil) {
		t.Error("AllErrors: un seul élément du tableau n'est pas nil: Fail")
	} else {
		t.Log("AllErrors: un seul élément du tableau n'est pas nil: OK")
	}

	if AllErrors(c, nil) {
		t.Log("AllErrors: tous les éléments du tableau sont nil: OK")
	} else {
		t.Error("AllErrors: tous les éléments du tableau sont nil: Fail")
	}
}

func Test_ExcelToTime(t *testing.T) {
	date := "43101"
	goDate := time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	excelDate, err := ExcelToTime(date)
	if excelDate.UTC() == goDate && err == nil {
		t.Log("ExcelDate: 43101 -> 01/01/2018: OK")
	} else {
		t.Error("ExcelDate: 43101 -> 01/01/2018: Fail")
	}

	date = "43a01"
	goDate = time.Time{}
	excelDate, err = ExcelToTime(date)
	if excelDate.UTC() == goDate && err != nil {
		t.Log("ExcelDate: 43a01 -> erreur: OK")
	} else {
		t.Error("ExcelDate: 43a01 -> erreur: Fail")
	}
}
func Test_SliceIndex(t *testing.T) {
	slice := []string{"a", "b", "c"}

	k := SliceIndex(len(slice), func(i int) bool { return slice[i] == "b" })
	l := SliceIndex(len(slice), func(i int) bool { return slice[i] == "d" })
	if k == 1 {
		t.Log("SliceIndex: recherche un élément existant: OK")
	} else {
		t.Error("SliceIndex: recherche un élément existant: Fail")
	}

	if l == -1 {
		t.Log("SliceIndex: recherche un élément manquant: OK")
	} else {
		t.Error("SliceIndex: recherche un élément manquant: Fail")
	}
}
