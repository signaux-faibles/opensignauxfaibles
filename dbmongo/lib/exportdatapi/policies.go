package exportdatapi

import daclient "github.com/signaux-faibles/datapi/client"

// GetPolicies returns policies to be exported in datapi
func GetPolicies(batch string) []daclient.Object {
	var policies []daclient.Object

	var accesFranceEntiere []string
	var accesFranceEntiereAncien []string
	for d := range departement {
		accesFranceEntiere = append(accesFranceEntiere, d)
		accesFranceEntiereAncien = append(accesFranceEntiereAncien, d)
	}

	for r := range region {
		accesFranceEntiere = append(accesFranceEntiere, r)
	}

	for r := range ancienneRegion {
		accesFranceEntiereAncien = append(accesFranceEntiereAncien, r)
	}

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès France Entière",
		},
		Value: map[string]interface{}{
			"match":   "(public|reference)",
			"scope":   []string{"France entière"},
			"promote": accesFranceEntiere,
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès France Entière Ancienne région",
		},
		Value: map[string]interface{}{
			"match":   "(public|reference)",
			"scope":   []string{"France entière anciennes régions"},
			"promote": accesFranceEntiereAncien,
		},
	})

	for r, d := range region {
		policies = append(policies, daclient.Object{
			Key: map[string]string{
				"type": "policy",
				"name": "Accès " + r,
			},
			Value: map[string]interface{}{
				"match":   "(public|reference)",
				"scope":   []string{r},
				"promote": d,
			},
		})
	}

	for r, d := range ancienneRegion {
		policies = append(policies, daclient.Object{
			Key: map[string]string{
				"type": "policy",
				"name": "Accès " + r,
			},
			Value: map[string]interface{}{
				"match":   "(public|reference)",
				"scope":   []string{r},
				"promote": d,
			},
		})
	}

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Limitation en écriture",
		},
		Value: map[string]interface{}{
			"match":  ".*",
			"key":    map[string]string{},
			"writer": []string{"datawriter"},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Limitation Accès",
		},
		Value: map[string]interface{}{
			"match":  "system",
			"key":    map[string]string{},
			"writer": []string{"manager"},
			"reader": []string{"manager"},
		},
	})

	policies = append(policies, urssafPolicies()...)
	return policies
}
