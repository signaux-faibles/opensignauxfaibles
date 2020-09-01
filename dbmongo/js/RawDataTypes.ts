import { EntréeDelai } from "./common/validDelai"

// Données importées pour une entreprise ou établissement

export type Scope = "etablissement" | "entreprise"

export type CompanyDataValues = {
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

export type CompanyDataValuesWithFlags = CompanyDataValues & IndexFlags

export type IndexFlags = {
  index: {
    algo1: boolean
    algo2: boolean
  }
}

// Données importées par les parseurs, pour chaque source de données

export type BatchKey = string

export type BatchValues = Record<BatchKey, BatchValue>

export type DataType = keyof BatchValueProps // => 'reporder' | 'effectif' | 'apconso' | ...

export type BatchValue = Partial<BatchValueProps>

type BatchValueProps = {
  reporder: Record<Periode, EntréeRepOrder> // RepOrder est généré par "compact", et non importé => Usage de Periode en guise de hash d'indexation
  effectif: ParHash<EntréeEffectif>
  apconso: ParHash<EntréeApConso>
  apdemande: ParHash<EntréeApDemande>
  compte: ParHash<EntréeCompte>
  interim: ParHash<EntréeInterim>
  delai: ParHash<EntréeDelai>
  altares: ParHash<EntréeDefaillances>
  procol: ParHash<EntréeDefaillances>
  cotisation: ParHash<EntréeCotisation>
  debit: ParHash<EntréeDebit>
  ccsf: ParHash<{ date_traitement: Date }>
  sirene: ParHash<EntréeSirene>
  sirene_ul: ParHash<EntréeSireneEntreprise>
  effectif_ent: ParHash<EntréeEffectif>
  bdf: ParHash<EntréeBdf>
  diane: ParHash<EntréeDiane>
}
