package exportdatapi

import (
	"log"

	daclient "github.com/signaux-faibles/datapi/client"
)

var urssafMapping = map[string]string{
	"116": "Île-de-France",
	"117": "Île-de-France",
	"200": "Corse",
	"217": "Champagne-Ardenne",
	"227": "Picardie",
	"237": "Haute-Normandie",
	"247": "Centre-Val de Loire",
	"257": "Basse-Normandie",
	"267": "Bourgogne",
	"311": "Midi-Pyrénées",
	"317": "Nord-Pas-de-Calais",
	"417": "Lorraine",
	"427": "Alsace",
	"437": "Franche-Comté",
	"451": "Centre-Val de Loire",
	"480": "Languedoc-Roussillon",
	"527": "Pays de la Loire",
	"537": "Bretagne",
	"547": "Poitou-Charentes",
	"595": "Nord-Pas-de-Calais",
	"693": "Rhône-Alpes",
	"727": "Aquitaine",
	"737": "Midi-Pyrénées",
	"747": "Limousin",
	"748": "Limousin",
	"827": "Rhône-Alpes",
	"837": "Auvergne",
	"917": "Languedoc-Roussillon",
	"937": "Provence-Alpes-Côte d'Azur",
}

// UrssafScope fournit l'identifiant du scope Urssaf
func UrssafScope(compte string, departement string) string {
	if len(compte) < 3 {
		log.Println("compte trop court: '" + compte + "', fallback sur le département")
		return "URSSAF " + regionFromDepartement(departement)
	}

	if scope, ok := urssafMapping[compte[0:3]]; ok {
		return "URSSAF " + scope
	}

	log.Println("pas de correspondance renseignée pour " + compte[0:3] + " (lib/exportdatapi/urssaf.go)")
	return "URSSAF " + regionFromDepartement(departement)
}

func listeDepartement(departement map[string]string) (liste []string) {
	for k := range departement {
		liste = append(liste, k)
	}
	return liste
}

func urssafPolicies() []daclient.Object {
	var urssafPolicies []daclient.Object

	for k := range reverseMap(urssafMapping) {
		urssafPolicies = append(urssafPolicies, daclient.Object{
			Key: map[string]string{
				"type": "policy",
				"name": "Accès URSSAF " + k,
			},
			Value: map[string]interface{}{
				"match":   "(public|reference)",
				"scope":   []string{"Agent URSSAF"},
				"promote": listeDepartement(departement),
			},
		})
	}
	return urssafPolicies
}
