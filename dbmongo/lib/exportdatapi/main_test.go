package exportdatapi

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// func Test_GetRegions(t *testing.T) {
// 	spew.Dump(GetRegions("1805"))
// }

func Test_GetPolicies(t *testing.T) {
	spew.Dump(GetPolicies("1805"))
}
