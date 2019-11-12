package exportdatapi

import (
	"fmt"
	"testing"
)

func Test_GetRegions(t *testing.T) {
	var regions []string
	for r := range reverseMap(urssafMapping) {
		regions = append(regions, r)
	}
	fmt.Println(regions)
}
