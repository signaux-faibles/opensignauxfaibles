import test from 'ava'
import { map } from './map.js'

test('one plus two equals three', t => {
  t.is(1 + 2, 3)
})

test('map is a function', t => {
  t.is(typeof map, 'function')
})
