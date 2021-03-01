import { Timestamp } from "../RawDataTypes"

export interface ParPériode<T> extends Map<Timestamp, T> {
  has(période: Date | Timestamp | string): boolean
  get(période: Date | Timestamp | string): T | undefined
  set(période: Date | Timestamp | string, val: T): this
}

export function newParPériode<T>(
  arg?: readonly (readonly [number, T])[] | null | undefined
): ParPériode<T> {
  // if ("_get" in Map.prototype && "put" in Map.prototype) {
  let data: Record<number, T> = {}

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
    // @ ts-expect-error Override MongoDB's Map implementation
    has(key: number) {
      return key in data
    }
    // @ ts-expect-error Override MongoDB's Map implementation
    get(key: number): T | undefined {
      return data[key]
    }
    // @ ts-expect-error Override MongoDB's Map implementation
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
    // @ ts-expect-error Override MongoDB's Map implementation
    *keys() {
      for (const k in data) {
        yield parseInt(k)
      }
    }
    // @ ts-expect-error Override MongoDB's Map implementation
    *values() {
      for (const val of Object.values(data)) {
        yield val
      }
    }
    // @ ts-expect-error Override MongoDB's Map implementation
    *entries(): Generator<[number, T]> {
      for (const k in data) {
        yield [parseInt(k), data[k] as T]
      }
    }
    // @ ts-expect-error Override MongoDB's Map implementation
    forEach(
      callbackfn: (value: T, key: number, map: any) => void,
      thisArg?: unknown
    ): void {
      // @ ts-expect-error entries() is defined above
      for (const [key, value] of this.entries()) {
        callbackfn.call(thisArg, value as T, key, this)
      }
    }
    // }
    [Symbol.iterator]() {
      return this.entries()
    }

    get [Symbol.toStringTag]() {
      return "Map"
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
