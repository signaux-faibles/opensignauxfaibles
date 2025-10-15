package urssaf

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	Siret                        string    `                         json:"-"                            sql:"siret"                            csv:"siret"`
	NumeroCompte                 string    `input:"num_cpte"         json:"numero_compte"                sql:"numero_compte"                    csv:"numéro_compte"`
	NumeroEcartNegatif           string    `input:"Num_Ecn"          json:"numero_ecart_negatif"         sql:"numero_ecart_negatif"             csv:"numéro_écart_négatif"`
	DateTraitement               time.Time `input:"Dt_trt_ecn"       json:"date_traitement"              sql:"date_traitement"                  csv:"date_traitement"`
	PartOuvriere                 float64   `input:"Mt_PO"            json:"part_ouvriere"                sql:"part_ouvriere"                    csv:"part_ouvrière"`
	PartPatronale                float64   `input:"Mt_PP"            json:"part_patronale"               sql:"part_patronale"                   csv:"part_patronale"`
	NumeroHistoriqueEcartNegatif *int      `input:"Num_Hist_Ecn"     json:"numero_historique"            sql:"numero_historique_ecart_negatif"  csv:"numéro_historique_écart_négatif"`
	EtatCompte                   *int      `input:"Etat_cpte"        json:"etat_compte"                  sql:"etat_compte"                      csv:"état_compte"`
	CodeProcedureCollective      string    `input:"Cd_pro_col"       json:"code_procedure_collective"    sql:"code_procedure_collective"        csv:"code_procédure_collective"`
	PeriodeDebut                 time.Time `input:"Periode"          json:"periode_debut"                sql:"periode_debut"                    csv:"période_début"`
	PeriodeFin                   time.Time `input:"Periode"          json:"periode_fin"                  sql:"periode_fin"                      csv:"période_fin"`
	CodeOperationEcartNegatif    string    `input:"Cd_op_ecn"        json:"code_operation_ecart_negatif" sql:"code_operation_ecart_negatif"     csv:"code_opération_écart_négatif"`
	CodeMotifEcartNegatif        string    `input:"Motif_ecn"        json:"code_motif_ecart_negatif"     sql:"code_motif_ecart_negatif"         csv:"code_motif_écart_négatif"`
	Recours                      bool      `input:"Recours_en_cours" json:"recours_en_cours"             sql:"recours_en_cours"                 csv:"recours"`
}

// Key _id de l'objet
func (debit Debit) Key() string {
	return debit.Siret
}

// Scope de l'objet
func (debit Debit) Scope() engine.Scope {
	return engine.ScopeEtablissement
}

// Type de l'objet
func (debit Debit) Type() engine.ParserType {
	return engine.Debit
}
