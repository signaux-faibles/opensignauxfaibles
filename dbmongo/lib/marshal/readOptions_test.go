package marshal

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cnf/structhash"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestResolveParserOptions(t *testing.T) {
	testParserOptionsDir := filepath.Join("testData", "testParserOptions")

	_, err := RegisteredParserOptions(testParserOptionsDir)
	if err != nil {
		t.Error(err)
	}

}

func compareParserOptions(t *testing.T, testFile string, goldenFile string, testDescription string) {

	actual := &ParserOptions{}
	// fmt.Println(actual)
	err := actual.ReadOptions(filepath.Join("testData", testFile))
	t.Log(actual)
	fullGoldenFile := filepath.Join("testData", goldenFile)
	if err != nil {
		t.Log(filepath.Join("testData", testFile))
		t.Error("test file could not be read " + err.Error())
	}
	if *update {
		ioutil.WriteFile(fullGoldenFile, structhash.Md5(actual, 1), 0644)
	}
	actualmd5 := structhash.Md5(actual, 1)
	expected, err := ioutil.ReadFile(fullGoldenFile)
	if err != nil {
		t.Fatal("Could not open golden file" + err.Error())
	}
	if string(actualmd5) != string(expected) {
		t.Errorf("Case failed: %s, %s", testFile, testDescription)
	}
}

func TestRead(t *testing.T) {
	cases := []struct {
		FileName        string
		goldenFile      string
		testDescription string
	}{
		{"test1_ParserOptions.yaml", "test1_md5.csv", "standard test"},
		{"test2_ParserOptions.yaml", "test2_md5.csv", "default data_type"},
		{"test3_ParserOptions.yaml", "test1_md5.csv", "default IfEmpty, IfInvalid"},
		{"test4_ParserOptions.yaml", "test1_md5.csv", "default jsonName"},
		{"test5_ParserOptions.yaml", "test1_md5.csv", "default validityRegexp"},
	}
	for _, tc := range cases {
		compareParserOptions(t, tc.FileName, tc.goldenFile, tc.testDescription)
	}
}
