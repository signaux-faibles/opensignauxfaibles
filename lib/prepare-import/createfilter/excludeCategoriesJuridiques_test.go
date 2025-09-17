package createfilter

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestExclusionList(t *testing.T) {
	ass := assert.New(t)

	// GIVEN
	sireneULPath := "./test_uniteLegale.csv"

	// WHEN
	excludedSirens := readExcludedSirens(sireneULPath)
	_, ok1 := excludedSirens["111111111"]
	_, ok2 := excludedSirens["222222222"]
	_, ok3 := excludedSirens["333333333"]
	_, ok4 := excludedSirens["444444444"]

	// THEN
	ass.False(ok1)
	ass.True(ok2)
	ass.True(ok3)
	ass.False(ok4)
}

func TestCategorieJuridiqueFilter(t *testing.T) {
	ass := assert.New(t)

	// GIVEN
	sireneULPath := "./test_uniteLegale.csv"
	testFilter := CategorieJuridiqueFilter(sireneULPath)
	initialPerimeter := map[string]struct{}{
		"111111111": {},
		"222222222": {},
		"333333333": {},
		"444444444": {},
	}
	expectedPerimeter := map[string]struct{}{
		"111111111": {},
		"444444444": {},
	}

	// WHEN
	actualPerimeter := applyFilter(initialPerimeter, testFilter)

	// THEN
	eq := reflect.DeepEqual(actualPerimeter, expectedPerimeter)
	ass.True(eq)
}
