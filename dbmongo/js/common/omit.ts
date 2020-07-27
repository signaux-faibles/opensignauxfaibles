// Fonction pour omettre des props, tout en retournant le bon type
export function omit<Source, Exclusions extends Array<keyof Source>>(
  object: Source,
  ...propNames: Exclusions
): Omit<Source, Exclusions[number]> {
  const result = Object.assign({}, object)
  for (const prop of propNames) {
    delete result[prop]
  }
  return result as Omit<Source, Exclusions[number]>
}
