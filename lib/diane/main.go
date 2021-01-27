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
	Annee                           *int      `json:"exercice_diane,omitempty" bson:"exercice_diane,omitempty"`
	NomEntreprise                   string    `json:"nom_entreprise" bson:"nom_entreprise,omitempty"`
	NumeroSiren                     string    `json:"numero_siren" bson:"numero_siren,omitempty"`
	StatutJuridique                 string    `json:"statut_juridique" bson:"statut_juridique,omitempty"`
	ProcedureCollective             bool      `json:"procedure_collective" bson:"procedure_collective,omitempty"`
	EffectifConsolide               *int      `json:"effectif_consolide" bson:"effectif_consolide,omitempty"`
	DetteFiscaleEtSociale           *float64  `json:"dette_fiscale_et_sociale" bson:"dette_fiscale_et_sociale,omitempty"`
	FraisDeRetD                     *float64  `json:"frais_de_RetD" bson:"frais_de_RetD,omitempty"`
	ConcesBrevEtDroitsSim           *float64  `json:"conces_brev_et_droits_sim" bson:"conces_brev_et_droits_sim,omitempty"`
	NombreEtabSecondaire            *int      `json:"nombre_etab_secondaire" bson:"nombre_etab_secondaire,omitempty"`
	NombreFiliale                   *int      `json:"nombre_filiale" bson:"nombre_filiale,omitempty"`
	TailleCompoGroupe               *int      `json:"taille_compo_groupe" bson:"taille_compo_groupe,omitempty"`
	ArreteBilan                     time.Time `json:"arrete_bilan_diane" bson:"arrete_bilan_diane,omitempty"`
	NombreMois                      *int      `json:"nombre_mois" bson:"nombre_mois,omitempty"`
	ConcoursBancaireCourant         *float64  `json:"concours_bancaire_courant" bson:"concours_bancaire_courant,omitempty"`
	EquilibreFinancier              *float64  `json:"equilibre_financier" bson:"equilibre_financier,omitempty"`
	IndependanceFinanciere          *float64  `json:"independance_financiere" bson:"independance_financiere,omitempty"`
	Endettement                     *float64  `json:"endettement" bson:"endettement,omitempty"`
	AutonomieFinanciere             *float64  `json:"autonomie_financiere" bson:"autonomie_financiere,omitempty"`
	DegreImmoCorporelle             *float64  `json:"degre_immo_corporelle" bson:"degre_immo_corporelle,omitempty"`
	FinancementActifCirculant       *float64  `json:"financement_actif_circulant" bson:"financement_actif_circulant,omitempty"`
	LiquiditeGenerale               *float64  `json:"liquidite_generale" bson:"liquidite_generale,omitempty"`
	LiquiditeReduite                *float64  `json:"liquidite_reduite" bson:"liquidite_reduite,omitempty"`
	RotationStocks                  *float64  `json:"rotation_stocks" bson:"rotation_stocks,omitempty"`
	CreditClient                    *float64  `json:"credit_client" bson:"credit_client,omitempty"`
	CreditFournisseur               *float64  `json:"credit_fournisseur" bson:"credit_fournisseur,omitempty"`
	CAparEffectif                   *float64  `json:"ca_par_effectif" bson:"ca_apar_effectif,omitempty"`
	TauxInteretFinancier            *float64  `json:"taux_interet_financier" bson:"taux_interet_financier,omitempty"`
	TauxInteretSurCA                *float64  `json:"taux_interet_sur_ca" bson:"taux_interet_sur_ca,omitempty"`
	EndettementGlobal               *float64  `json:"endettement_global" bson:"endettement_global,omitempty"`
	TauxEndettement                 *float64  `json:"taux_endettement" bson:"taux_endettement,omitempty"`
	CapaciteRemboursement           *float64  `json:"capacite_remboursement" bson:"capacite_remboursement,omitempty"`
	CapaciteAutofinancement         *float64  `json:"capacite_autofinancement" bson:"capacite_autofinancement,omitempty"`
	CouvertureCaFdr                 *float64  `json:"couverture_ca_fdr" bson:"couverture_ca_fdr,omitempty"`
	CouvertureCaBesoinFdr           *float64  `json:"couverture_ca_besoin_fdr" bson:"couverture_ca_besoin_fdr,omitempty"`
	PoidsBFRExploitation            *float64  `json:"poids_bfr_exploitation" bson:"poids_bfr_exploitation,omitempty"`
	Exportation                     *float64  `json:"exportation" bson:"exportation,omitempty"`
	EfficaciteEconomique            *float64  `json:"efficacite_economique" bson:"efficacite_economique,omitempty"`
	ProductivitePotentielProduction *float64  `json:"productivite_potentiel_production" bson:"productivite_potentiel_production,omitempty"`
	ProductiviteCapitalFinancier    *float64  `json:"productivite_capital_financier" bson:"productivite_capital_financier,omitempty"`
	ProductiviteCapitalInvesti      *float64  `json:"productivite_capital_investi" bson:"productivite_capital_investi,omitempty"`
	TauxDInvestissementProductif    *float64  `json:"taux_d_investissement_productif" bson:"taux_d_investissement_productif,omitempty"`
	RentabiliteEconomique           *float64  `json:"rentabilite_economique" bson:"rentabilite_economique,omitempty"`
	Performance                     *float64  `json:"performance" bson:"performance,omitempty"`
	RendementBrutFondsPropres       *float64  `json:"rendement_brut_fonds_propres" bson:"rendement_brut_fonds_propres,omitempty"`
	RentabiliteNette                *float64  `json:"rentabilite_nette" bson:"rentabilite_nette,omitempty"`
	RendementCapitauxPropres        *float64  `json:"rendement_capitaux_propres" bson:"rendement_capitaux_propres,omitempty"`
	RendementRessourcesDurables     *float64  `json:"rendement_ressources_durables" bson:"rendement_ressources_durables,omitempty"`
	TauxMargeCommerciale            *float64  `json:"taux_marge_commerciale" bson:"taux_marge_commerciale,omitempty"`
	TauxValeurAjoutee               *float64  `json:"taux_valeur_ajoutee" bson:"taux_valeur_ajoutee,omitempty"`
	PartSalaries                    *float64  `json:"part_salaries" bson:"part_salaries,omitempty"`
	PartEtat                        *float64  `json:"part_etat" bson:"part_etat,omitempty"`
	PartPreteur                     *float64  `json:"part_preteur" bson:"part_preteur,omitempty"`
	PartAutofinancement             *float64  `json:"part_autofinancement" bson:"part_autofinancement,omitempty"`
	CA                              *float64  `json:"ca" bson:"ca,omitempty"`
	CAExportation                   *float64  `json:"ca_exportation" bson:"ca_exportation,omitempty"`
	AchatMarchandises               *float64  `json:"achat_marchandises" bson:"achat_marchandises,omitempty"`
	AchatMatieresPremieres          *float64  `json:"achat_matieres_premieres" bson:"achat_matieres_premieres,omitempty"`
	Production                      *float64  `json:"production" bson:"production,omitempty"`
	MargeCommerciale                *float64  `json:"marge_commerciale" bson:"marge_commerciale,omitempty"`
	Consommation                    *float64  `json:"consommation" bson:"consommation,omitempty"`
	AutresAchatsChargesExternes     *float64  `json:"autres_achats_charges_externes" bson:"autres_achats_charges_externes,omitempty"`
	ValeurAjoutee                   *float64  `json:"valeur_ajoutee" bson:"valeur_ajoutee,omitempty"`
	ChargePersonnel                 *float64  `json:"charge_personnel" bson:"charge_personnel,omitempty"`
	ImpotsTaxes                     *float64  `json:"impots_taxes" bson:"impots_taxes,omitempty"`
	SubventionsDExploitation        *float64  `json:"subventions_d_exploitation" bson:"subventions_d_exploitation,omitempty"`
	ExcedentBrutDExploitation       *float64  `json:"excedent_brut_d_exploitation" bson:"excedent_brut_d_exploitation,omitempty"`
	AutresProduitsChargesReprises   *float64  `json:"autres_produits_charges_reprises" bson:"autres_produits_charges_reprises,omitempty"`
	DotationAmortissement           *float64  `json:"dotation_amortissement" bson:"dotation_amortissement,omitempty"`
	ResultatExpl                    *float64  `json:"resultat_expl" bson:"resultat_expl,omitempty"`
	OperationsCommun                *float64  `json:"operations_commun" bson:"operations_commun,omitempty"`
	ProduitsFinanciers              *float64  `json:"produits_financiers" bson:"produits_financiers,omitempty"`
	ChargesFinancieres              *float64  `json:"charges_financieres" bson:"charges_financieres,omitempty"`
	Interets                        *float64  `json:"interets" bson:"interets,omitempty"`
	ResultatAvantImpot              *float64  `json:"resultat_avant_impot" bson:"resultat_avant_impot,omitempty"`
	ProduitExceptionnel             *float64  `json:"produit_exceptionnel" bson:"produit_exceptionnel,omitempty"`
	ChargeExceptionnelle            *float64  `json:"charge_exceptionnelle" bson:"charge_exceptionnelle,omitempty"`
	ParticipationSalaries           *float64  `json:"participation_salaries" bson:"participation_salaries,omitempty"`
	ImpotBenefice                   *float64  `json:"impot_benefice" bson:"impot_benefice,omitempty"`
	BeneficeOuPerte                 *float64  `json:"benefice_ou_perte" bson:"benefice_ou_perte,omitempty"`
	// TODO: ajouter NotePreface ou le retirer de la documentation (cf https://github.com/signaux-faibles/documentation/blob/master/description-donnees.md#donn%C3%A9es-financi%C3%A8res-issues-des-bilans-d%C3%A9pos%C3%A9s-au-greffe-de-tribunaux-de-commerce)
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
}

func (parser *dianeParser) GetFileType() string {
	return "diane"
}

func (parser *dianeParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *dianeParser) Open(filePath string) (err error) {
	var reader *io.ReadCloser
	parser.closeFct, reader, err = openFile(filePath)
	if err != nil {
		return err
	}

	// init csv reader
	parser.reader = csv.NewReader(*reader)
	parser.reader.Comma = ';'
	parser.reader.LazyQuotes = true

	_, err = parser.reader.Read() // Discard header
	if err != nil {
		return errors.New("echec de lecture de l'en-tête du fichier en sortie du script: " + err.Error())
	}

	return err
}

func (parser *dianeParser) Close() error {
	return parser.closeFct()
}

func openFile(filePath string) (func() error, *io.ReadCloser, error) {

	// TODO: implémenter ces traitements en Go
	pipedCmds := []*exec.Cmd{
		exec.Command("cat", filePath),
		exec.Command("iconv", "--from-code", "UTF-16LE", "--to-code", "UTF-8"), // conversion d'encodage de fichier car awk ne supporte pas UTF-16LE
		exec.Command("sed", "s/\r$//"),                                         // forcer l'usage de retours charriot au format UNIX car le caractère \r cause une duplication de colonne depuis le script awk
		exec.Command("awk", awkScript),                                         // réorganisation des des données pour éviter la duplication de colonnes par année
		exec.Command("sed", "s/,/./g"),                                         // usage de points au lieu de virgules, pour que les nombres décimaux soient reconnus par csv.Reader
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
			parseDianeLine(row, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseDianeLine(row []string, parsedLine *marshal.ParsedLineResult) {
	parsedLine.AddTuple(parseDianeRow(row))
}

// parseDianeRow construit un objet Diane à partir d'une ligne de valeurs récupérée depuis un fichier
func parseDianeRow(row []string) (diane Diane) {
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
