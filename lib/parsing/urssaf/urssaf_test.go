package urssaf

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

var notFoundRegexp = "column [A-Za-z_]+ not found"

var update = flag.Bool("update", false, "update the expected test values in golden file")

func TestUrssaf(t *testing.T) {
	t.Run("URSSAF gzipped files can be decompressed on the fly", func(t *testing.T) {
		type TestCase struct {
			Parser     engine.Parser
			InputFile  string
			GoldenFile string
		}
		urssafFiles := []TestCase{
			{NewCCSFParser(), "ccsfTestData.csv", "expectedCcsf.json"},
			{NewDebitParser(), "debitTestData.csv", "expectedDebit.json"},
			{NewDelaiParser(), "delaiTestData.csv", "expectedDelai.json"},
			{NewProcolParser(), "procolTestData.csv", "expectedProcol.json"},
		}
		for _, testCase := range urssafFiles {
			t.Run(string(testCase.Parser.Type()), func(t *testing.T) {
				// Compression du fichier de données
				err := exec.Command("gzip", "--keep", filepath.Join("testData", testCase.InputFile)).Run() // créée une version gzippée du fichier
				assert.NoError(t, err)
				compressedFilePath := engine.NewBatchFile("testData", testCase.InputFile+".gz")
				t.Cleanup(func() { os.Remove(compressedFilePath.Path()) })

				// Création d'un fichier Golden temporaire mentionnant le nom du fichier compressé
				initialGoldenContent, err := os.ReadFile(filepath.Join("testData", testCase.GoldenFile))
				assert.NoError(t, err)
				goldenContent := bytes.ReplaceAll(initialGoldenContent, []byte(testCase.InputFile), []byte(testCase.InputFile+".gz"))
				tmpGoldenFile := engine.CreateTempFileWithContent(t, goldenContent)

				engine.TestParserOutput(t, testCase.Parser, compressedFilePath, tmpGoldenFile.Name(), false)
			})
		}
	})
}

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = engine.NewBatchFile("testData", "debitTestData.csv")

	engine.TestParserOutput(t, NewDebitParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDebitParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestDebitParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected Debit
	}{
		{
			csvRow: []string{
				"12345678901234", "123456789", "ECN001", "1180515", "201811",
				"123456", "789012", "1", "2", "01", "02", "MOT01", "true",
			},
			expected: Debit{
				Siret:                        "12345678901234",
				NumeroCompte:                 "123456789",
				NumeroEcartNegatif:           "ECN001",
				DateTraitement:               parsing.MustParseTime("2006-01-02", "2018-05-15"),
				PeriodePriseEnCompte:         parsing.MustParseTime("2006-01-02", "2018-05-01"),
				PartOuvriere:                 1234.56,
				PartPatronale:                7890.12,
				NumeroHistoriqueEcartNegatif: parsing.IntPtr(1),
				EtatCompte:                   parsing.IntPtr(2),
				CodeProcedureCollective:      "01",
				PeriodeDebut:                 parsing.MustParseTime("2006-01-02", "2018-01-01"),
				PeriodeFin:                   parsing.MustParseTime("2006-01-02", "2018-02-01"),
				CodeOperationEcartNegatif:    "02",
				CodeMotifEcartNegatif:        "MOT01",
				Recours:                      true,
				DebitID:                      "123456789012342018010120180201ECN001",
			},
		},
		{
			csvRow: []string{
				"98765432109876", "987654321", "ECN002", "1200325", "202010",
				"567890", "123450", "3", "4", "03", "04", "MOT02", "false",
			},
			expected: Debit{
				Siret:                        "98765432109876",
				NumeroCompte:                 "987654321",
				NumeroEcartNegatif:           "ECN002",
				DateTraitement:               parsing.MustParseTime("2006-01-02", "2020-03-25"),
				PeriodePriseEnCompte:         parsing.MustParseTime("2006-01-02", "2020-04-01"),
				PartOuvriere:                 5678.90,
				PartPatronale:                1234.50,
				NumeroHistoriqueEcartNegatif: parsing.IntPtr(3),
				EtatCompte:                   parsing.IntPtr(4),
				CodeProcedureCollective:      "03",
				PeriodeDebut:                 parsing.MustParseTime("2006-01-02", "2020-01-01"),
				PeriodeFin:                   parsing.MustParseTime("2006-01-02", "2020-04-01"),
				CodeOperationEcartNegatif:    "04",
				CodeMotifEcartNegatif:        "MOT02",
				Recours:                      false,
				DebitID:                      "987654321098762020010120200401ECN002",
			},
		},
	}

	parser := NewDebitParser()
	header := "Siret;num_cpte;Num_Ecn;Dt_trt_ecn;Periode;Mt_PO;Mt_PP;Num_Hist_Ecn;Etat_cpte;Cd_pro_col;Cd_op_ecn;Motif_ecn;Recours_en_cours"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		debit, ok := res.Tuples[0].(Debit)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, debit)
	}
}

func TestDebitMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDebitParser(), []string{"Siret;num_cpte"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})

	t.Run("should fail if Periode column is missing", func(t *testing.T) {
		headerRow := []string{"Siret;num_cpte;Num_Ecn;Dt_trt_ecn;Mt_PO;Mt_PP;Num_Hist_Ecn;Etat_cpte;Cd_pro_col;Cd_op_ecn;Motif_ecn;Recours_en_cours"}
		output := engine.RunParserInline(t, NewDebitParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Periode not found")
	})
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = engine.NewBatchFile("testData", "debitCorrompuTestData.csv")
	engine.TestParserOutput(t, NewDebitParser(), testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = engine.NewBatchFile("testData", "delaiTestData.csv")
	engine.TestParserOutput(t, NewDelaiParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDelaiParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestDelaiParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected Delai
	}{
		{
			csvRow: []string{
				"12345678901234", "123456789", "CONT001", "15/05/2018", "15/11/2018",
				"180", "Entreprise Test", "O", "2018", "12345,67", "01", "02",
			},
			expected: Delai{
				Siret:             "12345678901234",
				NumeroCompte:      "123456789",
				NumeroContentieux: "CONT001",
				DateCreation:      parsing.MustParseTime("02/01/2006", "15/05/2018"),
				DateEcheance:      parsing.MustParseTime("02/01/2006", "15/11/2018"),
				DureeDelai:        parsing.IntPtr(180),
				Denomination:      "Entreprise Test",
				Indic6m:           "O",
				AnneeCreation:     parsing.IntPtr(2018),
				MontantEcheancier: parsing.Float64Ptr(12345.67),
				Stade:             "01",
				Action:            "02",
			},
		},
		{
			csvRow: []string{
				"98765432109876", "987654321", "CONT002", "12/03/2020", "12/09/2020",
				"182", "Société ABC", "N", "2020", "98765,43", "03", "04",
			},
			expected: Delai{
				Siret:             "98765432109876",
				NumeroCompte:      "987654321",
				NumeroContentieux: "CONT002",
				DateCreation:      parsing.MustParseTime("02/01/2006", "12/03/2020"),
				DateEcheance:      parsing.MustParseTime("02/01/2006", "12/09/2020"),
				DureeDelai:        parsing.IntPtr(182),
				Denomination:      "Société ABC",
				Indic6m:           "N",
				AnneeCreation:     parsing.IntPtr(2020),
				MontantEcheancier: parsing.Float64Ptr(98765.43),
				Stade:             "03",
				Action:            "04",
			},
		},
	}

	parser := NewDelaiParser()
	header := "Siret;Numero_compte_externe;Numero_structure;Date_creation;Date_echeance;Duree_delai;Denomination_premiere_ligne;Indic_6M;Annee_creation;Montant_global_echeancier;Code_externe_stade;Code_externe_action"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		delai, ok := res.Tuples[0].(Delai)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, delai)
	}
}

func TestDelaiMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDelaiParser(), []string{"Siret;Numero_compte_externe"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})

	t.Run("should fail if Date_creation column is missing", func(t *testing.T) {
		headerRow := []string{"Siret;Numero_compte_externe;Numero_structure;Date_echeance"}
		output := engine.RunParserInline(t, NewDelaiParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Date_creation not found")
	})
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = engine.NewBatchFile("testData", "ccsfTestData.csv")
	engine.TestParserOutput(t, NewCCSFParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCCSFParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestCcsfParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected CCSF
	}{
		{
			csvRow: []string{"12345678901234", "123456789", "1180515", "01", "02"},
			expected: CCSF{
				Siret:          "12345678901234",
				NumeroCompte:   "123456789",
				DateTraitement: parsing.MustParseTime("2006-01-02", "2018-05-15"),
				Stade:          "01",
				Action:         "02",
			},
		},
		{
			csvRow: []string{"98765432109876", "987654321", "1200312", "03", "04"},
			expected: CCSF{
				Siret:          "98765432109876",
				NumeroCompte:   "987654321",
				DateTraitement: parsing.MustParseTime("2006-01-02", "2020-03-12"),
				Stade:          "03",
				Action:         "04",
			},
		},
		{
			csvRow: []string{"11122233344455", "111222333", "1150101", "05", "06"},
			expected: CCSF{
				Siret:          "11122233344455",
				NumeroCompte:   "111222333",
				DateTraitement: parsing.MustParseTime("2006-01-02", "2015-01-01"),
				Stade:          "05",
				Action:         "06",
			},
		},
	}

	parser := NewCCSFParser()
	header := "Siret;Compte;Date_de_traitement;Code_externe_stade;Code_externe_action"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		ccsf, ok := res.Tuples[0].(CCSF)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, ccsf)
	}
}

func TestCcsfMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCCSFParser(), []string{"Siret;Compte"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})

	t.Run("should fail if Date_de_traitement column is missing", func(t *testing.T) {
		headerRow := []string{"Siret;Compte;Code_externe_stade;Code_externe_action"}
		output := engine.RunParserInline(t, NewCCSFParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Date_de_traitement not found")
	})
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = engine.NewBatchFile("testData", "cotisationTestData.csv")
	engine.TestParserOutput(t, NewCotisationParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCotisationParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestCotisationParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected Cotisation
	}{
		{
			csvRow: []string{"12345678901234", "123456789", "201811", "1234,56", "7890,12"},
			expected: Cotisation{
				Siret:        "12345678901234",
				NumeroCompte: "123456789",
				PeriodeDebut: parsing.MustParseTime("2006-01-02", "2018-01-01"),
				PeriodeFin:   parsing.MustParseTime("2006-01-02", "2018-02-01"),
				Encaisse:     parsing.Float64Ptr(1234.56),
				Du:           parsing.Float64Ptr(7890.12),
			},
		},
		{
			csvRow: []string{"98765432109876", "987654321", "202010", "5678,90", "1234,50"},
			expected: Cotisation{
				Siret:        "98765432109876",
				NumeroCompte: "987654321",
				PeriodeDebut: parsing.MustParseTime("2006-01-02", "2020-01-01"),
				PeriodeFin:   parsing.MustParseTime("2006-01-02", "2020-04-01"),
				Encaisse:     parsing.Float64Ptr(5678.90),
				Du:           parsing.Float64Ptr(1234.50),
			},
		},
		{
			csvRow: []string{"11122233344455", "111222333", "201962", "999,99", "888,88"},
			expected: Cotisation{
				Siret:        "11122233344455",
				NumeroCompte: "111222333",
				PeriodeDebut: parsing.MustParseTime("2006-01-02", "2019-01-01"),
				PeriodeFin:   parsing.MustParseTime("2006-01-02", "2020-01-01"),
				Encaisse:     parsing.Float64Ptr(999.99),
				Du:           parsing.Float64Ptr(888.88),
			},
		},
	}

	parser := NewCotisationParser()
	header := "Siret;Compte;periode;enc_direct;cotis_due"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		cotisation, ok := res.Tuples[0].(Cotisation)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, cotisation)
	}
}

func TestCotisationMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCotisationParser(), []string{"Siret;Compte"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})

	t.Run("should fail if periode column is missing", func(t *testing.T) {
		headerRow := []string{"Siret;Compte;enc_direct;cotis_due"}
		output := engine.RunParserInline(t, NewCotisationParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "periode not found")
	})
}

func TestProcol(t *testing.T) {
	t.Run("Le fichier de test Procol est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedProcol.json")
		var testData = engine.NewBatchFile("testData", "procolTestData.csv")
		engine.TestParserOutput(t, NewProcolParser(), testData, golden, *update)
	})

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewProcolParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "not found")
	})

	t.Run("est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewProcolParser(), []string{"dT_eFfeT;lIb_aCtx_stDx;sIret"})
		assert.Len(t, output.Reports[0].HeadFatal, 0)
	})
}

func TestProcolParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected Procol
	}{
		{
			csvRow: []string{"12345678901234", "15May2018", "ouverture_liquidation"},
			expected: Procol{
				Siret:        "12345678901234",
				DateEffet:    parsing.MustParseTime("02Jan2006", "15May2018"),
				ActionProcol: "liquidation",
				StadeProcol:  "ouverture",
			},
		},
		{
			csvRow: []string{"98765432109876", "12Mar2020", "conversion_redressement"},
			expected: Procol{
				Siret:        "98765432109876",
				DateEffet:    parsing.MustParseTime("02Jan2006", "12Mar2020"),
				ActionProcol: "redressement",
				StadeProcol:  "conversion",
			},
		},
		{
			csvRow: []string{"11122233344455", "01Jan2019", "jugement_sauvegarde"},
			expected: Procol{
				Siret:        "11122233344455",
				DateEffet:    parsing.MustParseTime("02Jan2006", "01Jan2019"),
				ActionProcol: "sauvegarde",
				StadeProcol:  "jugement",
			},
		},
	}

	parser := NewProcolParser()
	header := "Siret;Dt_effet;Lib_actx_stdx"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		procol, ok := res.Tuples[0].(Procol)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, procol)
	}
}

func TestProcolMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewProcolParser(), []string{"Siret"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "not found")
	})

	t.Run("should fail if Dt_effet column is missing", func(t *testing.T) {
		headerRow := []string{"Siret;Lib_actx_stdx"}
		output := engine.RunParserInline(t, NewProcolParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Dt_effet not found")
	})
}
