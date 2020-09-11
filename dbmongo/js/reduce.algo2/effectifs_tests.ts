import { effectifs } from "./effectifs"
import test, { ExecutionContext } from "ava"
import { setGlobals } from "../test/helpers/setGlobals"

test.serial(
  "Effectif reporte la valeur du mois m à la période m + 1, si offset_effectif vaut -2",
  (t: ExecutionContext) => {
    setGlobals({ offset_effectif: -2 }) // TODO: offset_effectif ne devrait pas être négatif. et devrait être égale au nombre de mois manquants de l'effectif => dans ce test, on devrait avoir 1. => à redefinir dans reduce.go
    const periodes = [new Date("2020-01-01"), new Date("2020-02-01")]
    const entréeEffectif = {
      hash_periode0: {
        effectif: 24,
        periode: periodes[0],
        numero_compte: "123",
      },
    }
    const clé = "effectif"
    const résultat = effectifs(
      entréeEffectif,
      periodes.map((d) => d.getTime()),
      clé
    )
    t.log(résultat)
    t.deepEqual(
      résultat[periodes[1].getTime().toString()].effectif,
      entréeEffectif["hash_periode0"].effectif
    )
  }
)
