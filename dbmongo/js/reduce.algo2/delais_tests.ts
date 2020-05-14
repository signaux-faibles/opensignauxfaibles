import test from "ava"
import { delais } from "./delais"

test("delais est défini", (t) => {
  t.is(typeof delais, "function")
})
