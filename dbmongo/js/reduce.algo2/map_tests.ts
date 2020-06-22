import test, { ExecutionContext } from "ava"
import { map } from "./map"

test("map is a function", (t: ExecutionContext) => {
  t.is(typeof map, "function")
})
