package diane

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Diane Information financières
type Diane struct {
	Annee *int `col:"Annee" json:"exercice_diane,omitempty" bson:"exercice_diane,omitempty"`
	// Colonne sautée: "Marquée" (indice: 1)
	NomEntreprise         string   `col:"Nom de l'entreprise" json:"nom_entreprise" bson:"nom_entreprise,omitempty"`
	NumeroSiren           string   `col:"Numéro Siren" json:"numero_siren" bson:"numero_siren,omitempty"`
	StatutJuridique       string   `col:"Statut juridique " json:"statut_juridique" bson:"statut_juridique,omitempty"`
	ProcedureCollective   bool     `col:"Procédure collective" json:"procedure_collective" bson:"procedure_collective,omitempty"`
	EffectifConsolide     *int     `col:"Effectif consolidé" json:"effectif_consolide" bson:"effectif_consolide,omitempty"`
	DetteFiscaleEtSociale *float64 `col:"Dettes fiscales et sociales kEUR" json:"dette_fiscale_et_sociale" bson:"dette_fiscale_et_sociale,omitempty"`
	FraisDeRetD           *float64 `col:"Frais de R&D : net kEUR" json:"frais_de_RetD" bson:"frais_de_RetD,omitempty"`
	ConcesBrevEtDroitsSim *float64 `col:"Conces.. brev. et droits sim. : net kEUR" json:"conces_brev_et_droits_sim" bson:"conces_brev_et_droits_sim,omitempty"`
	NombreEtabSecondaire  *int     `col:"Nombre d\"ES" json:"nombre_etab_secondaire" bson:"nombre_etab_secondaire,omitempty"`
	NombreFiliale         *int     `col:"Nombre de filiales" json:"nombre_filiale" bson:"nombre_filiale,omitempty"`
	TailleCompoGroupe     *int     `col:"Taille de la Composition du Groupe" json:"taille_compo_groupe" bson:"taille_compo_groupe,omitempty"`
	// Colonne sautée: "Dernière année disponible" (indice: 13)
	ArreteBilan time.Time `col:"Date de clôture" json:"arrete_bilan_diane" bson:"arrete_bilan_diane,omitempty"`
	NombreMois  *int      `col:"Nombre de mois" json:"nombre_mois" bson:"nombre_mois,omitempty"`
	// TODO: retirer colonne redondante: "Conc.  banc. cour. & sold. cr. kEUR"
	ConcoursBancaireCourant *float64 `col:"Conc. banc. cour. & sold. cr. kEUR" json:"concours_bancaire_courant" bson:"concours_bancaire_courant,omitempty"`
	// TODO: retirer colonne redondante: "Conc. banc. cour. & sold. cr.  kEUR"
	EquilibreFinancier     *float64 `col:"Equilibre financier" json:"equilibre_financier" bson:"equilibre_financier,omitempty"`
	IndependanceFinanciere *float64 `col:"Indépendance fin. %" json:"independance_financiere" bson:"independance_financiere,omitempty"`
	// TODO: retirer colonne redondante: "Indépendance fin.  %"
	Endettement         *float64 `col:"Endettement %" json:"endettement" bson:"endettement,omitempty"`
	AutonomieFinanciere *float64 `col:"Autonomie fin. %" json:"autonomie_financiere" bson:"autonomie_financiere,omitempty"`
	DegreImmoCorporelle *float64 `col:"Degré d'amort. des immob. corp. %" json:"degre_immo_corporelle" bson:"degre_immo_corporelle,omitempty"`
	// TODO: retirer colonne redondante: "Degré d'amort.  des immob. corp. %"
	FinancementActifCirculant *float64 `col:"Financ. de l'actif circ. net" json:"financement_actif_circulant" bson:"financement_actif_circulant,omitempty"`
	LiquiditeGenerale         *float64 `col:"Liquidité générale" json:"liquidite_generale" bson:"liquidite_generale,omitempty"`
	LiquiditeReduite          *float64 `col:"Liquidité réduite" json:"liquidite_reduite" bson:"liquidite_reduite,omitempty"`
	RotationStocks            *float64 `col:"Rotation des stocks jours" json:"rotation_stocks" bson:"rotation_stocks,omitempty"`
	CreditClient              *float64 `col:"Crédit clients jours" json:"credit_client" bson:"credit_client,omitempty"`
	CreditFournisseur         *float64 `col:"Crédit fournisseurs jours" json:"credit_fournisseur" bson:"credit_fournisseur,omitempty"`
	// TODO: retirer colonne redondante: "C. A.  par effectif (milliers/pers.) kEUR"
	CAparEffectif           *float64 `col:"C. A. par effectif (milliers/pers.) kEUR" json:"ca_par_effectif" bson:"ca_apar_effectif,omitempty"`
	TauxInteretFinancier    *float64 `col:"Taux d'intérêt financier %" json:"taux_interet_financier" bson:"taux_interet_financier,omitempty"`
	TauxInteretSurCA        *float64 `col:"Intérêts / Chiffre d'affaires %" json:"taux_interet_sur_ca" bson:"taux_interet_sur_ca,omitempty"`
	EndettementGlobal       *float64 `col:"Endettement global jours" json:"endettement_global" bson:"endettement_global,omitempty"`
	TauxEndettement         *float64 `col:"Taux d'endettement %" json:"taux_endettement" bson:"taux_endettement,omitempty"`
	CapaciteRemboursement   *float64 `col:"Capacité de remboursement" json:"capacite_remboursement" bson:"capacite_remboursement,omitempty"`
	CapaciteAutofinancement *float64 `col:"Capacité d'autofin. %" json:"capacite_autofinancement" bson:"capacite_autofinancement,omitempty"`
	CouvertureCaFdr         *float64 `col:"Couv. du C.A. par le f.d.r. jours" json:"couverture_ca_fdr" bson:"couverture_ca_fdr,omitempty"`
	// TODO: retirer colonne redondante: "Couv. du C.A. par le f.d.r.  jours"
	CouvertureCaBesoinFdr           *float64 `col:"Couv. du C.A. par bes. en fdr jours" json:"couverture_ca_besoin_fdr" bson:"couverture_ca_besoin_fdr,omitempty"`
	PoidsBFRExploitation            *float64 `col:"Poids des BFR d'exploitation %" json:"poids_bfr_exploitation" bson:"poids_bfr_exploitation,omitempty"`
	Exportation                     *float64 `col:"Exportation %" json:"exportation" bson:"exportation,omitempty"`
	EfficaciteEconomique            *float64 `col:"Efficacité économique (milliers/pers.) kEUR" json:"efficacite_economique" bson:"efficacite_economique,omitempty"`
	ProductivitePotentielProduction *float64 `col:"Prod. du potentiel de production" json:"productivite_potentiel_production" bson:"productivite_potentiel_production,omitempty"`
	// TODO: retirer colonne redondante: "Prod.  du potentiel de production"
	ProductiviteCapitalFinancier *float64 `col:"Productivité du capital financier" json:"productivite_capital_financier" bson:"productivite_capital_financier,omitempty"`
	ProductiviteCapitalInvesti   *float64 `col:"Productivité du capital investi" json:"productivite_capital_investi" bson:"productivite_capital_investi,omitempty"`
	TauxDInvestissementProductif *float64 `col:"Taux d'invest. productif %" json:"taux_d_investissement_productif" bson:"taux_d_investissement_productif,omitempty"`
	RentabiliteEconomique        *float64 `col:"Rentabilité économique %" json:"rentabilite_economique" bson:"rentabilite_economique,omitempty"`
	Performance                  *float64 `col:"Performance %" json:"performance" bson:"performance,omitempty"`
	RendementBrutFondsPropres    *float64 `col:"Rend. brut des f. propres nets %" json:"rendement_brut_fonds_propres" bson:"rendement_brut_fonds_propres,omitempty"`
	// TODO: retirer colonne redondante: "Rend.  brut des f. propres nets %"
	RentabiliteNette         *float64 `col:"Rentabilité nette %" json:"rentabilite_nette" bson:"rentabilite_nette,omitempty"`
	RendementCapitauxPropres *float64 `col:"Rend. des capitaux propres nets %" json:"rendement_capitaux_propres" bson:"rendement_capitaux_propres,omitempty"`
	// TODO: retirer colonne redondante: Rend.  des capitaux propres nets %
	RendementRessourcesDurables *float64 `col:"Rend. des res. durables nettes %" json:"rendement_ressources_durables" bson:"rendement_ressources_durables,omitempty"`
	// TODO: retirer colonne redondante: "Rend. des res.  durables nettes %"
	TauxMargeCommerciale *float64 `col:"Taux de marge commerciale %" json:"taux_marge_commerciale" bson:"taux_marge_commerciale,omitempty"`
	TauxValeurAjoutee    *float64 `col:"Taux de valeur ajoutée %" json:"taux_valeur_ajoutee" bson:"taux_valeur_ajoutee,omitempty"`
	PartSalaries         *float64 `col:"Part des salariés %" json:"part_salaries" bson:"part_salaries,omitempty"`
	PartEtat             *float64 `col:"Part de l'Etat %" json:"part_etat" bson:"part_etat,omitempty"`
	PartPreteur          *float64 `col:"Part des prêteurs %" json:"part_preteur" bson:"part_preteur,omitempty"`
	PartAutofinancement  *float64 `col:"Part de l'autofin. %" json:"part_autofinancement" bson:"part_autofinancement,omitempty"`
	CA                   *float64 `col:"Chiffre d'affaires net (H.T.) kEUR" json:"ca" bson:"ca,omitempty"`
	CAExportation        *float64 `col:"Dont exportation kEUR" json:"ca_exportation" bson:"ca_exportation,omitempty"`
	AchatMarchandises    *float64 `col:"Achats march. et autres approv. kEUR" json:"achat_marchandises" bson:"achat_marchandises,omitempty"`
	// TODO: retirer colonne redondante: "Achats march.  et autres approv. kEUR"
	// Colonne sautée: "Achats de march. kEUR" (indice: 59)
	AchatMatieresPremieres *float64 `col:"Achats de mat. prem. et autres approv. kEUR" json:"achat_matieres_premieres" bson:"achat_matieres_premieres,omitempty"`
	// TODO: retirer colonne redondante: "Achats de mat. prem.  et autres approv. kEUR"
	Production                  *float64 `col:"Production de l'ex. kEUR" json:"production" bson:"production,omitempty"`
	MargeCommerciale            *float64 `col:"Marge commerciale kEUR" json:"marge_commerciale" bson:"marge_commerciale,omitempty"`
	Consommation                *float64 `col:"Consommation de l'ex. kEUR" json:"consommation" bson:"consommation,omitempty"`
	AutresAchatsChargesExternes *float64 `col:"Autres achats et charges externes kEUR" json:"autres_achats_charges_externes" bson:"autres_achats_charges_externes,omitempty"`
	ValeurAjoutee               *float64 `col:"Valeur ajoutée kEUR" json:"valeur_ajoutee" bson:"valeur_ajoutee,omitempty"`
	ChargePersonnel             *float64 `col:"Charges de personnel kEUR" json:"charge_personnel" bson:"charge_personnel,omitempty"`
	ImpotsTaxes                 *float64 `col:"Impôts. taxes et vers. assimil. kEUR" json:"impots_taxes" bson:"impots_taxes,omitempty"`
	// TODO: retirer colonne redondante: "Impôts. taxes et vers.  assimil. kEUR"
	SubventionsDExploitation      *float64 `col:"Subventions d'expl. kEUR" json:"subventions_d_exploitation" bson:"subventions_d_exploitation,omitempty"`
	ExcedentBrutDExploitation     *float64 `col:"Excédent brut d'exploitation kEUR" json:"excedent_brut_d_exploitation" bson:"excedent_brut_d_exploitation,omitempty"`
	AutresProduitsChargesReprises *float64 `col:"Autres Prod.. char. et Repr. kEUR" json:"autres_produits_charges_reprises" bson:"autres_produits_charges_reprises,omitempty"`
	// TODO: retirer colonne redondante: "Autres Prod.. char. et Repr.  kEUR"
	DotationAmortissement *float64 `col:"Dot. d'exploit. aux amort. et  prov. kEUR" json:"dotation_amortissement" bson:"dotation_amortissement,omitempty"`
	ResultatExpl          *float64 `col:"Résultat d'expl. kEUR" json:"resultat_expl" bson:"resultat_expl,omitempty"`
	OperationsCommun      *float64 `col:"Opérations en commun kEUR" json:"operations_commun" bson:"operations_commun,omitempty"`
	ProduitsFinanciers    *float64 `col:"Produits fin. kEUR" json:"produits_financiers" bson:"produits_financiers,omitempty"`
	ChargesFinancieres    *float64 `col:"Charges fin. kEUR" json:"charges_financieres" bson:"charges_financieres,omitempty"`
	Interets              *float64 `col:"Intérêts et charges assimilées kEUR" json:"interets" bson:"interets,omitempty"`
	ResultatAvantImpot    *float64 `col:"Résultat courant avant impôts kEUR" json:"resultat_avant_impot" bson:"resultat_avant_impot,omitempty"`
	ProduitExceptionnel   *float64 `col:"Produits except. kEUR" json:"produit_exceptionnel" bson:"produit_exceptionnel,omitempty"`
	ChargeExceptionnelle  *float64 `col:"Charges except. kEUR" json:"charge_exceptionnelle" bson:"charge_exceptionnelle,omitempty"`
	ParticipationSalaries *float64 `col:"Particip. des sal. aux résul. kEUR" json:"participation_salaries" bson:"participation_salaries,omitempty"`
	// TODO: retirer colonne redondante: "Particip. des sal. aux résul.  kEUR"
	ImpotBenefice *float64 `col:"Impôts sur le bénéf. et impôts diff. kEUR" json:"impot_benefice" bson:"impot_benefice,omitempty"`
	// TODO: retirer colonne redondante: "Impôts sur le bénéf. et impôts diff.  kEUR"
	BeneficeOuPerte *float64 `col:"Bénéfice ou perte kEUR" json:"benefice_ou_perte" bson:"benefice_ou_perte,omitempty"`
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
	var readCloser io.Reader
	parser.closeFct, readCloser, err = openFile(filePath)
	if err == nil {
		err = parser.initCsvReader(readCloser)
	}
	return err
}

func (parser *dianeParser) Close() error {
	return parser.closeFct()
}

func openFile(filePath string) (func() error, io.Reader, error) {

	pipedCmds := []*exec.Cmd{
		exec.Command("cat", filePath),                                          // TODO: implement this step in Go
		exec.Command("iconv", "--from-code", "UTF-16LE", "--to-code", "UTF-8"), // TODO: implement this step in Go
		exec.Command("awk", awkScript),                                         // TODO: implement this step in Go
		exec.Command("sed", "s/,/./g"),                                         // TODO: implement this step in Go
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

	return close, stdout, nil
}

func (parser *dianeParser) initCsvReader(reader io.Reader) (err error) {
	parser.reader = csv.NewReader(reader)
	parser.reader.Comma = ';'
	parser.reader.LazyQuotes = true

	header, err := parser.reader.Read()
	if err != nil {
		return errors.New("echec de lecture de l'en-tête du fichier en sortie du script: " + err.Error())
	}

	parser.idx, err = marshal.ValidateAndIndexColumnsFromColTags(header, Diane{})
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
			parsedLine.AddTuple(parseDianeRow(&parser.idx, row))
		}
		parsedLineChan <- parsedLine
	}
}

// parseDianeRow construit un objet Diane à partir d'une ligne de valeurs récupérée depuis un fichier
func parseDianeRow(idx *marshal.ColMapping, row []string) (diane Diane) {
	if i, err := strconv.Atoi(row[0]); err == nil {
		diane.Annee = &i
	}
	diane.NomEntreprise = row[2]
	diane.NumeroSiren = row[3]
	diane.StatutJuridique = row[4]
	diane.ProcedureCollective = (row[5] == "Oui")

	if i, err := strconv.Atoi(row[6]); err == nil {
		diane.EffectifConsolide = &i
	}
	if i, err := strconv.ParseFloat(row[7], 64); err == nil {
		diane.DetteFiscaleEtSociale = &i
	}
	if i, err := strconv.ParseFloat(row[8], 64); err == nil {
		diane.FraisDeRetD = &i
	}
	if i, err := strconv.ParseFloat(row[9], 64); err == nil {
		diane.ConcesBrevEtDroitsSim = &i
	}
	if i, err := strconv.Atoi(row[10]); err == nil {
		diane.NombreEtabSecondaire = &i
	}
	if i, err := strconv.Atoi(row[11]); err == nil {
		diane.NombreFiliale = &i
	}
	if i, err := strconv.Atoi(row[12]); err == nil {
		diane.TailleCompoGroupe = &i
	}
	if i, err := time.Parse("02/01/2006", row[14]); err == nil {
		diane.ArreteBilan = i
	}
	if i, err := strconv.Atoi(row[15]); err == nil {
		diane.NombreMois = &i
	}
	if i, err := strconv.ParseFloat(row[16], 64); err == nil {
		diane.ConcoursBancaireCourant = &i
	}
	if i, err := strconv.ParseFloat(row[17], 64); err == nil {
		diane.EquilibreFinancier = &i
	}
	if i, err := strconv.ParseFloat(row[18], 64); err == nil {
		diane.IndependanceFinanciere = &i
	}
	if i, err := strconv.ParseFloat(row[19], 64); err == nil {
		diane.Endettement = &i
	}
	if i, err := strconv.ParseFloat(row[20], 64); err == nil {
		diane.AutonomieFinanciere = &i
	}
	if i, err := strconv.ParseFloat(row[21], 64); err == nil {
		diane.DegreImmoCorporelle = &i
	}
	if i, err := strconv.ParseFloat(row[22], 64); err == nil {
		diane.FinancementActifCirculant = &i
	}
	if i, err := strconv.ParseFloat(row[23], 64); err == nil {
		diane.LiquiditeGenerale = &i
	}
	if i, err := strconv.ParseFloat(row[24], 64); err == nil {
		diane.LiquiditeReduite = &i
	}
	if i, err := strconv.ParseFloat(row[25], 64); err == nil {
		diane.RotationStocks = &i
	}
	if i, err := strconv.ParseFloat(row[26], 64); err == nil {
		diane.CreditClient = &i
	}
	if i, err := strconv.ParseFloat(row[27], 64); err == nil {
		diane.CreditFournisseur = &i
	}
	if i, err := strconv.ParseFloat(row[28], 64); err == nil {
		diane.CAparEffectif = &i
	}
	if i, err := strconv.ParseFloat(row[29], 64); err == nil {
		diane.TauxInteretFinancier = &i
	}
	if i, err := strconv.ParseFloat(row[30], 64); err == nil {
		diane.TauxInteretSurCA = &i
	}
	if i, err := strconv.ParseFloat(row[31], 64); err == nil {
		diane.EndettementGlobal = &i
	}
	if i, err := strconv.ParseFloat(row[32], 64); err == nil {
		diane.TauxEndettement = &i
	}
	if i, err := strconv.ParseFloat(row[33], 64); err == nil {
		diane.CapaciteRemboursement = &i
	}
	if i, err := strconv.ParseFloat(row[34], 64); err == nil {
		diane.CapaciteAutofinancement = &i
	}
	if i, err := strconv.ParseFloat(row[35], 64); err == nil {
		diane.CouvertureCaFdr = &i
	}
	if i, err := strconv.ParseFloat(row[36], 64); err == nil {
		diane.CouvertureCaBesoinFdr = &i
	}
	if i, err := strconv.ParseFloat(row[37], 64); err == nil {
		diane.PoidsBFRExploitation = &i
	}
	if i, err := strconv.ParseFloat(row[38], 64); err == nil {
		diane.Exportation = &i
	}
	if i, err := strconv.ParseFloat(row[39], 64); err == nil {
		diane.EfficaciteEconomique = &i
	}
	if i, err := strconv.ParseFloat(row[40], 64); err == nil {
		diane.ProductivitePotentielProduction = &i
	}
	if i, err := strconv.ParseFloat(row[41], 64); err == nil {
		diane.ProductiviteCapitalFinancier = &i
	}
	if i, err := strconv.ParseFloat(row[42], 64); err == nil {
		diane.ProductiviteCapitalInvesti = &i
	}
	if i, err := strconv.ParseFloat(row[43], 64); err == nil {
		diane.TauxDInvestissementProductif = &i
	}
	if i, err := strconv.ParseFloat(row[44], 64); err == nil {
		diane.RentabiliteEconomique = &i
	}
	if i, err := strconv.ParseFloat(row[45], 64); err == nil {
		diane.Performance = &i
	}
	if i, err := strconv.ParseFloat(row[46], 64); err == nil {
		diane.RendementBrutFondsPropres = &i
	}
	if i, err := strconv.ParseFloat(row[47], 64); err == nil {
		diane.RentabiliteNette = &i
	}
	if i, err := strconv.ParseFloat(row[48], 64); err == nil {
		diane.RendementCapitauxPropres = &i
	}
	if i, err := strconv.ParseFloat(row[49], 64); err == nil {
		diane.RendementRessourcesDurables = &i
	}
	if i, err := strconv.ParseFloat(row[50], 64); err == nil {
		diane.TauxMargeCommerciale = &i
	}
	if i, err := strconv.ParseFloat(row[51], 64); err == nil {
		diane.TauxValeurAjoutee = &i
	}
	if i, err := strconv.ParseFloat(row[52], 64); err == nil {
		diane.PartSalaries = &i
	}
	if i, err := strconv.ParseFloat(row[53], 64); err == nil {
		diane.PartEtat = &i
	}
	if i, err := strconv.ParseFloat(row[54], 64); err == nil {
		diane.PartPreteur = &i
	}
	if i, err := strconv.ParseFloat(row[55], 64); err == nil {
		diane.PartAutofinancement = &i
	}
	if i, err := strconv.ParseFloat(row[56], 64); err == nil {
		diane.CA = &i
	}
	if i, err := strconv.ParseFloat(row[57], 64); err == nil {
		diane.CAExportation = &i
	}
	if i, err := strconv.ParseFloat(row[58], 64); err == nil {
		diane.AchatMarchandises = &i
	}
	if i, err := strconv.ParseFloat(row[60], 64); err == nil {
		diane.AchatMatieresPremieres = &i
	}
	if i, err := strconv.ParseFloat(row[61], 64); err == nil {
		diane.Production = &i
	}
	if i, err := strconv.ParseFloat(row[62], 64); err == nil {
		diane.MargeCommerciale = &i
	}
	if i, err := strconv.ParseFloat(row[63], 64); err == nil {
		diane.Consommation = &i
	}
	if i, err := strconv.ParseFloat(row[64], 64); err == nil {
		diane.AutresAchatsChargesExternes = &i
	}
	if i, err := strconv.ParseFloat(row[65], 64); err == nil {
		diane.ValeurAjoutee = &i
	}
	if i, err := strconv.ParseFloat(row[66], 64); err == nil {
		diane.ChargePersonnel = &i
	}
	if i, err := strconv.ParseFloat(row[67], 64); err == nil {
		diane.ImpotsTaxes = &i
	}
	if i, err := strconv.ParseFloat(row[68], 64); err == nil {
		diane.SubventionsDExploitation = &i
	}
	if i, err := strconv.ParseFloat(row[69], 64); err == nil {
		diane.ExcedentBrutDExploitation = &i
	}
	if i, err := strconv.ParseFloat(row[70], 64); err == nil {
		diane.AutresProduitsChargesReprises = &i
	}
	if i, err := strconv.ParseFloat(row[71], 64); err == nil {
		diane.DotationAmortissement = &i
	}
	if i, err := strconv.ParseFloat(row[72], 64); err == nil {
		diane.ResultatExpl = &i
	}
	if i, err := strconv.ParseFloat(row[73], 64); err == nil {
		diane.OperationsCommun = &i
	}
	if i, err := strconv.ParseFloat(row[74], 64); err == nil {
		diane.ProduitsFinanciers = &i
	}
	if i, err := strconv.ParseFloat(row[75], 64); err == nil {
		diane.ChargesFinancieres = &i
	}
	if i, err := strconv.ParseFloat(row[76], 64); err == nil {
		diane.Interets = &i
	}
	if i, err := strconv.ParseFloat(row[77], 64); err == nil {
		diane.ResultatAvantImpot = &i
	}
	if i, err := strconv.ParseFloat(row[78], 64); err == nil {
		diane.ProduitExceptionnel = &i
	}
	if i, err := strconv.ParseFloat(row[79], 64); err == nil {
		diane.ChargeExceptionnelle = &i
	}
	if i, err := strconv.ParseFloat(row[80], 64); err == nil {
		diane.ParticipationSalaries = &i
	}
	if i, err := strconv.ParseFloat(row[81], 64); err == nil {
		diane.ImpotBenefice = &i
	}
	if i, err := strconv.ParseFloat(row[82], 64); err == nil {
		diane.BeneficeOuPerte = &i
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
