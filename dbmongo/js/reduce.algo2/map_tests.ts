import test from "ava"
import { map } from "./map.js"

test("map is a function", (t) => {
  t.is(typeof map, "function")
})
