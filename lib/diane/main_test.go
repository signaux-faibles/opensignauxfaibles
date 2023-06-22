package diane

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {

	t.Run("initCsvReader() doit détecter s'il manque des colonnes", func(t *testing.T) {
		var parser = Parser
		err := parser.initCsvReader(strings.NewReader("AnneeS;Marquée;Nom de l'entreprise;Numéro Siren;Statut juridique ;Procédure collective;Effectif consolidé;Dettes fiscales et sociales kEUR;Frais de R&D : net kEUR;Conces., brev. et droits sim. : net kEUR;Nombre d\"ES;Nombre de filiales;Taille de la Composition du Groupe;Dernière année disponible;Date de clôture;Nombre de mois;Conc.  banc. cour. & sold. cr. kEUR;Conc. banc. cour. & sold. cr. kEUR;Conc. banc. cour. & sold. cr.  kEUR;Equilibre financier;Indépendance fin. %;Indépendance fin.  %;Endettement %;Autonomie fin. %;Degré d'amort. des immob. corp. %;Degré d'amort.  des immob. corp. %;Financ. de l'actif circ. net;Liquidité générale;Liquidité réduite;Rotation des stocks jours;Crédit clients jours;Crédit fournisseurs jours;C. A.  par effectif (milliers/pers.) kEUR;C. A. par effectif (milliers/pers.) kEUR;Taux d'intérêt financier %;Intérêts / Chiffre d'affaires %;Endettement global jours;Taux d'endettement %;Capacité de remboursement;Capacité d'autofin. %;Couv. du C.A. par le f.d.r. jours;Couv. du C.A. par le f.d.r.  jours;Couv. du C.A. par bes. en fdr jours;Poids des BFR d'exploitation %;Exportation %;Efficacité économique (milliers/pers.) kEUR;Prod. du potentiel de production;Prod.  du potentiel de production;Productivité du capital financier;Productivité du capital investi;Taux d'invest. productif %;Rentabilité économique %;Performance %;Rend. brut des f. propres nets %;Rend.  brut des f. propres nets %;Rentabilité nette %;Rend. des capitaux propres nets %;Rend.  des capitaux propres nets %;Rend. des res. durables nettes %;Rend. des res.  durables nettes %;Taux de marge commerciale %;Taux de valeur ajoutée %;Part des salariés %;Part de l'Etat %;Part des prêteurs %;Part de l'autofin. %;Chiffre d'affaires net (H.T.) kEUR;Dont exportation kEUR;Achats march. et autres approv. kEUR;Achats march.  et autres approv. kEUR;Achats de march. kEUR;Achats de mat. prem. et autres approv. kEUR;Achats de mat. prem.  et autres approv. kEUR;Production de l'ex. kEUR;Marge commerciale kEUR;Consommation de l'ex. kEUR;Autres achats et charges externes kEUR;Valeur ajoutée kEUR;Charges de personnel kEUR;Impôts, taxes et vers. assimil. kEUR;Impôts, taxes et vers.  assimil. kEUR;Subventions d'expl. kEUR;Excédent brut d'exploitation kEUR;Autres Prod., char. et Repr. kEUR;Autres Prod., char. et Repr.  kEUR;Dot. d'exploit. aux amort. et  prov. kEUR;Résultat d'expl. kEUR;Opérations en commun kEUR;Produits fin. kEUR;Charges fin. kEUR;Intérêts et charges assimilées kEUR;Résultat courant avant impôts kEUR;Produits except. kEUR;Charges except. kEUR;Particip. des sal. aux résul. kEUR;Particip. des sal. aux résul.  kEUR;Impôts sur le bénéf. et impôts diff. kEUR;Impôts sur le bénéf. et impôts diff.  kEUR;Bénéfice ou perte kEUR")) // un S a été ajoutée à "Annee"
		if assert.Error(t, err, "initCsvReader() devrait échouer") {
			assert.Contains(t, err.Error(), "Colonne Annee non trouvée")
		}
	})

	t.Run("preprocessDianeFile() doit retourner un flux csv avec un en-tête sans duplication de caractères d'espacement", func(t *testing.T) {
		var testData = filepath.Join("testData", "dianeTestData.txt")
		_, reader, err := preprocessDianeFile(testData)
		if assert.NoError(t, err) {
			csvReader := csv.NewReader(reader)
			csvReader.Comma = ';'
			csvReader.LazyQuotes = true
			header, err := csvReader.Read() // Discard header
			if assert.NoError(t, err) {
				for _, field := range header {
					assert.NotContains(t, field, "  ")
				}
			}
		}
	})

	t.Run("preprocessDianeFile() doit produire un fichier csv intermédiaire conforme", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDianeConvert.csv")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		_, reader, err := preprocessDianeFile(testData)
		if assert.NoError(t, err) {
			output, err := ioutil.ReadAll(reader)
			if assert.NoError(t, err) {
				diffWithGoldenFile(t, output, golden, *update)
			}
		}
	})

	t.Run("Diane parser (JSON output)", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDiane.json")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
	})
}

func diffWithGoldenFile(t *testing.T, output []byte, goldenFileName string, updateGoldenFile bool) {
	t.Helper()
	if updateGoldenFile {
		os.WriteFile(goldenFileName, output, 0644)
	}
	expected, err := os.ReadFile(goldenFileName)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, output)
	}
}
