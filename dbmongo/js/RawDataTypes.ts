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
  effectif: Record<DataHash, EntréeEffectif>
  apconso: Record<DataHash, EntréeApConso>
  apdemande: Record<DataHash, EntréeApDemande>
  compte: Record<DataHash, EntréeCompte>
  interim: Record<DataHash, EntréeInterim>
  delai: Record<DataHash, EntréeDelai>
  altares: Record<DataHash, EntréeDefaillances>
  procol: Record<DataHash, EntréeDefaillances>
  cotisation: Record<DataHash, EntréeCotisation>
  debit: Record<DataHash, EntréeDebit>
  ccsf: Record<DataHash, { date_traitement: Date }>
  sirene: Record<DataHash, EntréeSirene>
  sirene_ul: Record<DataHash, EntréeSireneEntreprise>
  effectif_ent: Record<DataHash, EntréeEffectif>
  bdf: Record<DataHash, EntréeBdf>
  diane: Record<DataHash, EntréeDiane>
}>

export type DataType = Exclude<keyof BatchValue, "compact"> // => 'reporder' | 'effectif' | 'apconso' | ...
