import { effectifs, SortieEffectifs } from "./effectifs"
import test, { ExecutionContext } from "ava"
import { setGlobals } from "../test/helpers/setGlobals"
import { EntréeEffectif } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"
import { ParPériode } from "../common/makePeriodeMap"

type SortieEffectifsEtab = SortieEffectifs<"effectif">

function assertEffectif(
  t: ExecutionContext,
  résultat: ParPériode<SortieEffectifsEtab>,
  effectifsAttendus: Array<[number | null, boolean]>
): void {
  const périodes = [...résultat.keys()]
  for (let i = 0; i < périodes.length; i++) {
    const période = périodes[i] as number
    t.is(
      résultat.get(période)?.effectif,
      effectifsAttendus[i]?.[0],
      `valeur inattendue pour la période ${i}`
    )
    t.is(
      résultat.get(période)?.effectif_reporte,
      effectifsAttendus[i]?.[1] ? 1 : 0,
      `flag de report inattendu pour la période ${i}`
    )
  }
}

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
    const résultat = effectifs(
      entréeEffectif as ParHash<EntréeEffectif>,
      periodes.map((d) => d.getTime()),
      "effectif" // clé: établissement
    )
    assertEffectif(t, résultat, [
      [entréeEffectif["hash_periode0"].effectif, false],
      [entréeEffectif["hash_periode0"].effectif, true],
    ])
  }
)

test.serial(
  "Effectif ne reporte pas de valeur si le nombre de mois avec effectifs manquants est strictement supérieur au nombre de mois avec effectif manquant attendus (offset_effectif)",
  (t: ExecutionContext) => {
    setGlobals({ offset_effectif: -2 })
    const periodes = [
      new Date("2020-01-01"),
      new Date("2020-02-01"),
      new Date("2020-03-01"),
    ]
    const entréeEffectif = {
      hash_periode0: {
        effectif: 24,
        periode: periodes[0],
        numero_compte: "123",
      },
    }
    const résultat = effectifs(
      entréeEffectif as ParHash<EntréeEffectif>,
      periodes.map((d) => d.getTime()),
      "effectif" // clé: établissement
    )
    assertEffectif(t, résultat, [
      [entréeEffectif["hash_periode0"].effectif, false],
      [null, true], // TODO: reporté devrait être false
      [null, true],
    ])
  }
)

test.serial(
  "Effectif reporte la dernière valeur connue si le nombre de mois avec effectifs manquants est égal au nombre de mois avec effectif manquant attendu (offset_effectif)",
  (t: ExecutionContext) => {
    setGlobals({ offset_effectif: -3 }) // car 2 mois inconnus
    const periodes = [
      new Date("2020-01-01"),
      new Date("2020-02-01"),
      new Date("2020-03-01"),
      new Date("2020-04-01"),
    ]
    const entréeEffectif = {
      hash_periode0: {
        effectif: 24,
        periode: periodes[0],
        numero_compte: "123",
      },
      hash_periode1: {
        effectif: 25,
        periode: periodes[1],
        numero_compte: "123",
      },
    }
    const résultat = effectifs(
      entréeEffectif as ParHash<EntréeEffectif>,
      periodes.map((d) => d.getTime()),
      "effectif" // clé: établissement
    )
    assertEffectif(t, résultat, [
      [entréeEffectif["hash_periode0"].effectif, false],
      [entréeEffectif["hash_periode1"].effectif, false],
      [entréeEffectif["hash_periode1"].effectif, true],
      [entréeEffectif["hash_periode1"].effectif, true],
    ])
  }
)

// TODO:
// * Tester la clé "effectif_ent"
// * Tester les valeurs past (similairement à ce que fait débit)
