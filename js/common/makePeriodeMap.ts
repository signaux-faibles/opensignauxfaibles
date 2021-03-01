import { Timestamp } from "../RawDataTypes"

/**
 * ParPériode est une extension de Map permettant de valider les dates fournies
 * avant de les sérialiser sous forme de timestamp (numérique), type employé
 * pour l'indexation des données par période dans ce Map.
 */
export interface ParPériode<T> extends Map<Timestamp, T> {
  has(période: Date | Timestamp | string): boolean
  get(période: Date | Timestamp | string): T | undefined
  set(période: Date | Timestamp | string, val: T): this
}

/**
 * makePeriodeMap() retourne une nouvelle instance de la classe ParPériode
 * (équivalente à Map<Timestamp, T>). Cette fonction a été fournie à défaut
 * d'être parvenu à inclure directement la classe ParPériode dans le scope
 * transmis à MongoDB depuis le traitement map-reduce lancé par le code Go.
 * @param arg (optionnel) - pour initialiser la Map avec un tableau d'entries.
 */
export function makePeriodeMap<T>(
  arg?: readonly (readonly [number, T])[] | null | undefined
): ParPériode<T> {
  /**
   * IntMap est une ré-implémentation partielle de Map<number, T> utilisant un
   * objet JavaScript pour indexer les entrées, et rendue nécéssaire par le
   * fait que MongoDB fournit une version non standard de la classe Map.
   */
  class IntMap {
    private data: Record<number, T> = {}
    constructor(
      entries?: readonly (readonly [number, T])[] | null | undefined
    ) {
      if (entries) {
        for (const [key, value] of entries) {
          this.data[key] = value
        }
      }
    }
    has(key: number) {
      return key in this.data
    }
    get(key: number): T | undefined {
      return this.data[key]
    }
    set(key: number, value: T): this {
      this.data[key] = value
      return this
    }
    get size() {
      return Object.keys(this.data).length
    }
    clear() {
      this.data = {}
    }
    delete(key: number): boolean {
      const exists = key in this.data
      delete this.data[key]
      return exists
    }
    *keys() {
      for (const k in this.data) {
        yield parseInt(k)
      }
    }
    *values() {
      for (const val of Object.values(this.data)) {
        yield val
      }
    }
    *entries(): Generator<[number, T]> {
      for (const k in this.data) {
        yield [parseInt(k), this.data[k] as T]
      }
    }
    forEach(
      callbackfn: (value: T, key: number, map: this) => void,
      thisArg?: unknown
    ): void {
      for (const [key, value] of this.entries()) {
        callbackfn.call(thisArg, value as T, key, this)
      }
    }
    [Symbol.iterator]() {
      return this.entries()
    }
    get [Symbol.toStringTag]() {
      return "IntMap"
    }
  }

  /**
   * Cette classe étend Map<Timestamp, T> pour valider les dates passées
   * en tant que clés et supporter diverses représentations de ces dates
   * (ex: instance Date, timestamp numérique ou chaine de caractères), tout en
   * evitant que des chaines de caractères arbitaires y soient passées.
   */
  class ParPériodeImpl extends IntMap implements ParPériode<T> {
    /** Extraie le timestamp d'une date, quelque soit sa représentation. */
    private getNumericValue(période: Date | Timestamp | string): number {
      if (typeof période === "number") return période
      if (typeof période === "string") return parseInt(période)
      if (période instanceof Date) return période.getTime()
      throw new TypeError("type non supporté: " + typeof période)
    }
    /** Vérifie que le timestamp retourné par getNumericValue est valide. */
    private getTimestamp(période: Date | Timestamp | string): Timestamp {
      const timestamp = this.getNumericValue(période)
      if (isNaN(timestamp) || new Date(timestamp).getTime() !== timestamp) {
        throw new RangeError("valeur invalide: " + période)
      }
      return timestamp
    }
    /** @throws TypeError ou RangeError si la période n'est pas valide. */
    has(période: Date | Timestamp | string): boolean {
      return super.has(this.getTimestamp(période))
    }
    /** @throws TypeError ou RangeError si la période n'est pas valide. */
    get(période: Date | Timestamp | string): T | undefined {
      return super.get(this.getTimestamp(période))
    }
    /** @throws TypeError ou RangeError si la période n'est pas valide. */
    set(période: Date | Timestamp | string, val: T): this {
      const timestamp = this.getTimestamp(période)
      super.set(timestamp, val)
      return this
    }
  }

  return new ParPériodeImpl(arg)
}
