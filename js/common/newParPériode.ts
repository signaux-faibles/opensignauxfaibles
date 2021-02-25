import { Timestamp } from "../RawDataTypes"

export interface ParPériode<T> extends Map<Timestamp, T> {
  has(période: Date | Timestamp | string): boolean
  get(période: Date | Timestamp | string): T | undefined
  set(période: Date | Timestamp | string, val: T): this
}

declare function friendlyEqual(): void

export function newParPériode<T>(
  arg?: readonly (readonly [number, T])[] | null | undefined
): ParPériode<T> {
  // @ts-expect-error To prevent "ReferenceError: friendlyEqual is not defined" errors, see https://jira.mongodb.org/browse/SERVER-19169 and https://github.com/mongodb/mongo/blob/master/src/mongo/shell/types.js#L584
  friendlyEqual = () => {} // eslint-disable-line

  Map.prototype.keys =
    Map.prototype.keys ||
    function () {
      const keys = []
      // @ts-expect-error Polyfill for Map.prototype.keys on MongoDB's Map implementation
      for (const k in this._data) {
        keys.push(k)
      }
      return keys
    }

  /**
   * Cette classe est une Map<Timestamp, T> qui valide (et convertit,
   * si besoin) la période passée aux différentes méthodes.
   */
  class ParPériodeImpl<T> extends Map<Timestamp, T> implements ParPériode<T> {
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
      "put" in Map.prototype
        ? // @ts-expect-error MongoDB provides Map.prototype.put() instead of Map.prototype.set()
          super.put(timestamp, val)
        : super.set(timestamp, val)
      return this
    }
  }

  return new ParPériodeImpl(arg)
}
