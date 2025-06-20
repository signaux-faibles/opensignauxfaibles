import test from "ava"
import { fraisFinancier, ChampsDiane } from "./fraisFinancier"
import { EntréeDiane } from "../GeneratedTypes"

const fakeDiane = () => ({
  interets: 50,
  excedent_brut_d_exploitation: 1000,
  produits_financiers: 100,
  produit_exceptionnel: 120,
  charge_exceptionnelle: 450,
  charges_financieres: 160,
})

test(`fraisFinancier est calculé selon la formule:
 interets / (
  excedent_brut_d_exploitation +
  produits_financiers +
  produit_exceptionnel -
  charge_exceptionnelle -
  charges_financieres ) * 100`, (t: test) => {
  const diane = fakeDiane()
  const resultat =
    (diane.interets /
      (diane.excedent_brut_d_exploitation +
        diane.produits_financiers +
        diane.produit_exceptionnel -
        diane.charge_exceptionnelle -
        diane.charges_financieres)) *
    100
  t.is(fraisFinancier(diane as EntréeDiane), resultat)
})

const proprietes: (keyof ChampsDiane)[] = [
  "interets",
  "excedent_brut_d_exploitation",
  "produits_financiers",
  "produit_exceptionnel",
  "charge_exceptionnelle",
  "charges_financieres",
]
proprietes.forEach((propriete) =>
  test(`fraisFinancier est nul si "${propriete}" n'est pas disponible dans Diane`, (t: test) => {
    const diane = fakeDiane()
    delete diane[propriete]
    t.is(fraisFinancier(diane), null)
  })
)

// $ cd js && $(npm bin)/ava ./reduce.algo2/fraisFinancier_tests.ts
