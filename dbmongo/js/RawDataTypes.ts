// Données importées pour une entreprise ou établissement

export type Scope = "etablissement" | "entreprise"

export type CompanyDataValues = {
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

export type CompanyDataValuesWithFlags = CompanyDataValues & {
  index: {
    algo1: boolean
    algo2: boolean
  }
}

// Données importées par les parseurs, pour chaque source de données

export type BatchKey = string

export type BatchValues = Record<BatchKey, BatchValue>

export type BatchValue = Partial<{
  reporder: Record<Periode, EntréeRepOrder> // RepOrder est généré, et non importé => Usage de Periode en guise de hash d'indexation
  compact: { delete: Partial<Record<DataType, DataHash[]>> }
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
}>

export type DataType = Exclude<keyof BatchValue, "compact"> // => 'reporder' | 'effectif' | 'apconso' | ...
