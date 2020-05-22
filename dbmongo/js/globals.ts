// Déclaration des fonctions globales fournies par MongoDB
declare function emit(key: any, value: any): void
declare function print(...any): void

// Déclaration des fonctions globales fournies par JSC
declare function debug(string) // supported by jsc, to print in stdout

// Déclaration de variables globales
// eslint-disable-next-line @typescript-eslint/no-unused-vars
let f: {
  [key: string]: Function
}
