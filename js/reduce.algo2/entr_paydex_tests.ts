import test from "ava"
import { entr_paydex } from "./entr_paydex"

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
