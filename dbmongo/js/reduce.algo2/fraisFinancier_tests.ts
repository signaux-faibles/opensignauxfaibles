import test from 'ava'
import { fraisFinancier } from './fraisFinancier.js'

test(`fraisFinancier est calculé selon la formule:
 interets / (
  excedent_brut_d_exploitation +
  produits_financiers +
  produit_exceptionnel -
  charge_exceptionnelle -
  charges_financieres ) * 100`, t => {
  const diane = {
    interets: 50,
    excedent_brut_d_exploitation: 1000,
    produits_financiers: 100,
    produit_exceptionnel: 120,
    charge_exceptionnelle: 450,
    charges_financieres: 160
  }
  const resultat = diane.interets / (diane.excedent_brut_d_exploitation +
    diane.produits_financiers + diane.produit_exceptionnel -
    diane.charge_exceptionnelle - diane.charges_financieres) * 100
  t.is(fraisFinancier(diane), resultat)
})

const proprietes = ["interets", "excedent_brut_d_exploitation", "produits_financiers", "produit_exceptionnel", "charge_exceptionnelle", "charges_financieres"]
proprietes.forEach((propriete) =>
  test(`fraisFinancier est nul si "${propriete}" n'est pas disponible dans Diane`, t => {
    const diane = {
      interets: 50,
      excedent_brut_d_exploitation: 1000,
      produits_financiers: 100,
      produit_exceptionnel: 120,
      charge_exceptionnelle: 450,
      charges_financieres: 160
    }
    diane[propriete] = null
    t.is(fraisFinancier(diane), null)
  })
)

// $ node_modules/.bin/ava-ts ./reduce.algo2/fraisFinancier_tests.ts
