package exportdatapi

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Detection correspond aux données retournées pour l'export Datapi
type Detection struct {
	ID            map[string]string `json:"_id" bson:"_id"`
	Score         float64           `json:"score" bson:"score"`
	Diff          float64           `json:"diff" bson:"diff"`
	Timestamp     time.Time         `json:"timestamp" bson:"timestamp"`
	Alert         string            `json:"alert" bson:"alert"`
	Connu         bool              `json:"connu" bson:"connu"`
	Periode       time.Time         `json:"periode" bson:"periode"`
	Etablissement *Etablissement    `json:"etablissement,omitempty" bson:"etablissement,omitempty"`
	Entreprise    *Entreprise       `json:"entreprise,omitempty" bson:"entreprise,omitempty"`
	Algo          string            `json:"algo" bson:"algo"`
}

// Score objet de la base Scores
type Score struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	Score     float64       `json:"score" bson:"score"`
	Diff      float64       `json:"diff" bson:"diff"`
	Timestamp time.Time     `json:"-" bson:"timestamp"`
	Alert     string        `json:"alert" bson:"alert"`
	Periode   string        `json:"periode" bson:"periode"`
	Batch     string        `json:"batch" bson:"batch"`
	Algo      string        `json:"algo" bson:"algo"`
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
	ID    string `bson:"_id"`
	Value struct {
		Key        string      `json:"key" bson:"key"`
		Sirene     Sirene      `json:"sirene" bson:"sirene"`
		Cotisation []float64   `json:"cotisation" bson:"cotisation"`
		Debit      []Debit     `json:"debit" bson:"debit"`
		APDemande  []APDemande `json:"apdemande" bson:"apdemande"`
		APConso    []APConso   `json:"apconso" bson:"apconso"`
		Compte     struct {
			Siret   string    `json:"siret" bson:"siret"`
			Numero  string    `json:"numero_compte" bson:"numero_compte"`
			Periode time.Time `json:"periode" bson:"periode"`
		} `json:"compte" bson:"compte"`
		Effectif        []Effectif    `json:"effectif" bson:"effectif"`
		DernierEffectif Effectif      `json:"dernier_effectif" bson:"dernier_effectif"`
		Delai           []interface{} `json:"delai" bson:"delai"`
		Procol          []Procol      `json:"procol" bson:"procol"`
		LastProcol      Procol        `json:"last_procol" bson:"last_procol"`
	} `bson:"value"`
	Scores     []Score     `json:"scores" bson:"scores"`
	Entreprise *Entreprise `json:"entreprise" bson:"entreprise"`
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
