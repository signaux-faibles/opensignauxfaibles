import test from 'ava'
import { fraisFinancier } from './fraisFinancier.js'

test(`fraisFinancier est nul si "interets" n'est pas disponible dans Diane`, t => {
  const diane = {}
  t.is(fraisFinancier(diane), null)
})
