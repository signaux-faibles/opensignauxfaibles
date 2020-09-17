package marshal

import (
	"errors"
	"strings"

	"github.com/globalsign/mgo/bson"
)

//CheckMarshallingMap checks with the header row if the marshalling map is
//correct
func CheckMarshallingMap(headerRow []string, marshallingMap map[string]int) error {

	errorString := "Following fields do not match the specification:"
	var failingFields []string

	for k, v := range marshallingMap {
		ok := (headerRow[v] == k)
		if !ok {
			failingFields = append(failingFields, k)
		}
	}
	if len(failingFields) == 0 {
		return nil
	}
	return errors.New(errorString + strings.Join(failingFields, ", "))
}

// Object ...
type Object struct {
	Data     bson.M
	key      string
	scope    string
	datatype string
}

// Scope ...
func (obj Object) Scope() string {
	return obj.scope
}

// Key ...
func (obj Object) Key() string {
	return obj.key
}

// Type ...
func (obj Object) Type() string {
	return obj.datatype
}
