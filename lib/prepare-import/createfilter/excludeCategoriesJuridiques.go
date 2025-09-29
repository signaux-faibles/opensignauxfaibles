package createfilter

import (
	"context"
	"github.com/signaux-faibles/goSirene"
	"strings"
)

func isExcludedCategorieJuridique(categorieJuridique string) bool {
	var excludedCategoriesJuridiques = []string{
		"7490",
		"7430",
		"7470",
		"7410",
		"7379",
		"7348",
		"7346",
		"7210",
		"7220",
		"4140",
		"7373",
		"7366",
		"7389",
		"4110",
		"4120",
		"7383",
		"4160",
	}
	for _, excluded := range excludedCategoriesJuridiques {
		if categorieJuridique == excluded {
			return true
		}
	}
	return false
}

func isExcludedActivity(activity string) bool {
	var excludedActivityPrefix = []string{
		"84",
		"85",
	}
	for _, excluded := range excludedActivityPrefix {
		if strings.HasPrefix(activity, excluded) {
			return true
		}
	}
	return false
}

func readExcludedSirens(path string) map[string]struct{} {
	sireneUL := goSirene.SireneULParser(context.Background(), path)
	var excludedSirens = make(map[string]struct{})
	for s := range sireneUL {
		if isExcludedCategorieJuridique(s.CategorieJuridiqueUniteLegale) {
			excludedSirens[s.Siren] = struct{}{}
		}
		if isExcludedActivity(s.ActivitePrincipaleUniteLegale) {
			excludedSirens[s.Siren] = struct{}{}
		}
	}
	return excludedSirens
}

func CategorieJuridiqueFilter(path string) filter {
	var excludedSirens = readExcludedSirens(path)
	return func(siren string) bool {
		_, ok := excludedSirens[siren]
		return !ok
	}
}
