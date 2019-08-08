package jstests

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

func mapReduceTestHelper(t *testing.T, filenames ...string) string {
	var test_file string
	filenames = append([]string{"testing.js"}, filenames...)
	for _, file := range filenames {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err.Error())
		}
		fc := strings.Replace(string(fileContent), "f.", "", -1)
		test_file = test_file + fc
	}
	return test_file
}

func testResultsHelper(t *testing.T, vm *otto.Otto) {

	aux_res, err := vm.Get("test_results")
	res, _ := aux_res.Export()
	res_bool := res.([]bool)
	if err != nil {
		t.Fatal("Failed to compute test results: " + err.Error())
	}
	for ind, test := range res_bool {
		if !test {
			t.Fatal("A test failed: " + strconv.Itoa(ind))
		}
	}
}

func Test_lookAhead(t *testing.T) {
	vm := otto.New()
	_, err := vm.Run(mapReduceTestHelper(t, "lookAhead.js", "lookAhead_test.js"))

	if err != nil {
		t.Fatal("Is the script ECMA5 compatible ? ECMA6 features are not compatible. " + err.Error())
	}
	testResultsHelper(t, vm)
}

func Test_cibleApprentissage(t *testing.T) {
	vm := otto.New()
	_, err := vm.Run(mapReduceTestHelper(t, "lookAhead.js", "cibleApprentissage.js", "cibleApprentissage_test.js"))
	if err != nil {
		t.Fatal("Is the script ECMA5 compatible ? ECMA6 features are not compatible. " + err.Error())
	}
	testResultsHelper(t, vm)
}

func Test_add(t *testing.T) {
	vm := otto.New()
	_, err := vm.Run(mapReduceTestHelper(t, "add.js", "add_test.js"))
	if err != nil {
		t.Fatal("Is the script ECMA5 compatible ? ECMA6 features are not compatible. " + err.Error())
	}
	testResultsHelper(t, vm)
}
