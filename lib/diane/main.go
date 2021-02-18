package diane

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Diane Information financières
type Diane struct {
	Annee                           *int      `col:"Annee" json:"exercice_diane,omitempty" bson:"exercice_diane,omitempty"`
	NomEntreprise                   string    `col:"Nom de l'entreprise" json:"nom_entreprise" bson:"nom_entreprise,omitempty"`
	NumeroSiren                     string    `col:"Numéro Siren" json:"numero_siren" bson:"numero_siren,omitempty"`
	StatutJuridique                 string    `col:"Statut juridique " json:"statut_juridique" bson:"statut_juridique,omitempty"`
	ProcedureCollective             bool      `col:"Procédure collective" json:"procedure_collective" bson:"procedure_collective,omitempty"`
	EffectifConsolide               *int      `col:"Effectif consolidé" json:"effectif_consolide" bson:"effectif_consolide,omitempty"`
	DetteFiscaleEtSociale           *float64  `col:"Dettes fiscales et sociales kEUR" json:"dette_fiscale_et_sociale" bson:"dette_fiscale_et_sociale,omitempty"`
	FraisDeRetD                     *float64  `col:"Frais de R&D : net kEUR" json:"frais_de_RetD" bson:"frais_de_RetD,omitempty"`
	ConcesBrevEtDroitsSim           *float64  `col:"Conces., brev. et droits sim. : net kEUR" json:"conces_brev_et_droits_sim" bson:"conces_brev_et_droits_sim,omitempty"`
	NombreEtabSecondaire            *int      `col:"Nombre d’ES" json:"nombre_etab_secondaire" bson:"nombre_etab_secondaire,omitempty"` // Ancien nom: "Nombre d"ES"
	NombreFiliale                   *int      `col:"Nombre de filiales" json:"nombre_filiale" bson:"nombre_filiale,omitempty"`
	TailleCompoGroupe               *int      `col:"Taille de la Composition du Groupe" json:"taille_compo_groupe" bson:"taille_compo_groupe,omitempty"`
	ArreteBilan                     time.Time `col:"Date de clôture" json:"arrete_bilan_diane" bson:"arrete_bilan_diane,omitempty"`
	NombreMois                      *int      `col:"Nombre de mois" json:"nombre_mois" bson:"nombre_mois,omitempty"`
	ConcoursBancaireCourant         *float64  `col:"Conc. banc. cour. & sold. cr. kEUR" json:"concours_bancaire_courant" bson:"concours_bancaire_courant,omitempty"`
	EquilibreFinancier              *float64  `col:"Equilibre financier" json:"equilibre_financier" bson:"equilibre_financier,omitempty"`
	IndependanceFinanciere          *float64  `col:"Indépendance fin. %" json:"independance_financiere" bson:"independance_financiere,omitempty"`
	Endettement                     *float64  `col:"Endettement %" json:"endettement" bson:"endettement,omitempty"`
	AutonomieFinanciere             *float64  `col:"Autonomie fin. %" json:"autonomie_financiere" bson:"autonomie_financiere,omitempty"`
	DegreImmoCorporelle             *float64  `col:"Degré d'amort. des immob. corp. %" json:"degre_immo_corporelle" bson:"degre_immo_corporelle,omitempty"`
	FinancementActifCirculant       *float64  `col:"Financ. de l'actif circ. net" json:"financement_actif_circulant" bson:"financement_actif_circulant,omitempty"`
	LiquiditeGenerale               *float64  `col:"Liquidité générale" json:"liquidite_generale" bson:"liquidite_generale,omitempty"`
	LiquiditeReduite                *float64  `col:"Liquidité réduite" json:"liquidite_reduite" bson:"liquidite_reduite,omitempty"`
	RotationStocks                  *float64  `col:"Rotation des stocks jours" json:"rotation_stocks" bson:"rotation_stocks,omitempty"`
	CreditClient                    *float64  `col:"Crédit clients jours" json:"credit_client" bson:"credit_client,omitempty"`
	CreditFournisseur               *float64  `col:"Crédit fournisseurs jours" json:"credit_fournisseur" bson:"credit_fournisseur,omitempty"`
	CAparEffectif                   *float64  `col:"C. A. par effectif (milliers/pers.) kEUR" json:"ca_par_effectif" bson:"ca_par_effectif,omitempty"`
	TauxInteretFinancier            *float64  `col:"Taux d'intérêt financier %" json:"taux_interet_financier" bson:"taux_interet_financier,omitempty"`
	TauxInteretSurCA                *float64  `col:"Intérêts / Chiffre d'affaires %" json:"taux_interet_sur_ca" bson:"taux_interet_sur_ca,omitempty"`
	EndettementGlobal               *float64  `col:"Endettement global jours" json:"endettement_global" bson:"endettement_global,omitempty"`
	TauxEndettement                 *float64  `col:"Taux d'endettement %" json:"taux_endettement" bson:"taux_endettement,omitempty"`
	CapaciteRemboursement           *float64  `col:"Capacité de remboursement" json:"capacite_remboursement" bson:"capacite_remboursement,omitempty"`
	CapaciteAutofinancement         *float64  `col:"Capacité d'autofin. %" json:"capacite_autofinancement" bson:"capacite_autofinancement,omitempty"`
	CouvertureCaFdr                 *float64  `col:"Couv. du C.A. par le f.d.r. jours" json:"couverture_ca_fdr" bson:"couverture_ca_fdr,omitempty"`
	CouvertureCaBesoinFdr           *float64  `col:"Couv. du C.A. par bes. en fdr jours" json:"couverture_ca_besoin_fdr" bson:"couverture_ca_besoin_fdr,omitempty"`
	PoidsBFRExploitation            *float64  `col:"Poids des BFR d'exploitation %" json:"poids_bfr_exploitation" bson:"poids_bfr_exploitation,omitempty"`
	Exportation                     *float64  `col:"Exportation %" json:"exportation" bson:"exportation,omitempty"`
	EfficaciteEconomique            *float64  `col:"Efficacité économique (milliers/pers.) kEUR" json:"efficacite_economique" bson:"efficacite_economique,omitempty"`
	ProductivitePotentielProduction *float64  `col:"Prod. du potentiel de production" json:"productivite_potentiel_production" bson:"productivite_potentiel_production,omitempty"`
	ProductiviteCapitalFinancier    *float64  `col:"Productivité du capital financier" json:"productivite_capital_financier" bson:"productivite_capital_financier,omitempty"`
	ProductiviteCapitalInvesti      *float64  `col:"Productivité du capital investi" json:"productivite_capital_investi" bson:"productivite_capital_investi,omitempty"`
	TauxDInvestissementProductif    *float64  `col:"Taux d'invest. productif %" json:"taux_d_investissement_productif" bson:"taux_d_investissement_productif,omitempty"`
	RentabiliteEconomique           *float64  `col:"Rentabilité économique %" json:"rentabilite_economique" bson:"rentabilite_economique,omitempty"`
	Performance                     *float64  `col:"Performance %" json:"performance" bson:"performance,omitempty"`
	RendementBrutFondsPropres       *float64  `col:"Rend. brut des f. propres nets %" json:"rendement_brut_fonds_propres" bson:"rendement_brut_fonds_propres,omitempty"`
	RentabiliteNette                *float64  `col:"Rentabilité nette %" json:"rentabilite_nette" bson:"rentabilite_nette,omitempty"`
	RendementCapitauxPropres        *float64  `col:"Rend. des capitaux propres nets %" json:"rendement_capitaux_propres" bson:"rendement_capitaux_propres,omitempty"`
	RendementRessourcesDurables     *float64  `col:"Rend. des res. durables nettes %" json:"rendement_ressources_durables" bson:"rendement_ressources_durables,omitempty"`
	TauxMargeCommerciale            *float64  `col:"Taux de marge commerciale %" json:"taux_marge_commerciale" bson:"taux_marge_commerciale,omitempty"`
	TauxValeurAjoutee               *float64  `col:"Taux de valeur ajoutée %" json:"taux_valeur_ajoutee" bson:"taux_valeur_ajoutee,omitempty"`
	PartSalaries                    *float64  `col:"Part des salariés %" json:"part_salaries" bson:"part_salaries,omitempty"`
	PartEtat                        *float64  `col:"Part de l'Etat %" json:"part_etat" bson:"part_etat,omitempty"`
	PartPreteur                     *float64  `col:"Part des prêteurs %" json:"part_preteur" bson:"part_preteur,omitempty"`
	PartAutofinancement             *float64  `col:"Part de l'autofin. %" json:"part_autofinancement" bson:"part_autofinancement,omitempty"`
	CA                              *float64  `col:"Chiffre d'affaires net (H.T.) kEUR" json:"ca" bson:"ca,omitempty"`
	CAExportation                   *float64  `col:"Dont exportation kEUR" json:"ca_exportation" bson:"ca_exportation,omitempty"`
	AchatMarchandises               *float64  `col:"Achats march. et autres approv. kEUR" json:"achat_marchandises" bson:"achat_marchandises,omitempty"`
	AchatMatieresPremieres          *float64  `col:"Achats de mat. prem. et autres approv. kEUR" json:"achat_matieres_premieres" bson:"achat_matieres_premieres,omitempty"`
	Production                      *float64  `col:"Production de l'ex. kEUR" json:"production" bson:"production,omitempty"`
	MargeCommerciale                *float64  `col:"Marge commerciale kEUR" json:"marge_commerciale" bson:"marge_commerciale,omitempty"`
	Consommation                    *float64  `col:"Consommation de l'ex. kEUR" json:"consommation" bson:"consommation,omitempty"`
	AutresAchatsChargesExternes     *float64  `col:"Autres achats et charges externes kEUR" json:"autres_achats_charges_externes" bson:"autres_achats_charges_externes,omitempty"`
	ValeurAjoutee                   *float64  `col:"Valeur ajoutée kEUR" json:"valeur_ajoutee" bson:"valeur_ajoutee,omitempty"`
	ChargePersonnel                 *float64  `col:"Charges de personnel kEUR" json:"charge_personnel" bson:"charge_personnel,omitempty"`
	ImpotsTaxes                     *float64  `col:"Impôts, taxes et vers. assimil. kEUR" json:"impots_taxes" bson:"impots_taxes,omitempty"`
	SubventionsDExploitation        *float64  `col:"Subventions d'expl. kEUR" json:"subventions_d_exploitation" bson:"subventions_d_exploitation,omitempty"`
	ExcedentBrutDExploitation       *float64  `col:"Excédent brut d'exploitation kEUR" json:"excedent_brut_d_exploitation" bson:"excedent_brut_d_exploitation,omitempty"`
	AutresProduitsChargesReprises   *float64  `col:"Autres Prod., char. et Repr. kEUR" json:"autres_produits_charges_reprises" bson:"autres_produits_charges_reprises,omitempty"`
	DotationAmortissement           *float64  `col:"Dot. d'exploit. aux amort. et prov. kEUR" json:"dotation_amortissement" bson:"dotation_amortissement,omitempty"`
	ResultatExpl                    *float64  `col:"Résultat d'expl. kEUR" json:"resultat_expl" bson:"resultat_expl,omitempty"`
	OperationsCommun                *float64  `col:"Opérations en commun kEUR" json:"operations_commun" bson:"operations_commun,omitempty"`
	ProduitsFinanciers              *float64  `col:"Produits fin. kEUR" json:"produits_financiers" bson:"produits_financiers,omitempty"`
	ChargesFinancieres              *float64  `col:"Charges fin. kEUR" json:"charges_financieres" bson:"charges_financieres,omitempty"`
	Interets                        *float64  `col:"Intérêts et charges assimilées kEUR" json:"interets" bson:"interets,omitempty"`
	ResultatAvantImpot              *float64  `col:"Résultat courant avant impôts kEUR" json:"resultat_avant_impot" bson:"resultat_avant_impot,omitempty"`
	ProduitExceptionnel             *float64  `col:"Produits except. kEUR" json:"produit_exceptionnel" bson:"produit_exceptionnel,omitempty"`
	ChargeExceptionnelle            *float64  `col:"Charges except. kEUR" json:"charge_exceptionnelle" bson:"charge_exceptionnelle,omitempty"`
	ParticipationSalaries           *float64  `col:"Particip. des sal. aux résul. kEUR" json:"participation_salaries" bson:"participation_salaries,omitempty"`
	ImpotBenefice                   *float64  `col:"Impôts sur le bénéf. et impôts diff. kEUR" json:"impot_benefice" bson:"impot_benefice,omitempty"`
	BeneficeOuPerte                 *float64  `col:"Bénéfice ou perte kEUR" json:"benefice_ou_perte" bson:"benefice_ou_perte,omitempty"`
	// TODO: ajouter NotePreface ou le retirer de la documentation (cf https://github.com/signaux-faibles/documentation/blob/master/description-donnees.md#donn%C3%A9es-financi%C3%A8res-issues-des-bilans-d%C3%A9pos%C3%A9s-au-greffe-de-tribunaux-de-commerce)

	// Colonnes non utilisées:
	// 01 "Marquée"; (propriété `marquee`)
	// 13 "Dernière année disponible";
	// 59 "Achats de march. kEUR";
}

// Key id de l'objet
func (diane Diane) Key() string {
	return diane.NumeroSiren
}

// Type de données
func (diane Diane) Type() string {
	return "diane"
}

// Scope de l'objet
func (diane Diane) Scope() string {
	return "entreprise"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &dianeParser{}

type dianeParser struct {
	closeFct func() error
	reader   *csv.Reader
	idx      marshal.ColMapping
}

func (parser *dianeParser) GetFileType() string {
	return "diane"
}

func (parser *dianeParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *dianeParser) Open(filePath string) (err error) {
	var reader *io.ReadCloser
	parser.closeFct, reader, err = preprocessDianeFile(filePath)
	if err != nil {
		return err
	}
	return parser.initCsvReader(*reader)
}

func (parser *dianeParser) Close() error {
	return parser.closeFct()
}

func preprocessDianeFile(filePath string) (func() error, *io.ReadCloser, error) {

	// TODO: implémenter ces traitements en Go
	pipedCmds := []*exec.Cmd{
		exec.Command("cat", filePath),
		exec.Command("iconv", "--from-code", "UTF-16LE", "--to-code", "UTF-8"), // conversion d'encodage de fichier car awk ne supporte pas UTF-16LE
		exec.Command("sed", "s/\r$//"),                                         // forcer l'usage de retours charriot au format UNIX car le caractère \r cause une duplication de colonne depuis le script awk
		exec.Command("sed", "s/  / /g"),                                        // dé-dupliquer les caractères d'espacement, notamment dans les en-têtes de colonnes
		exec.Command("awk", awkScript),                                         // réorganisation des des données pour éviter la duplication de colonnes par année
	}
	lastCmd := pipedCmds[len(pipedCmds)-1]

	var err error
	close := func() error {
		if err := lastCmd.Wait(); err != nil {
			return errors.New("[convert_diane.sh] failed with " + err.Error())
		}
		return nil
	}

	// pipe streams between commands
	for i := 0; i < len(pipedCmds)-1; i++ {
		pipedCmds[i].Stderr = os.Stderr
		pipedCmds[i+1].Stdin, err = pipedCmds[i].StdoutPipe()
		if err != nil {
			return close, nil, err
		}
	}
	lastCmd.Stderr = os.Stderr
	stdout, err := lastCmd.StdoutPipe()
	if err != nil {
		return close, nil, err
	}

	// start piped commands, from last to first
	for i := len(pipedCmds) - 1; i >= 0; i-- {
		err = pipedCmds[i].Start()
		if err != nil {
			return close, nil, err
		}
	}

	return close, &stdout, nil
}

func (parser *dianeParser) initCsvReader(reader io.Reader) (err error) {
	parser.reader = csv.NewReader(reader)
	parser.reader.Comma = ';'
	parser.reader.LazyQuotes = true
	parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Diane{})
	return err
}

func (parser *dianeParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else if len(row) < 83 {
			parsedLine.AddRegularError(errors.New("Ligne invalide"))
		} else {
			parsedLine.AddTuple(parseDianeRow(parser.idx, row))
		}
		parsedLineChan <- parsedLine
	}
}

// parseDianeRow construit un objet Diane à partir d'une ligne de valeurs récupérée depuis un fichier
func parseDianeRow(idx marshal.ColMapping, row []string) (diane Diane) {

	idxRow := idx.IndexRow(row)

	if val, err := idxRow.GetInt("Annee"); err == nil {
		diane.Annee = val
	}
	diane.NomEntreprise = idxRow.GetVal("Nom de l'entreprise")
	diane.NumeroSiren = idxRow.GetVal("Numéro Siren")
	diane.StatutJuridique = idxRow.GetVal("Statut juridique ")
	diane.ProcedureCollective = (idxRow.GetVal("Procédure collective") == "Oui")

	if val, err := idxRow.GetInt("Effectif consolidé"); err == nil {
		diane.EffectifConsolide = val
	}
	if val, err := idxRow.GetCommaFloat64("Dettes fiscales et sociales kEUR"); err == nil {
		diane.DetteFiscaleEtSociale = val
	}
	if val, err := idxRow.GetCommaFloat64("Frais de R&D : net kEUR"); err == nil {
		diane.FraisDeRetD = val
	}
	if val, err := idxRow.GetCommaFloat64("Conces., brev. et droits sim. : net kEUR"); err == nil {
		diane.ConcesBrevEtDroitsSim = val
	}
	if val, err := idxRow.GetInt("Nombre d’ES"); err == nil {
		diane.NombreEtabSecondaire = val
	}
	if val, err := idxRow.GetInt("Nombre de filiales"); err == nil {
		diane.NombreFiliale = val
	}
	if val, err := idxRow.GetInt("Taille de la Composition du Groupe"); err == nil {
		diane.TailleCompoGroupe = val
	}
	if i, err := time.Parse("02/01/2006", idxRow.GetVal("Date de clôture")); err == nil {
		diane.ArreteBilan = i
	}
	if val, err := idxRow.GetInt("Nombre de mois"); err == nil {
		diane.NombreMois = val
	}
	if val, err := idxRow.GetCommaFloat64("Conc. banc. cour. & sold. cr. kEUR"); err == nil {
		diane.ConcoursBancaireCourant = val
	}
	if val, err := idxRow.GetCommaFloat64("Equilibre financier"); err == nil {
		diane.EquilibreFinancier = val
	}
	if val, err := idxRow.GetCommaFloat64("Indépendance fin. %"); err == nil {
		diane.IndependanceFinanciere = val
	}
	if val, err := idxRow.GetCommaFloat64("Endettement %"); err == nil {
		diane.Endettement = val
	}
	if val, err := idxRow.GetCommaFloat64("Autonomie fin. %"); err == nil {
		diane.AutonomieFinanciere = val
	}
	if val, err := idxRow.GetCommaFloat64("Degré d'amort. des immob. corp. %"); err == nil {
		diane.DegreImmoCorporelle = val
	}
	if val, err := idxRow.GetCommaFloat64("Financ. de l'actif circ. net"); err == nil {
		diane.FinancementActifCirculant = val
	}
	if val, err := idxRow.GetCommaFloat64("Liquidité générale"); err == nil {
		diane.LiquiditeGenerale = val
	}
	if val, err := idxRow.GetCommaFloat64("Liquidité réduite"); err == nil {
		diane.LiquiditeReduite = val
	}
	if val, err := idxRow.GetCommaFloat64("Rotation des stocks jours"); err == nil {
		diane.RotationStocks = val
	}
	if val, err := idxRow.GetCommaFloat64("Crédit clients jours"); err == nil {
		diane.CreditClient = val
	}
	if val, err := idxRow.GetCommaFloat64("Crédit fournisseurs jours"); err == nil {
		diane.CreditFournisseur = val
	}
	if val, err := idxRow.GetCommaFloat64("C. A. par effectif (milliers/pers.) kEUR"); err == nil {
		diane.CAparEffectif = val
	}
	if val, err := idxRow.GetCommaFloat64("Taux d'intérêt financier %"); err == nil {
		diane.TauxInteretFinancier = val
	}
	if val, err := idxRow.GetCommaFloat64("Intérêts / Chiffre d'affaires %"); err == nil {
		diane.TauxInteretSurCA = val
	}
	if val, err := idxRow.GetCommaFloat64("Endettement global jours"); err == nil {
		diane.EndettementGlobal = val
	}
	if val, err := idxRow.GetCommaFloat64("Taux d'endettement %"); err == nil {
		diane.TauxEndettement = val
	}
	if val, err := idxRow.GetCommaFloat64("Capacité de remboursement"); err == nil {
		diane.CapaciteRemboursement = val
	}
	if val, err := idxRow.GetCommaFloat64("Capacité d'autofin. %"); err == nil {
		diane.CapaciteAutofinancement = val
	}
	if val, err := idxRow.GetCommaFloat64("Couv. du C.A. par le f.d.r. jours"); err == nil {
		diane.CouvertureCaFdr = val
	}
	if val, err := idxRow.GetCommaFloat64("Couv. du C.A. par bes. en fdr jours"); err == nil {
		diane.CouvertureCaBesoinFdr = val
	}
	if val, err := idxRow.GetCommaFloat64("Poids des BFR d'exploitation %"); err == nil {
		diane.PoidsBFRExploitation = val
	}
	if val, err := idxRow.GetCommaFloat64("Exportation %"); err == nil {
		diane.Exportation = val
	}
	if val, err := idxRow.GetCommaFloat64("Efficacité économique (milliers/pers.) kEUR"); err == nil {
		diane.EfficaciteEconomique = val
	}
	if val, err := idxRow.GetCommaFloat64("Prod. du potentiel de production"); err == nil {
		diane.ProductivitePotentielProduction = val
	}
	if val, err := idxRow.GetCommaFloat64("Productivité du capital financier"); err == nil {
		diane.ProductiviteCapitalFinancier = val
	}
	if val, err := idxRow.GetCommaFloat64("Productivité du capital investi"); err == nil {
		diane.ProductiviteCapitalInvesti = val
	}
	if val, err := idxRow.GetCommaFloat64("Taux d'invest. productif %"); err == nil {
		diane.TauxDInvestissementProductif = val
	}
	if val, err := idxRow.GetCommaFloat64("Rentabilité économique %"); err == nil {
		diane.RentabiliteEconomique = val
	}
	if val, err := idxRow.GetCommaFloat64("Performance %"); err == nil {
		diane.Performance = val
	}
	if val, err := idxRow.GetCommaFloat64("Rend. brut des f. propres nets %"); err == nil {
		diane.RendementBrutFondsPropres = val
	}
	if val, err := idxRow.GetCommaFloat64("Rentabilité nette %"); err == nil {
		diane.RentabiliteNette = val
	}
	if val, err := idxRow.GetCommaFloat64("Rend. des capitaux propres nets %"); err == nil {
		diane.RendementCapitauxPropres = val
	}
	if val, err := idxRow.GetCommaFloat64("Rend. des res. durables nettes %"); err == nil {
		diane.RendementRessourcesDurables = val
	}
	if val, err := idxRow.GetCommaFloat64("Taux de marge commerciale %"); err == nil {
		diane.TauxMargeCommerciale = val
	}
	if val, err := idxRow.GetCommaFloat64("Taux de valeur ajoutée %"); err == nil {
		diane.TauxValeurAjoutee = val
	}
	if val, err := idxRow.GetCommaFloat64("Part des salariés %"); err == nil {
		diane.PartSalaries = val
	}
	if val, err := idxRow.GetCommaFloat64("Part de l'Etat %"); err == nil {
		diane.PartEtat = val
	}
	if val, err := idxRow.GetCommaFloat64("Part des prêteurs %"); err == nil {
		diane.PartPreteur = val
	}
	if val, err := idxRow.GetCommaFloat64("Part de l'autofin. %"); err == nil {
		diane.PartAutofinancement = val
	}
	if val, err := idxRow.GetCommaFloat64("Chiffre d'affaires net (H.T.) kEUR"); err == nil {
		diane.CA = val
	}
	if val, err := idxRow.GetCommaFloat64("Dont exportation kEUR"); err == nil {
		diane.CAExportation = val
	}
	if val, err := idxRow.GetCommaFloat64("Achats march. et autres approv. kEUR"); err == nil {
		diane.AchatMarchandises = val
	}
	if val, err := idxRow.GetCommaFloat64("Achats de mat. prem. et autres approv. kEUR"); err == nil {
		diane.AchatMatieresPremieres = val
	}
	if val, err := idxRow.GetCommaFloat64("Production de l'ex. kEUR"); err == nil {
		diane.Production = val
	}
	if val, err := idxRow.GetCommaFloat64("Marge commerciale kEUR"); err == nil {
		diane.MargeCommerciale = val
	}
	if val, err := idxRow.GetCommaFloat64("Consommation de l'ex. kEUR"); err == nil {
		diane.Consommation = val
	}
	if val, err := idxRow.GetCommaFloat64("Autres achats et charges externes kEUR"); err == nil {
		diane.AutresAchatsChargesExternes = val
	}
	if val, err := idxRow.GetCommaFloat64("Valeur ajoutée kEUR"); err == nil {
		diane.ValeurAjoutee = val
	}
	if val, err := idxRow.GetCommaFloat64("Charges de personnel kEUR"); err == nil {
		diane.ChargePersonnel = val
	}
	if val, err := idxRow.GetCommaFloat64("Impôts, taxes et vers. assimil. kEUR"); err == nil {
		diane.ImpotsTaxes = val
	}
	if val, err := idxRow.GetCommaFloat64("Subventions d'expl. kEUR"); err == nil {
		diane.SubventionsDExploitation = val
	}
	if val, err := idxRow.GetCommaFloat64("Excédent brut d'exploitation kEUR"); err == nil {
		diane.ExcedentBrutDExploitation = val
	}
	if val, err := idxRow.GetCommaFloat64("Autres Prod., char. et Repr. kEUR"); err == nil {
		diane.AutresProduitsChargesReprises = val
	}
	if val, err := idxRow.GetCommaFloat64("Dot. d'exploit. aux amort. et prov. kEUR"); err == nil {
		diane.DotationAmortissement = val
	}
	if val, err := idxRow.GetCommaFloat64("Résultat d'expl. kEUR"); err == nil {
		diane.ResultatExpl = val
	}
	if val, err := idxRow.GetCommaFloat64("Opérations en commun kEUR"); err == nil {
		diane.OperationsCommun = val
	}
	if val, err := idxRow.GetCommaFloat64("Produits fin. kEUR"); err == nil {
		diane.ProduitsFinanciers = val
	}
	if val, err := idxRow.GetCommaFloat64("Charges fin. kEUR"); err == nil {
		diane.ChargesFinancieres = val
	}
	if val, err := idxRow.GetCommaFloat64("Intérêts et charges assimilées kEUR"); err == nil {
		diane.Interets = val
	}
	if val, err := idxRow.GetCommaFloat64("Résultat courant avant impôts kEUR"); err == nil {
		diane.ResultatAvantImpot = val
	}
	if val, err := idxRow.GetCommaFloat64("Produits except. kEUR"); err == nil {
		diane.ProduitExceptionnel = val
	}
	if val, err := idxRow.GetCommaFloat64("Charges except. kEUR"); err == nil {
		diane.ChargeExceptionnelle = val
	}
	if val, err := idxRow.GetCommaFloat64("Particip. des sal. aux résul. kEUR"); err == nil {
		diane.ParticipationSalaries = val
	}
	if val, err := idxRow.GetCommaFloat64("Impôts sur le bénéf. et impôts diff. kEUR"); err == nil {
		diane.ImpotBenefice = val
	}
	if val, err := idxRow.GetCommaFloat64("Bénéfice ou perte kEUR"); err == nil {
		diane.BeneficeOuPerte = val
	}
	return diane
}

// This awk spreads company data so that each year of data has its own row.
const awkScript = `
BEGIN { # Semi-column separated csv as input and output
  FS = ";"
  OFS = ";"
  RE_YEAR = "[[:digit:]][[:digit:]][[:digit:]][[:digit:]]"
  RE_YEAR_SUFFIX = / ([[:digit:]][[:digit:]][[:digit:]][[:digit:]])$/
  first_year = last_year = 0
}
FNR==1 { # Heading row => coalesce yearly fields
  printf "%s", "\"Annee\""
  for (field = 1; field <= NF; ++field) {
    if ($field !~ RE_YEAR_SUFFIX) { # Field without year
      fields[++nb_fields] = field
      printf "%s%s",  OFS, $field
    } else { # Field with year
      match($field, RE_YEAR, year)
      field_name = gensub(" "year[0], "", "g", $field) # Remove year from column name
      first_year = !first_year || year[0] < first_year ? year[0] : first_year
      last_year = !last_year || year[0] > last_year ? year[0] : last_year
      if (!yearly_fields[field_name]) {
        ++nb_fields
        ++yearly_fields[field_name]
        printf "%s%s", OFS, field_name;
      }
      fields[nb_fields, year[0]] = field
    }
  }
  printf "%s", ORS
}
FNR>1 && $1 !~ "Marquée" { # Data row
  for (current_year = first_year; current_year <= last_year; ++current_year) {
    printf "%i", current_year
    for (field = 1; field <= nb_fields; ++field) {
      if (fields[field, current_year] && $(fields[field, current_year])) {
        printf "%s%s", OFS, $(fields[field, current_year]);
      } else if (fields[field] && $(fields[field])) {
        printf "%s%s", OFS, $(fields[field]);
      } else {
        printf "%s%s", OFS, "\"\"";
      }
    }
    printf "%s", ORS # Each year on a new line
  }
}`
