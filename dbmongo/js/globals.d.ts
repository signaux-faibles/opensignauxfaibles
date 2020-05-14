// Déclaration de variables globales accessible via l'objet globalThis.
interface Global {
  f: {
    [key: string]: Function
  }
}
