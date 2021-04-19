// Package misc fournit les fonctions et types habituels dans le projet
// opensignauxfaibles
package misc

import (
	"errors"
	"strconv"
	"time"
)

// ParsePInt parse un entier et retourne un pointeur
func ParsePInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(s)
	return &i, err
}

// ParsePFloat parse un flottant et retourne un pointeur
func ParsePFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.ParseFloat(s, 64)
	return &i, err
}

// Max retourne le plus grand des deux entiers passés en argument
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AllErrors retourne true si le tableau contient au moins une erreur non nil
func AllErrors(slice []error, item interface{}) bool {
	for _, i := range slice {
		if i != item {
			return false
		}
	}
	return true
}

// ExcelToTime convertit une date excel en time.Time
func ExcelToTime(excel string) (time.Time, error) {
	excelInt, err := strconv.ParseInt(excel, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("valeur non autorisée")
	}
	return time.Unix((excelInt-25569)*3600*24, 0), nil
}

// SliceIndex retourne la position du premier élément qui satisfait le
// prédicat, -1 si aucun élément n'est trouvé.
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Periode est un type temporel avec un début et une fin employé dans les types et
// fonctions opensignauxfaibles manipulant des périodes temporelles. La date de fin
// est exclue de la période.
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

// GenereSeriePeriode génère une liste de dates pour les mois entre la date de début (incluse) et la date de fin (exclue)
func GenereSeriePeriode(debut time.Time, fin time.Time) []time.Time {
	var serie []time.Time
	for fin.After(debut) {
		serie = append(serie, debut)
		debut = debut.AddDate(0, 1, 0)
	}
	return serie
}

// GenereSeriePeriodeAnnuelle génère une liste de dates pour les années entre le début de l'année de la date de début et la fin de l'année de la date de fin
func GenereSeriePeriodeAnnuelle(debut time.Time, fin time.Time) []int {
	var serie []int
	for debut.Year() <= fin.Year() {
		serie = append(serie, debut.Year())
		debut = debut.AddDate(1, 0, 0)
	}
	return serie
}
