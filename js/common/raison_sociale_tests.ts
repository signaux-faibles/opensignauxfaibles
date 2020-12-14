import test, { ExecutionContext } from "ava"
import { raison_sociale } from "./raison_sociale"

type Param = string | null | undefined

type TestCase = {
  data: [Param, Param, Param, Param, Param, Param, Param]
  expected: string
}

const testCases: TestCase[] = [
  {
    data: ["nom_entreprise", null, null, null, null, null, null],
    expected: "nom_entreprise",
  },
  {
    data: ["nom_entreprise", "roger", null, null, null, null, null],
    expected: "nom_entreprise",
  },
  {
    data: [null, "roger", null, "pierre", "henry", null, null],
    expected: "roger*pierre henry/",
  },
  {
    data: [null, "roger", undefined, "pierre", "henry", undefined, undefined],
    expected: "roger*pierre henry/",
  },
  {
    data: [null, "toto", "titi", "mathilde", "louisette", "fanny", "géraldine"],
    expected: "toto*titi/mathilde louisette fanny géraldine/",
  },
]

testCases.forEach(({ data, expected }) => {
  test.serial(
    `raison_sociale(${data.map((param) =>
      param ? param.toString() : typeof param
    )}) === ${expected}`,
    (t: ExecutionContext) => {
      const actualResults = raison_sociale(...data)
      t.deepEqual(actualResults, expected)
    }
  )
})
