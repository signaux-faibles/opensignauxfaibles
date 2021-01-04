import test from "ava"
import { entr_paydex } from "./entr_paydex"
import { f } from "./functions"

test(`doit retourner paydex_nb_jours pour chaque période`, (t) => {
  const dateDébut = new Date("2015-12-01T00:00Z")
  const dateFin = new Date("2016-01-01T00:00Z")
  const périodes = [dateDébut, dateFin]
  const entréesPaydex = {
    decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
    janvier: { date_valeur: new Date("2016-01-15T00:00Z"), nb_jours: 2 },
  }
  const sortiePaydex = entr_paydex(entréesPaydex, périodes)
  t.is(
    Object.values(sortiePaydex).length,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  Object.values(entréesPaydex).forEach((entréePaydex, i) =>
    t.is(
      sortiePaydex[(périodes[i] as Date).getTime()]?.paydex_nb_jours,
      entréePaydex.nb_jours,
      "le nombre de jours paydex doit être transmis pour chaque période"
    )
  )
})

test(`doit donner accès au nombre de jours de la période précédente`, (t) => {
  const dateDébut = new Date("2015-12-01T00:00Z")
  const dateFin = new Date("2016-01-01T00:00Z")
  const périodes = [dateDébut, dateFin]
  const entréesPaydex = {
    decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
  }
  const sortiePaydex = entr_paydex(entréesPaydex, périodes)
  t.is(
    Object.values(sortiePaydex).length,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  const [decembre, janvier] = Object.values(sortiePaydex)
  t.is(janvier?.paydex_nb_jours_past_1, entréesPaydex.decembre.nb_jours)
  t.is(decembre?.paydex_nb_jours_past_1, null)
})

test(`doit donner accès au nombre de jours d'il y a 12 mois`, (t) => {
  const dateDébut = new Date("2015-12-01T00:00Z")
  const dateFin = new Date("2017-01-01T00:00Z")
  const périodes = f.generatePeriodSerie(dateDébut, dateFin)
  const entréesPaydex = {
    decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
  }
  const sortiePaydex = entr_paydex(entréesPaydex, périodes)
  t.is(
    Object.values(sortiePaydex).length,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  const [dernièrePériode, ...autres] = Object.values(sortiePaydex).reverse()
  t.is(
    dernièrePériode?.paydex_nb_jours_past_12,
    entréesPaydex.decembre.nb_jours
  )
  autres.forEach((autrePériode) =>
    t.is(autrePériode?.paydex_nb_jours_past_12, null)
  )
})
