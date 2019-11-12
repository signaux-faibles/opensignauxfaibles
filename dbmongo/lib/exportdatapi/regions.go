package exportdatapi

import (
	daclient "github.com/signaux-faibles/datapi/client"
)

var region = map[string][]string{
	"Provence-Alpes-Côte d'Azur": []string{"04", "05", "06", "13", "83", "84"},
	"Bourgogne-Franche-Comté":    []string{"25", "70", "39", "90", "21", "58", "71", "89"},
	"Pays de la Loire":           []string{"44", "49", "53", "72", "85"},
	"Auvergne-Rhône-Alpes":       []string{"01", "03", "07", "15", "26", "38", "42", "43", "63", "69", "73", "74"},
	"Nouvelle-Aquitaine":         []string{"16", "17", "19", "23", "24", "33", "40", "47", "64", "79", "86", "87"},
	"Guadeloupe":                 []string{"971"},
	"Martinique":                 []string{"972"},
	"Guyane":                     []string{"973"},
	"La Réunion":                 []string{"974"},
	"Mayotte":                    []string{"976"},
	"Île-de-France":              []string{"75", "77", "78", "91", "92", "93", "94", "95"},
	"Centre-Val de Loire":        []string{"18", "28", "36", "37", "41", "45"},
	"Normandie":                  []string{"14", "27", "50", "61", "76"},
	"Hauts-de-France":            []string{"02", "59", "60", "62", "80"},
	"Grand Est":                  []string{"08", "10", "51", "52", "54", "55", "57", "67", "68", "88"},
	"Bretagne":                   []string{"22", "29", "35", "56"},
	"Occitanie":                  []string{"09", "11", "12", "30", "31", "32", "34", "46", "48", "65", "66", "81", "82"},
	"Corse":                      []string{"2A", "2B"},
}

var lienRegion = map[string][]string{
	"Bourgogne-Franche-Comté": []string{"Bourgogne", "Franche-Comté"},
	"Auvergne-Rhône-Alpes":    []string{"Auvergne", "Rhône-Alpes"},
	"Nouvelle-Aquitaine":      []string{"Aquitaine", "Limousin", "Poitou-Charentes"},
	"Grand-Est":               []string{"Alsace", "Lorraine", "Champagne-Ardenne"},
}

var ancienneRegion = map[string][]string{
	"Alsace":               []string{"67", "68"},
	"Aquitaine":            []string{"24", "33", "40", "47", "64"},
	"Auvergne":             []string{"03", "15", "43", "63"},
	"Basse-Normandie":      []string{"14", "50", "61"},
	"Bourgogne":            []string{"21", "58", "71", "89"},
	"Bretagne":             []string{"22", "29", "35", "56"},
	"Centre":               []string{"18", "28", "36", "37", "41", "45"},
	"Champagne-Ardenne":    []string{"08", "10", "51", "52"},
	"Franche-Comté":        []string{"25", "70", "39", "90"},
	"Haute-Normandie":      []string{"27", "76"},
	"Languedoc-Roussillon": []string{"30", "34", "48", "66"},
	"Limousin":             []string{"19", "23", "87"},
	"Lorraine":             []string{"54", "55", "57", "88"},
	"Midi-Pyrénées":        []string{"09", "12", "31", "32", "46", "65", "81", "82"},
	"Nord-Pas-de-Calais":   []string{"59", "62"},
	"Picardie":             []string{"02", "60", "80"},
	"Poitou-Charentes":     []string{"16", "17", "79", "86"},
	"Rhône-Alpes":          []string{"01", "07", "26", "38", "42", "69", "73", "74"},
}

var departement = map[string]string{
	"01":  "Ain",
	"02":  "Aisne",
	"03":  "Allier",
	"04":  "Alpes-de-Haute-Provence",
	"05":  "Hautes-Alpes",
	"06":  "Alpes-Maritimes",
	"07":  "Ardèche",
	"08":  "Ardennes",
	"09":  "Ariège",
	"10":  "Aube",
	"11":  "Aude",
	"12":  "Aveyron",
	"13":  "Bouches-du-Rhône",
	"14":  "Calvados",
	"15":  "Cantal",
	"16":  "Charente",
	"17":  "Charente-Maritime",
	"18":  "Cher",
	"19":  "Corrèze",
	"21":  "Côte-d'Or",
	"22":  "Côtes-d'Armor",
	"23":  "Creuse",
	"24":  "Dordogne",
	"25":  "Doubs",
	"26":  "Drôme",
	"27":  "Eure",
	"28":  "Eure-et-Loir",
	"29":  "Finistère",
	"2A":  "Corse-du-Sud",
	"2B":  "Haute-Corse",
	"30":  "Gard",
	"31":  "Haute-Garonne",
	"32":  "Gers",
	"33":  "Gironde",
	"34":  "Hérault",
	"35":  "Ille-et-Vilaine",
	"36":  "Indre",
	"37":  "Indre-et-Loire",
	"38":  "Isère",
	"39":  "Jura",
	"40":  "Landes",
	"41":  "Loir-et-Cher",
	"42":  "Loire",
	"43":  "Haute-Loire",
	"44":  "Loire-Atlantique",
	"45":  "Loiret",
	"46":  "Lot",
	"47":  "Lot-et-Garonne",
	"48":  "Lozère",
	"49":  "Maine-et-Loire",
	"50":  "Manche",
	"51":  "Marne",
	"52":  "Haute-Marne",
	"53":  "Mayenne",
	"54":  "Meurthe-et-Moselle",
	"55":  "Meuse",
	"56":  "Morbihan",
	"57":  "Moselle",
	"58":  "Nièvre",
	"59":  "Nord",
	"60":  "Oise",
	"61":  "Orne",
	"62":  "Pas-de-Calais",
	"63":  "Puy-de-Dôme",
	"64":  "Pyrénées-Atlantiques",
	"65":  "Hautes-Pyrénées",
	"66":  "Pyrénées-Orientales",
	"67":  "Bas-Rhin",
	"68":  "Haut-Rhin",
	"69":  "Rhône",
	"70":  "Haute-Saône",
	"71":  "Saône-et-Loire",
	"72":  "Sarthe",
	"73":  "Savoie",
	"74":  "Haute-Savoie",
	"75":  "Paris",
	"76":  "Seine-Maritime",
	"77":  "Seine-et-Marne",
	"78":  "Yvelines",
	"79":  "Deux-Sèvres",
	"80":  "Somme",
	"81":  "Tarn",
	"82":  "Tarn-et-Garonne",
	"83":  "Var",
	"84":  "Vaucluse",
	"85":  "Vendée",
	"86":  "Vienne",
	"87":  "Haute-Vienne",
	"88":  "Vosges",
	"89":  "Yonne",
	"90":  "Territoire de Belfort",
	"91":  "Essonne",
	"92":  "Hauts-de-Seine",
	"93":  "Seine-Saint-Denis",
	"94":  "Val-de-Marne",
	"95":  "Val-d'Oise",
	"971": "Guadeloupe",
	"972": "Martinique",
	"973": "Guyane",
	"974": "La Réunion",
	"976": "Mayotte",
}

// GetRegions retournes les objets de référentiel liés aux régions.
func GetRegions(batch string, algo string) (regions []daclient.Object) {
	for r, d := range region {
		regions = append(regions, daclient.Object{
			Key: map[string]string{
				"type":   "region",
				"region": r,
				"batch":  batch + "." + algo,
			},
			Scope: []string{r},
			Value: map[string]interface{}{
				"departements": d,
			},
		})
	}

	for r, d := range ancienneRegion {
		regions = append(regions, daclient.Object{
			Key: map[string]string{
				"type":   "region",
				"region": r,
				"batch":  batch + "." + algo,
			},
			Scope: []string{r},
			Value: map[string]interface{}{
				"departements": d,
			},
		})
	}

	var franceEntiere []string
	for n, d := range departement {
		franceEntiere = append(franceEntiere, n)
		regions = append(regions, daclient.Object{
			Key: map[string]string{
				"type":  "departements",
				"batch": batch + "." + algo,
			},
			Scope: []string{n},
			Value: map[string]interface{}{
				n: d,
			},
		})
	}

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "France entière",
			"batch":  batch + "." + algo,
		},
		Scope: []string{"France entière"},
		Value: map[string]interface{}{
			"departements": franceEntiere,
		},
	})

	return regions
}
