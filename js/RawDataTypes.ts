/* eslint-disable no-use-before-define */

import {
  EntréeApConso,
  EntréeApDemande,
  EntréeBdf,
  EntréeCcsf,
  EntréeCompte,
  EntréeCotisation,
  EntréeDebit,
  EntréeDelai,
  EntréeDéfaillances,
  EntréeDiane,
  EntréeEffectif,
  EntréeEffectifEnt,
  EntréeEllisphere,
  EntréePaydex,
  EntréeSirene,
  EntréeSireneEntreprise,
} from "./GeneratedTypes"

// Types de données de base

export type Timestamp = number // Date.getTime()
export type Periode = Timestamp

/**
 * Cette classe encapsule un Map<Timestamp, T>, pour valider (et
 * convertir, si besoin) la période passée aux méthodes get() et set().
 */
export class ParPériode<T> {
  private map = new Map<Timestamp, T>()
  private getNumericValue(période: Date | Timestamp | string): number {
    if (typeof période === "number") return période
    if (typeof période === "string") return parseInt(période)
    if (période instanceof Date) return période.getTime()
    throw new TypeError("type non supporté: " + typeof période)
  }
  // pour vérifier que le timestamp retourné par getNumericValue est valide
  private getTimestamp(période: Date | Timestamp | string): Timestamp {
    const timestamp = this.getNumericValue(période)
    if (isNaN(timestamp) || new Date(timestamp).getTime() !== timestamp) {
      throw new RangeError("valeur invalide: " + période)
    }
    return timestamp
  }
  constructor(entries?: readonly (readonly [number, T])[] | null | undefined) {
    this.map = new Map<Timestamp, T>(entries)
  }
  /**
   * Informe sur la présence d'une valeur associée à la période donnée.
   * @throws TypeError si la période n'est pas valide.
   */
  has(période: Date | Timestamp | string): boolean {
    return this.map.has(this.getTimestamp(période))
  }
  /**
   * Retourne la valeur associée à la période donnée.
   * @throws TypeError si la période n'est pas valide.
   */
  get(période: Date | Timestamp | string): T | undefined {
    return this.map.get(this.getTimestamp(période))
  }
  /**
   * Définit la valeur associée à la période donnée.
   * @throws TypeError si la période n'est pas valide.
   */
  set(période: Date | Timestamp | string, val: T): this {
    const timestamp = this.getTimestamp(période)
    this.map.set(timestamp, val)
    return this
  }

  keys(): IterableIterator<Timestamp> {
    return this.map.keys()
  }

  values(): IterableIterator<T> {
    return this.map.values()
  }

  entries(): IterableIterator<[Timestamp, T]> {
    return this.map.entries()
  }

  forEach(
    callbackfn: (value: T, key: number, map: Map<Timestamp, T>) => void,
    thisArg?: unknown
  ): void {
    return this.map.forEach(callbackfn, thisArg)
  }

}

export type Departement = string

export type Siren = string
export type Siret = string
export type SiretOrSiren = Siret | Siren // TODO: supprimer ce type, une fois que tous les champs auront été séparés
export type CodeAPE = string

export type DataHash = string
export type ParHash<T> = Record<DataHash, T>

// Données importées pour une entreprise ou établissement

export type EntrepriseDataValues = {
  key: Siren
  scope: "entreprise"
  batch: Record<BatchKey, Partial<BatchValueProps>> // TODO: remplacer par `Partial<EntrepriseBatchProps>>`, une fois que tous les champs auront été séparés
}

export type EtablissementDataValues = {
  key: Siret
  scope: "etablissement"
  batch: Record<BatchKey, Partial<BatchValueProps>> // TODO: remplacer par Partial<EtablissementBatchProps>, une fois que tous les champs auront été séparés
}

export type Scope = (EntrepriseDataValues | EtablissementDataValues)["scope"]

export type CompanyDataValues = EntrepriseDataValues | EtablissementDataValues

export type CompanyDataValuesWithFlags = CompanyDataValues & IndexFlags

export type IndexFlags = {
  index: {
    algo2: boolean // pour spécifier quelles données seront à calculer puis inclure dans Features, par Reduce.algo2
  }
}

// Données importées par les parseurs, pour chaque source de données

export type BatchKey = string

export type BatchValues = Record<BatchKey, BatchValue>

export type DataType = keyof BatchValueProps // => 'reporder' | 'effectif' | 'apconso' | ...

export type BatchValue = Partial<BatchValueProps>

type CommonBatchProps = {
  reporder: ParHash<EntréeRepOrder> // RepOrder est généré par "compact", et non importé => Usage de Periode en guise de hash d'indexation
}

export type EntrepriseBatchProps = CommonBatchProps & {
  paydex: ParHash<EntréePaydex>
}

export type EtablissementBatchProps = CommonBatchProps & {
  apconso: ParHash<EntréeApConso>
}

// TODO: continuer d'extraire les propriétés vers EntrepriseBatchProps et EtablissementBatchProps, puis supprimer BatchValueProps et les types qui en dépendent
export type BatchValueProps = CommonBatchProps &
  EntrepriseBatchProps &
  EtablissementBatchProps & {
    effectif: ParHash<EntréeEffectif>
    apdemande: ParHash<EntréeApDemande>
    compte: ParHash<EntréeCompte>
    delai: ParHash<EntréeDelai>
    procol: ParHash<EntréeDéfaillances>
    cotisation: ParHash<EntréeCotisation>
    debit: ParHash<EntréeDebit>
    ccsf: ParHash<EntréeCcsf>
    sirene: ParHash<EntréeSirene>
    sirene_ul: ParHash<EntréeSireneEntreprise>
    effectif_ent: ParHash<EntréeEffectifEnt>
    bdf: ParHash<EntréeBdf>
    diane: ParHash<EntréeDiane>
    ellisphere: ParHash<EntréeEllisphere>
  }

// Détail des types de données

export type EntréeRepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}
