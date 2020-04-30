var test_cases = [
  {data : ["nom_entreprise", null, null, null, null, null, null],
    expected : "nom_entreprise"},
  {data : ["nom_entreprise", "roger", null, null, null, null, null],
    expected : "nom_entreprise"},
  {data : [null, "roger", null, "pierre", "henry", null, null],
    expected : "roger*pierre henry/"},
  {data : [null, "roger", undefined, "pierre", "henry", undefined, undefined],
    expected : "roger*pierre henry/"},
  {data : [null, "toto", "titi", "mathilde", "louisette", "fanny", "géraldine"],
    expected : "toto*titi/mathilde louisette fanny géraldine/"}
]

test_results = test_cases.map(tc => {
  return(raison_sociale(tc.data[0], tc.data[1], tc.data[2], tc.data[3], tc.data[4], tc.data[5], tc.data[6]) == tc.expected)
}
)
print(test_results.every(t=>t))
