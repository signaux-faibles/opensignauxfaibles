package exportdatapi

import (
	"errors"
	"strconv"
	"time"

	daclient "github.com/signaux-faibles/datapi/client"

	"github.com/globalsign/mgo/bson"
)

// GetPipeline construit le pipeline d'aggregation
func GetPipeline(batch string) (pipeline []bson.M) {
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"algo":  "algo",
		"batch": batch,
	}})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"siret":     1,
		"periode":   1,
		"timestamp": -1,
	}})

	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"siret":   "$siret",
				"periode": "$periode",
				"batch":   "$batch",
				"algo":    "$algo",
			},
			"score": bson.M{
				"$first": "$score",
			},
			"alert": bson.M{
				"$first": "$alert",
			},
			"diff": bson.M{
				"$first": "$diff",
			},
			"timestamp": bson.M{
				"$first": "$timestamp",
			},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"siret":   "$_id.siret",
			"periode": "$_id.periode",
			"batch":   "$_id.batch",
			"algo":    "$_id.algo",
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,
		},
	})

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id": bson.D{
			{Name: "scope", Value: "etablissement"},
			{Name: "key", Value: "$siret"},
			{Name: "batch", Value: "$batch"},
		},
		"_idEntreprise": bson.D{
			{Name: "scope", Value: "entreprise"},
			{Name: "key", Value: bson.M{"$substr": []interface{}{"$siret", 0, 9}}},
			{Name: "batch", Value: "$batch"},
		},
		"score": "$score",
		"alert": "$alert",
		"diff":  "$diff",
		"connu": "$connu",
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "etablissement"}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "_idEntreprise",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"etablissement": bson.M{"$arrayElemAt": []interface{}{"$etablissement", 0}},
		"entreprise":    bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}},
	}})

	return pipeline
}

// Detection correspond aux données retournées pour l'export Datapi
type Detection struct {
	ID            map[string]string `json:"_id" bson:"_id"`
	Score         float64           `json:"score" bson:"score"`
	Diff          float64           `json:"diff" bson:"diff"`
	Timestamp     time.Time         `json:"timestamp" bson:"timestamp"`
	Alert         string            `json:"alert" bson:"alert"`
	Connu         bool              `json:"connu" bson:"connu"`
	Periode       time.Time         `json:"periode" bson:"periode"`
	Etablissement Etablissement     `json:"etablissement" bson:"etablissement"`
	Entreprise    Entreprise        `json:"entreprise" bson:"entreprise"`
}

// CRP est une ligne de fichier CRP
type CRP struct {
	Siret       string `json:"siret" bson:"siret"`
	Difficulte  string `json:"difficulte" bson:"difficulte"`
	EtatDossier string `json:"etat_dossier" bson:"etat_dossier"`
	Actions     string `json:"actions" bson:"actions"`
	Statut      string `json:"statut" bson:"statut"`
	Fichier     string `json:"fichier" bson:"fichier"`
}

// Procol donne le statut et la date de statut pour une entreprise en matière de procédures collectives
type Procol struct {
	Etat string    `json:"etat" bson:"etat"`
	Date time.Time `json:"date_procol"`
}

// Etablissement is an object
type Etablissement struct {
	ID    map[string]string `bson:"_id"`
	Value struct {
		Sirene          Sirene        `json:"sirene" bson:"sirene"`
		Cotisation      []float64     `json:"cotisation" bson:"cotisation"`
		Debit           []Debit       `json:"debit" bson:"debit"`
		APDemande       []APDemande   `json:"apdemande" bson:"apdemande"`
		APConso         []APConso     `json:"apconso" bson:"apconso"`
		Effectif        []Effectif    `json:"effectif" bson:"effectif"`
		DernierEffectif Effectif      `json:"dernier_effectif" bson:"dernier_effectif"`
		Delai           []interface{} `json:"delai" bson:"delai"`
		Procol          []Procol      `json:"procol" bson:"procol"`
		LastProcol      Procol        `json:"last_procol" bson:"last_procol"`
	} `bson:"value"`
}

// Entreprise object
type Entreprise struct {
	ID    map[string]string `json:"_id" bson:"_id"`
	Value struct {
		Diane    []Diane       `json:"diane" bson:"diane"`
		BDF      []interface{} `json:"bdf" bson:"bdf"`
		SireneUL SireneUL      `json:"sirene_ul" bson:"sirene_ul"`
		CRP      CRP           `json:"crp" bson:"crp"`
	} `bson:"value"`
}

// Effectif detail
type Effectif struct {
	Periode  time.Time `json:"periode" bson:"periode"`
	Effectif int       `json:"effectif" bson:"effectif"`
}

// Debit detail
type Debit struct {
	PartOuvriere  float64   `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale float64   `json:"part_patronale" bson:"part_patronale"`
	Periode       time.Time `json:"periode" bson:"periode"`
}

// APConso detail
type APConso struct {
	IDConso       string    `json:"id_conso" bson:"id_conso"`
	HeureConsomme float64   `json:"heure_consomme" bson:"heure_consomme"`
	Montant       float64   `json:"montant" bson:"montant"`
	Effectif      int       `json:"int" bson:"int"`
	Periode       time.Time `json:"periode" bson:"periode"`
}

// SireneUL detail
type SireneUL struct {
	Siren           string `json:"siren" bson:"siren"`
	RaisonSociale   string `json:"raison_sociale" bson:"raison_sociale"`
	StatutJuridique string `json:"statut_juridique" bson:"statut_juridique"`
}

// APDemande detail
type APDemande struct {
	DateStatut time.Time `json:"date_statut" bson:"date_statut"`
	Periode    struct {
		Start time.Time `json:"start" bson:"start"`
		End   time.Time `json:"end" bson:"end"`
	} `json:"periode" bson:"periode"`
	EffectifAutorise int     `json:"effectif_autorise" bson:"effectif_autorise"`
	EffectifConsomme int     `json:"effectif_consomme" bson:"effectif_consomme"`
	IDDemande        string  `json:"id_conso" bson:"id_conso"`
	Effectif         int     `json:"int" bson:"int"`
	MTA              float64 `json:"mta" bson:"mta"`
	HTA              float64 `json:"hta" bson:"hta"`
	MotifRecoursSE   int     `json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsomme    float64 `json:"heure_consomme" bson:"heure_consomme"`
	Montant          float64 `json:"montant" bson:"montant"`
}

// Diane detail
type Diane struct {
	ChiffreAffaire                  float64   `json:"ca,omitempty" bson:"ca,omitempty"`
	Exercice                        float64   `json:"exercice_diane,omitempty" bson:"exercice_diane,omitempty"`
	NomEntreprise                   string    `json:"nom_entreprise,omitempty" bson:"nom_entreprise,omitempty"`
	NumeroSiren                     string    `json:"numero_siren,omitempty" bson:"numero_siren,omitempty"`
	StatutJuridique                 string    `json:"statut_juridique,omitempty" bson:"statut_juridique,omitempty"`
	ProcedureCollective             bool      `json:"procedure_collective,omitempty" bson:"procedure_collective,omitempty"`
	EffectifConsolide               *int      `json:"effectif_consolide,omitempty" bson:"effectif_consolide,omitempty"`
	DetteFiscaleEtSociale           *float64  `json:"dette_fiscale_et_sociale,omitempty" bson:"dette_fiscale_et_sociale,omitempty"`
	FraisDeRetD                     *float64  `json:"frais_de_RetD,omitempty" bson:"frais_de_RetD,omitempty"`
	ConcesBrevEtDroitsSim           *float64  `json:"conces_brev_et_droits_sim,omitempty" bson:"conces_brev_et_droits_sim,omitempty"`
	NombreEtabSecondaire            *int      `json:"nombre_etab_secondaire,omitempty" bson:"nombre_etab_secondaire,omitempty"`
	NombreFiliale                   *int      `json:"nombre_filiale,omitempty" bson:"nombre_filiale,omitempty"`
	TailleCompoGroupe               *int      `json:"taille_compo_groupe,omitempty" bson:"taille_compo_groupe,omitempty"`
	ArreteBilan                     time.Time `json:"arrete_bilan_diane,omitempty" bson:"arrete_bilan_diane,omitempty"`
	NombreMois                      *int      `json:"nombre_mois,omitempty" bson:"nombre_mois,omitempty"`
	ConcoursBancaireCourant         *float64  `json:"concours_bancaire_courant,omitempty" bson:"concours_bancaire_courant,omitempty"`
	EquilibreFinancier              *float64  `json:"equilibre_financier,omitempty" bson:"equilibre_financier,omitempty"`
	IndependanceFinanciere          *float64  `json:"independance_financiere,omitempty" bson:"independance_financiere,omitempty"`
	Endettement                     *float64  `json:"endettement,omitempty" bson:"endettement,omitempty"`
	AutonomieFinanciere             *float64  `json:"autonomie_financiere,omitempty" bson:"autonomie_financiere,omitempty"`
	DegreImmoCorporelle             *float64  `json:"degre_immo_corporelle,omitempty" bson:"degre_immo_corporelle,omitempty"`
	FinancementActifCirculant       *float64  `json:"financement_actif_circulant,omitempty" bson:"financement_actif_circulant,omitempty"`
	LiquiditeGenerale               *float64  `json:"liquidite_generale,omitempty" bson:"liquidite_generale,omitempty"`
	LiquiditeReduite                *float64  `json:"liquidite_reduite,omitempty" bson:"liquidite_reduite,omitempty"`
	RotationStocks                  *float64  `json:"rotation_stocks,omitempty" bson:"rotation_stocks,omitempty"`
	CreditClient                    *float64  `json:"credit_client,omitempty" bson:"credit_client,omitempty"`
	CreditFournisseur               *float64  `json:"credit_fournisseur,omitempty" bson:"credit_fournisseur,omitempty"`
	CAparEffectif                   *float64  `json:"ca_par_effectif,omitempty" bson:"ca_apar_effectif,omitempty"`
	TauxInteretFinancier            *float64  `json:"taux_interet_financier,omitempty" bson:"taux_interet_financier,omitempty"`
	TauxInteretSurCA                *float64  `json:"taux_interet_sur_ca,omitempty" bson:"taux_interet_sur_ca,omitempty"`
	EndettementGlobal               *float64  `json:"endettement_global,omitempty" bson:"endettement_global,omitempty"`
	TauxEndettement                 *float64  `json:"taux_endettement,omitempty" bson:"taux_endettement,omitempty"`
	CapaciteRemboursement           *float64  `json:"capacite_remboursement,omitempty" bson:"capacite_remboursement,omitempty"`
	CapaciteAutofinancement         *float64  `json:"capacite_autofinancement,omitempty" bson:"capacite_autofinancement,omitempty"`
	CouvertureCaFdr                 *float64  `json:"couverture_ca_fdr,omitempty" bson:"couverture_ca_fdr,omitempty"`
	CouvertureCaBesoinFdr           *float64  `json:"couverture_ca_besoin_fdr,omitempty" bson:"couverture_ca_besoin_fdr,omitempty"`
	PoidsBFRExploitation            *float64  `json:"poids_bfr_exploitation,omitempty" bson:"poids_bfr_exploitation,omitempty"`
	Exportation                     *float64  `json:"exportation,omitempty" bson:"exportation,omitempty"`
	EfficaciteEconomique            *float64  `json:"efficacite_economique,omitempty" bson:"efficacite_economique,omitempty"`
	ProductivitePotentielProduction *float64  `json:"productivite_potentiel_production,omitempty" bson:"productivite_potentiel_production,omitempty"`
	ProductiviteCapitalFinancier    *float64  `json:"productivite_capital_financier,omitempty" bson:"productivite_capital_financier,omitempty"`
	ProductiviteCapitalInvesti      *float64  `json:"productivite_capital_investi,omitempty" bson:"productivite_capital_investi,omitempty"`
	TauxDInvestissementProductif    *float64  `json:"taux_d_investissement_productif,omitempty" bson:"taux_d_investissement_productif,omitempty"`
	RentabiliteEconomique           *float64  `json:"rentabilite_economique,omitempty" bson:"rentabilite_economique,omitempty"`
	Performance                     *float64  `json:"performance,omitempty" bson:"performance,omitempty"`
	RendementBrutFondsPropres       *float64  `json:"rendement_brut_fonds_propres,omitempty" bson:"rendement_brut_fonds_propres,omitempty"`
	RentabiliteNette                *float64  `json:"rentabilite_nette,omitempty" bson:"rentabilite_nette,omitempty"`
	RendementCapitauxPropres        *float64  `json:"rendement_capitaux_propres,omitempty" bson:"rendement_capitaux_propres,omitempty"`
	RendementRessourcesDurables     *float64  `json:"rendement_ressources_durables,omitempty" bson:"rendement_ressources_durables,omitempty"`
	TauxMargeCommerciale            *float64  `json:"taux_marge_commerciale,omitempty" bson:"taux_marge_commerciale,omitempty"`
	TauxValeurAjoutee               *float64  `json:"taux_valeur_ajoutee,omitempty" bson:"taux_valeur_ajoutee,omitempty"`
	PartSalaries                    *float64  `json:"part_salaries,omitempty" bson:"part_salaries,omitempty"`
	PartEtat                        *float64  `json:"part_etat,omitempty" bson:"part_etat,omitempty"`
	PartPreteur                     *float64  `json:"part_preteur,omitempty" bson:"part_preteur,omitempty"`
	PartAutofinancement             *float64  `json:"part_autofinancement,omitempty" bson:"part_autofinancement,omitempty"`
	CAExportation                   *float64  `json:"ca_exportation,omitempty" bson:"ca_exportation,omitempty"`
	AchatMarchandises               *float64  `json:"achat_marchandises,omitempty" bson:"achat_marchandises,omitempty"`
	AchatMatieresPremieres          *float64  `json:"achat_matieres_premieres,omitempty" bson:"achat_matieres_premieres,omitempty"`
	Production                      *float64  `json:"production,omitempty" bson:"production,omitempty"`
	MargeCommerciale                *float64  `json:"marge_commerciale,omitempty" bson:"marge_commerciale,omitempty"`
	Consommation                    *float64  `json:"consommation,omitempty" bson:"consommation,omitempty"`
	AutresAchatsChargesExternes     *float64  `json:"autres_achats_charges_externes,omitempty" bson:"autres_achats_charges_externes,omitempty"`
	ValeurAjoutee                   *float64  `json:"valeur_ajoutee,omitempty" bson:"valeur_ajoutee,omitempty"`
	ChargePersonnel                 *float64  `json:"charge_personnel,omitempty" bson:"charge_personnel,omitempty"`
	ImpotsTaxes                     *float64  `json:"impots_taxes,omitempty" bson:"impots_taxes,omitempty"`
	SubventionsDExploitation        *float64  `json:"subventions_d_exploitation,omitempty" bson:"subventions_d_exploitation,omitempty"`
	ExcedentBrutDExploitation       *float64  `json:"excedent_brut_d_exploitation,omitempty" bson:"excedent_brut_d_exploitation,omitempty"`
	AutresProduitsChargesReprises   *float64  `json:"autres_produits_charges_reprises,omitempty" bson:"autres_produits_charges_reprises,omitempty"`
	DotationAmortissement           *float64  `json:"dotation_amortissement,omitempty" bson:"dotation_amortissement,omitempty"`
	ResultatExploitation            float64   `json:"resultat_expl" bson:"resultat_expl"`
	OperationsCommun                *float64  `json:"operations_commun,omitempty" bson:"operations_commun,omitempty"`
	ProduitsFinanciers              *float64  `json:"produits_financiers,omitempty" bson:"produits_financiers,omitempty"`
	ChargesFinancieres              *float64  `json:"charges_financieres,omitempty" bson:"charges_financieres,omitempty"`
	Interets                        *float64  `json:"interets,omitempty" bson:"interets,omitempty"`
	ResultatAvantImpot              *float64  `json:"resultat_avant_impot,omitempty" bson:"resultat_avant_impot,omitempty"`
	ProduitExceptionnel             *float64  `json:"produit_exceptionnel,omitempty" bson:"produit_exceptionnel,omitempty"`
	ChargeExceptionnelle            *float64  `json:"charge_exceptionnelle,omitempty" bson:"charge_exceptionnelle,omitempty"`
	ParticipationSalaries           *float64  `json:"participation_salaries,omitempty" bson:"participation_salaries,omitempty"`
	ImpotBenefice                   *float64  `json:"impot_benefice,omitempty" bson:"impot_benefice,omitempty"`
}

// Sirene detail
type Sirene struct {
	Region          string   `json:"region" bson:"region"`
	Commune         string   `json:"commune" bson:"commune"`
	RaisonSociale   string   `json:"raison_sociale" bson:"raison_sociale"`
	TypeVoie        string   `json:"type_voie" bson:"type_voie"`
	Siren           string   `json:"siren" bson:"siren"`
	CodePostal      string   `json:"code_postal" bson:"code_postal"`
	Lattitude       float64  `json:"lattitude" bson:"lattitude"`
	Adresse         []string `json:"adresse" bson:"adresse"`
	Departement     string   `json:"departement" bson:"departement"`
	NatureJuridique string   `json:"nature_juridique" bson:"nature_juridique"`
	NumeroVoie      string   `json:"numero_voie" bson:"numero_voie"`
	Ape             string   `json:"ape" bson:"ape"`
	Longitude       float64  `json:"longitude" bson:"longitude"`
	Nic             string   `json:"nic" bson:"nic"`
	NicSiege        string   `json:"nic_siege" bson:"nic_siege"`
}

// DatapiDetection résultat de l'aggregation
type DatapiDetection struct {
	Key   map[string]string
	Scope []string
	Value map[string]interface{}
}

func computeDetection(detection Detection) (detections []daclient.Object) {
	caVal, caVar, reVal, reVar, annee := computeDiane(detection)
	dernierEffectif, variationEffectif := computeEffectif(detection)

	key := map[string]string{
		"siret": detection.ID["key"],
		"batch": detection.ID["batch"],
		"type":  "detection",
	}

	var acteurs []string
	if detection.Connu {
		acteurs = append(acteurs, "connu")
	}

	scopeB := []string{"detection", detection.Etablissement.Value.Sirene.Departement}
	valueB := map[string]interface{}{
		"acteurs":                 acteurs,
		"raison_sociale":          detection.Entreprise.Value.SireneUL.RaisonSociale,
		"activite":                detection.Etablissement.Value.Sirene.Ape,
		"urssaf":                  computeUrssaf(detection),
		"activite_partielle":      computeActivitePartielle(detection),
		"dernier_effectif":        &dernierEffectif,
		"variation_effectif":      &variationEffectif,
		"annee_ca":                annee,
		"ca":                      caVal,
		"variation_ca":            caVar,
		"resultat_expl":           reVal,
		"variation_resultat_expl": reVar,
		"departement":             detection.Etablissement.Value.Sirene.Departement,
		"etat_procol":             detection.Etablissement.Value.LastProcol.Etat,
		"date_procol":             detection.Etablissement.Value.LastProcol.Date,
	}

	scopeA := []string{"detection", "score", detection.Etablissement.Value.Sirene.Departement}
	valueA := map[string]interface{}{
		"score": detection.Score,
		"alert": detection.Alert,
		"diff":  detection.Diff,
	}

	detections = append(detections, daclient.Object{
		Key:   key,
		Scope: scopeA,
		Value: valueA,
	})

	detections = append(detections, daclient.Object{
		Key:   key,
		Scope: scopeB,
		Value: valueB,
	})

	return detections
}

func computeEtablissement(detection Detection) (objects []daclient.Object) {
	key := map[string]string{
		"siret": detection.Etablissement.Value.Sirene.Siren + detection.Etablissement.Value.Sirene.Nic,
		"siren": detection.Etablissement.Value.Sirene.Siren,
		"batch": detection.ID["batch"],
		"type":  "detail",
	}

	scope := []string{detection.Etablissement.Value.Sirene.Departement}
	sirene := detection.Etablissement.Value.Sirene
	sirene.RaisonSociale = detection.Entreprise.Value.SireneUL.RaisonSociale
	value := map[string]interface{}{
		"diane":                detection.Entreprise.Value.Diane,
		"effectif":             detection.Etablissement.Value.Effectif,
		"sirene":               sirene,
		"procedure_collective": detection.Etablissement.Value.Procol,
	}

	scopeURSSAF := []string{"urssaf", detection.Etablissement.Value.Sirene.Departement}
	valueURSSAF := map[string]interface{}{
		"debit":      detection.Etablissement.Value.Debit,
		"delai":      detection.Etablissement.Value.Delai,
		"cotisation": detection.Etablissement.Value.Cotisation,
	}

	scopeDGEFP := []string{"dgefp", detection.Etablissement.Value.Sirene.Departement}
	valueDGEFP := map[string]interface{}{
		"apconso":   detection.Etablissement.Value.APConso,
		"apdemande": detection.Etablissement.Value.APDemande,
	}

	scopeBDF := []string{"bdf", detection.Etablissement.Value.Sirene.Departement}
	valueBDF := map[string]interface{}{
		"bdf": detection.Entreprise.Value.BDF,
	}

	object := daclient.Object{
		Key:   key,
		Scope: scope,
		Value: value,
	}

	objectURSSAF := daclient.Object{
		Key:   key,
		Scope: scopeURSSAF,
		Value: valueURSSAF,
	}

	if detection.Alert != "Pas d'alerte" {
		objectDGEFP := daclient.Object{
			Key:   key,
			Scope: scopeDGEFP,
			Value: valueDGEFP,
		}
		objects = append(objects, objectDGEFP)
	}

	objectBDF := daclient.Object{
		Key:   key,
		Scope: scopeBDF,
		Value: valueBDF,
	}

	objects = append(objects, object, objectURSSAF, objectBDF)
	return objects
}

// Compute traite un objet detection pour produire les objets datapi
func Compute(detection Detection) ([]daclient.Object, error) {

	if detection.Etablissement.Value.Sirene.Departement != "" {
		var objects []daclient.Object
		if detection.Alert != "Pas d'alerte" {
			objects = append(objects, computeDetection(detection)...)
		}
		objects = append(objects, computeEtablissement(detection)...)
		return objects, nil
	}

	return nil, errors.New("pas d'information sirene, objet ignoré")
}

func computeEffectif(detection Detection) (dernierEffectif int, variationEffectif float64) {
	l := len(detection.Etablissement.Value.Effectif)
	if l > 2 {
		dernierEffectif := detection.Etablissement.Value.Effectif[l-1].Effectif
		variationEffectif := float64(detection.Etablissement.Value.Effectif[l-1].Effectif) / float64(detection.Etablissement.Value.Effectif[l-2].Effectif)
		return dernierEffectif, variationEffectif
	}
	return 0, 0
}

func computeUrssaf(detection Detection) bool {
	debits := detection.Etablissement.Value.Debit
	if len(debits) == 24 {
		for i := 24 - 3; i < 24; i++ {
			if (debits[i].PartOuvriere+debits[i].PartPatronale)/(debits[i-1].PartOuvriere+debits[i-1].PartPatronale) > 1.01 {
				return true
			}
		}
	}
	return false
}

func computeActivitePartielle(detection Detection) bool {
	batch := detection.ID["batch"]
	date, err := batchToTime(batch)
	if err != nil {
		return false
	}
	for _, v := range detection.Etablissement.Value.APConso {
		if v.Periode.Add(1 * time.Second).After(date) {
			return true
		}
	}
	for _, v := range detection.Etablissement.Value.APDemande {
		if v.Periode.End.Add(1 * time.Second).After(date) {
			return true
		}
	}

	return false
}

func computeDiane(detection Detection) (caVal *float64, caVar *float64, reVal *float64, reVar *float64, annee *float64) {
	for i := 1; i < len(detection.Entreprise.Value.Diane); i++ {
		if detection.Entreprise.Value.Diane[i-1].ChiffreAffaire != 0 &&
			detection.Entreprise.Value.Diane[i].ChiffreAffaire != 0 {
			d1 := detection.Entreprise.Value.Diane[i-1]
			d2 := detection.Entreprise.Value.Diane[i]
			annee = &detection.Entreprise.Value.Diane[i-1].Exercice
			cavar := d1.ChiffreAffaire / d2.ChiffreAffaire
			caVal = &d1.ChiffreAffaire
			caVar = &cavar

			if d2.ResultatExploitation*d1.ResultatExploitation != 0 {
				reVal = &d1.ResultatExploitation
				revar := d1.ResultatExploitation / d2.ResultatExploitation
				reVar = &revar
			}

			break
		}
	}

	return caVal, caVar, reVal, reVar, annee
}

func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
}
