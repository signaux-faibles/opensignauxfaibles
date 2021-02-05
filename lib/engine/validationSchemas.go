/* Generated by bundle-schemas.go - DO NOT EDIT */

package engine

var validationSchemas = map[string]string{
"apconso.schema.json": `{
  "title": "EntréeApConso",
  "description": "Champs importés par le parseur lib/apconso/main.go de sfdata.",
  "bsonType": "object",
  "required": ["id_conso", "periode", "heure_consomme"],
  "properties": {
    "id_conso": {
      "bsonType": "string"
    },
    "heure_consomme": {
      "bsonType": "number"
    },
    "montant": {
      "bsonType": "number"
    },
    "effectif": {
      "bsonType": "number"
    },
    "periode": {
      "bsonType": "date"
    }
  },
  "additionalProperties": false
}
`,
"apdemande.schema.json": `{
  "title": "EntréeApDemande",
  "description": "Champs importés par le parseur lib/apdemande/main.go de sfdata.",
  "bsonType": "object",
  "required": ["id_demande", "periode", "hta", "motif_recours_se"],
  "properties": {
    "id_demande": {
      "bsonType": "string"
    },
    "periode": {
      "bsonType": "object",
      "required": ["start", "end"],
      "properties": {
        "start": { "bsonType": "date" },
        "end": { "bsonType": "date" }
      },
      "additionalProperties": false
    },
    "hta": {
      "description": "Nombre total d'heures autorisées",
      "bsonType": "number"
    },
    "motif_recours_se": {
      "description": "Cause d'activité partielle",
      "bsonType": "number"
    },
    "effectif_entreprise": {
      "bsonType": "number"
    },
    "effectif": {
      "bsonType": "number"
    },
    "date_statut": {
      "bsonType": "date"
    },
    "mta": {
      "bsonType": "number"
    },
    "effectif_autorise": {
      "bsonType": "number"
    },
    "heure_consommee": {
      "bsonType": "number"
    },
    "montant_consommee": {
      "bsonType": "number"
    },
    "effectif_consomme": {
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"bdf.schema.json": `{
  "title": "EntréeBdf",
  "description": "Note: CE SCHEMA EST INCOMPLET POUR L'INSTANT. Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/143",
  "bsonType": "object",
  "required": ["siren"],
  "properties": {
    "siren": {
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    }
  },
  "additionalProperties": false
}
`,
"ccsf.schema.json": `{
  "title": "EntréeCcsf",
  "description": "Champs importés par le parseur lib/urssaf/ccsf.go de sfdata.",
  "bsonType": "object",
  "required": ["date_traitement", "stade", "action"],
  "properties": {
    "date_traitement": {
      "bsonType": "date",
      "description": "Date de début de la procédure CCSF"
    },
    "stade": {
      "bsonType": "string",
      "description": "TODO: choisir un type plus précis"
    },
    "action": {
      "bsonType": "string",
      "description": "TODO: choisir un type plus précis"
    }
  },
  "additionalProperties": false
}
`,
"compte.schema.json": `{
  "title": "EntréeCompte",
  "description": "Champs importés par le parseur lib/urssaf/compte.go de sfdata.",
  "bsonType": "object",
  "required": ["periode", "siret", "numero_compte"],
  "properties": {
    "periode": {
      "bsonType": "date",
      "description": "Date à laquelle cet établissement est associé à ce numéro de compte URSSAF."
    },
    "siret": {
      "bsonType": "string",
      "description": "Numéro SIRET de l'établissement. Les numéros avec des Lettres sont des sirets provisoires."
    },
    "numero_compte": {
      "bsonType": "string",
      "description": "Compte administratif URSSAF."
    }
  },
  "additionalProperties": false
}
`,
"delai.schema.json": `{
  "title": "EntréeDelai",
  "description": "Champs importés par le parseur lib/urssaf/delai.go de sfdata.",
  "bsonType": "object",
  "required": [
    "numero_compte",
    "numero_contentieux",
    "date_creation",
    "date_echeance",
    "duree_delai",
    "denomination",
    "indic_6m",
    "annee_creation",
    "montant_echeancier",
    "stade",
    "action"
  ],
  "properties": {
    "numero_compte": {
      "bsonType": "string",
      "description": "Compte administratif URSSAF."
    },
    "numero_contentieux": {
      "bsonType": "string",
      "description": "Le numéro de structure est l'identifiant d'un dossier contentieux."
    },
    "date_creation": {
      "bsonType": "date",
      "description": "Date de création du délai."
    },
    "date_echeance": {
      "bsonType": "date",
      "description": "Date d'échéance du délai."
    },
    "duree_delai": {
      "bsonType": "number",
      "description": "Durée du délai en jours: nombre de jours entre date_creation et date_echeance.",
      "minimum": 1
    },
    "denomination": {
      "bsonType": "string",
      "description": "Raison sociale de l'établissement."
    },
    "indic_6m": {
      "bsonType": "string",
      "description": "Délai inférieur ou supérieur à 6 mois ? Modalités INF et SUP."
    },
    "annee_creation": {
      "bsonType": "number",
      "description": "Année de création du délai."
    },
    "montant_echeancier": {
      "bsonType": "number",
      "description": "Montant global de l'échéancier, en euros.",
      "minimum": 0.01
    },
    "stade": {
      "bsonType": "string",
      "description": "Code externe du stade."
    },
    "action": {
      "bsonType": "string",
      "description": "Code externe de l'action."
    }
  },
  "additionalProperties": false
}
`,
"detect_invalid_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": "$dataPerHash.v"
    }
  },
  { "$match": { "dataPerHash": { "$not": { "$type": 3 } } } }
]
`,
"effectif.schema.json": `{
  "title": "EntréeEffectif",
  "description": "Champs importés par le parseur lib/urssaf/effectif.go de sfdata.",
  "bsonType": "object",
  "required": ["numero_compte", "periode", "effectif"],
  "properties": {
    "numero_compte": {
      "description": "Compte administratif URSSAF.",
      "bsonType": "string"
    },
    "periode": {
      "bsonType": "date"
    },
    "effectif": {
      "description": "Nombre de personnes employées par l'établissement.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"effectif_ent.schema.json": `{
  "title": "EntréeEffectifEnt",
  "description": "Champs importés par le parseur lib/urssaf/effectif_ent.go de sfdata.",
  "bsonType": "object",
  "required": ["periode", "effectif"],
  "properties": {
    "periode": {
      "bsonType": "date"
    },
    "effectif": {
      "description": "Nombre de personnes employées par l'entreprise.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"flatten_data_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": { "$objectToArray": "$dataPerHash.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": 1,
      "dataHash": "$dataPerHash.k",
      "dataObject": "$dataPerHash.v"
    }
  }
]
`,
"procol.schema.json": `{
  "title": "EntréeDéfaillances",
  "description": "Champs importés par le parseur lib/urssaf/procol.go de sfdata.",
  "bsonType": "object",
  "required": ["action_procol", "stade_procol", "date_effet"],
  "properties": {
    "action_procol": {
      "description": "Nature de la procédure de défaillance.",
      "bsonType": "string",
      "enum": ["liquidation", "redressement", "sauvegarde"]
    },
    "stade_procol": {
      "description": "Evénement survenu dans le cadre de cette procédure.",
      "bsonType": "string",
      "enum": [
        "abandon_procedure",
        "solde_procedure",
        "fin_procedure",
        "plan_continuation",
        "ouverture",
        "inclusion_autre_procedure",
        "cloture_insuffisance_actif"
      ]
    },
    "date_effet": {
      "bsonType": "date",
      "description": "Date effet de la procédure collective."
    }
  },
  "additionalProperties": false
}
`,
}
