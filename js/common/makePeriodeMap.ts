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
  // if ("_get" in Map.prototype && "put" in Map.prototype) {
  let data: Record<number, T> = {}

  /**
   * MyMap est une ré-implémentation partielle de la classe Map standard de
   * JavaScript, utilisant un objet JavaScript pour indexer les entrées.
   * Implémenter l'interface de Map permet de valider les dates passées
   * en tant que clés, et de supporter plusieurs représentations de ces dates
   * (ex: instance Date, timestamp numérique ou chaine de caractères), tout en
   * evitant que des chaines de caractères arbitaires y soient passées.
   */
  class MyMap {
    constructor(
      entries?: readonly (readonly [number, T])[] | null | undefined
    ) {
      if (entries) {
        for (const [key, value] of entries) {
          data[key] = value
        }
      }
    }
    has(key: number) {
      return key in data
    }
    get(key: number): T | undefined {
      return data[key]
    }
    set(key: number, value: T): this {
      data[key] = value
      return this
    }
    get size() {
      return Object.keys(data).length
    }
    clear() {
      data = {}
    }
    delete(key: number): boolean {
      const exists = key in data
      delete data[key]
      return exists
    }
    *keys() {
      for (const k in data) {
        yield parseInt(k)
      }
    }
    *values() {
      for (const val of Object.values(data)) {
        yield val
      }
    }
    *entries(): Generator<[number, T]> {
      for (const k in data) {
        yield [parseInt(k), data[k] as T]
      }
    }
    forEach(
      callbackfn: (value: T, key: number, map: any) => void,
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
      return "MyMap"
    }
  }

  /**
   * Cette classe est une Map<Timestamp, T> qui valide (et convertit,
   * si besoin) la période passée aux différentes méthodes.
   */
  class ParPériodeImpl
    extends MyMap /*<Timestamp, T>*/
    implements ParPériode<T> {
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
    /**
     * Informe sur la présence d'une valeur associée à la période donnée.
     * @throws TypeError si la période n'est pas valide.
     */
    has(période: Date | Timestamp | string): boolean {
      return super.has(this.getTimestamp(période))
    }
    /**
     * Retourne la valeur associée à la période donnée.
     * @throws TypeError si la période n'est pas valide.
     */
    get(période: Date | Timestamp | string): T | undefined {
      return super.get(this.getTimestamp(période))
    }
    /**
     * Définit la valeur associée à la période donnée.
     * @throws TypeError si la période n'est pas valide.
     */
    set(période: Date | Timestamp | string, val: T): this {
      const timestamp = this.getTimestamp(période)
      super.set(timestamp, val)
      return this
    }
  }

  return new ParPériodeImpl(arg)
}
