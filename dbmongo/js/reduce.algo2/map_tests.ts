import test, {ExecutionContext}  from 'ava'
import { map } from './map.js'

test('map is a function', (t: ExecutionContext) => {
  t.is(typeof map, 'function')
})
