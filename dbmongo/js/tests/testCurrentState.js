var actual_res = currentState([
  { "apdemande": ["a", "b"], "apconso": ["c"]},
  { "apdemande": ["d"], "compact": {"delete" : {"apdemande":["a", "b"]}}}
])

var expected_res = {
  "apdemande" : new Set(["d"]),
  "apconso" : new Set(["c"])
}

function eqSet(as, bs) {
    if (as.size !== bs.size) return false;
    for (var a of as) if (!bs.has(a)) return false;
    return true;
}

if (!eqSet(actual_res["apdemande"], expected_res["apdemande"]) ||
    !eqSet(actual_res["apconso"], expected_res["apconso"])) {
  print("testCurrentState failed")
  print("actual:")
  print(JSON.stringify([...actual_res["apdemande"]], null, 2))
  print("expected:")
  print(JSON.stringify([...expected_res["apdemande"]], null, 2))
} else { 
  print("testCurrentState run successfully")
}
