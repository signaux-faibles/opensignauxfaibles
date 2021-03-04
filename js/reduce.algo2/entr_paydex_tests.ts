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
    sortiePaydex.size,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  Object.values(entréesPaydex).forEach((entréePaydex, i) =>
    t.is(
      sortiePaydex.get(périodes[i] as Date)?.paydex_nb_jours,
      entréePaydex.nb_jours,
      "le nombre de jours paydex doit être transmis pour chaque période"
    )
  )
})

test(`doit donner accès au nombre de jours d'il y a 3 mois`, (t) => {
  const dateDébut = new Date("2015-12-01T00:00Z")
  const dateFin = new Date("2016-03-01T00:00Z")
  const périodes = [dateDébut, dateFin]
  const entréesPaydex = {
    decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
  }
  const sortiePaydex = entr_paydex(entréesPaydex, périodes)
  t.is(
    sortiePaydex.size,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  const [decembre, mars] = sortiePaydex.values()
  t.is(mars?.paydex_nb_jours_past_3, entréesPaydex.decembre.nb_jours)
  t.is(decembre?.paydex_nb_jours_past_3, null)
})

test(`doit donner accès au nombre de jours d'il y a 6 mois`, (t) => {
  const dateDébut = new Date("2015-12-01T00:00Z")
  const dateFin = new Date("2016-06-01T00:00Z")
  const périodes = [dateDébut, dateFin]
  const entréesPaydex = {
    decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
  }
  const sortiePaydex = entr_paydex(entréesPaydex, périodes)
  t.is(
    sortiePaydex.size,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  const [decembre, juin] = sortiePaydex.values()
  t.is(juin?.paydex_nb_jours_past_6, entréesPaydex.decembre.nb_jours)
  t.is(decembre?.paydex_nb_jours_past_6, null)
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
    sortiePaydex.size,
    périodes.length,
    "entr_paydex() doit émettre un objet par période"
  )
  const [dernièrePériode, ...autres] = [...sortiePaydex.values()].reverse()
  t.is(
    dernièrePériode?.paydex_nb_jours_past_12,
    entréesPaydex.decembre.nb_jours
  )
  autres.forEach((autrePériode) =>
    t.is(autrePériode?.paydex_nb_jours_past_12, null)
  )
})
