// Types partagés

type Periode = string // Date.toString()
type Timestamp = number // Date.getTime()
type ParPériode<T> = Record<Periode, T>

type Departement = string

type Siret = string
type SiretOrSiren = Siret | string
type CodeAPE = string

type DataHash = string
type ParHash<T> = Record<DataHash, T>
