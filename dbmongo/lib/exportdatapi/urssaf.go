package exportdatapi

import "errors"

// // Provence-Alpes-Côte d'Azur
// "Provence-Alpes-Côte d'Azur" "937"

// // Bourgogne-Franche-Comté
// "Bourgogne"		"267"
// "Franche Comté"		"437"

// // Pays de la Loire
// "Pays de la Loire"		"527"

// // Auvergne-Rhône-Alpes
// "Auvergne"		"837"
// "Rhône Alpes"		"827" "693"

// // Nouvelle-Aquitaine
// "Aquitaine"		"727"
// "Limousin"		"747"
// "Poitou-Charentes"		"547"

// // Île-de-France
// "Ile de France"		"117"  "116"

// // Centre-Val de Loire
// "Centre-Val de Loire"		"247" "451"

// // Normandie
// "Basse-Normandie": "257"
// "Haute Normandie": "237"

// // "Hauts-de-France"
// "Nord Pas de Calais"		"317" "595"
// "Picardie"		"227"

// // Grand Est
// "Alsace"		"427"
// "Champagne-Ardenne"		"217"
// "Lorraine"		"417"

// // Occitanie
// "Midi-Pyrénées"		"737"  "311"

// // Bretagne
// "Bretagne"		"537"

// // Corse
// "Corse"		"200"

// // Occitanie
// "Languedoc-Roussillon"		"917"

func urssafScope(compte string) (string, error) {
	var urssafMapping = map[string]string{
		"937": "Provence-Alpes-Côte d'Azur",
		"267": "Bourgogne-Franche-Comté",
		"437": "Bourgogne-Franche-Comté",
		"527": "Pays de la Loire",
		"837": "Auvergne-Rhône-Alpes",
		"827": "Auvergne-Rhône-Alpes",
		"693": "Auvergne-Rhône-Alpes",
		"727": "Nouvelle-Aquitaine",
		"747": "Nouvelle-Aquitaine",
		"547": "Nouvelle-Aquitaine",
		"117": "Île-de-France",
		"116": "Île-de-France",
		"247": "Centre-Val de Loire",
		"451": "Centre-Val de Loire",
		"257": "Normandie",
		"237": "Normandie",
		"317": "Hauts-de-France",
		"595": "Hauts-de-France",
		"227": "Hauts-de-France",
		"427": "Grand Est",
		"217": "Grand Est",
		"417": "Grand Est",
		"737": "Occitanie",
		"311": "Occitanie",
		"537": "Bretagne",
		"200": "Corse",
		"917": "Languedoc-Roussillon",
	}

	if len(compte) < 3 {
		return "", errors.New("le numéro de compte est trop court")
	}

	if scope, ok := urssafMapping[compte[0:3]]; ok {
		return scope, nil
	}

	return "", errors.New("aucun ")
}
