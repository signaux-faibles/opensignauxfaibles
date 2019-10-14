package marshal

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:generate go run generateParser.go

//Field and options
type Field struct {
	GoName        string `yaml:"go_name"`
	CSVName       string `yaml:"csv_name"`
	CSVCol        int    `yaml:"csv_col"`
	JSONName      string `yaml:"json_name"`
	Parser        string `yaml:"parser"`
	IfEmpty       string `yaml:"if_empty"`
	IfInvalid     string `yaml:"if_invalid"`
	Mapping       string `yaml:"mapping"`
	ValidityRegex string `yaml:"validity_regex"`
	TimeFormat    string `yaml:"time_format"`
}

//ParserOptions lists parser general options
type ParserOptions struct {
	FileType   string  `yaml:"file_type"`
	Delimiter  string  `yaml:"delimiter"`
	DataType   string  `yaml:"data_type"`
	StructName string  `yaml:"struct_name"`
	Scope      string  `yaml:"scope"`
	Key        string  `yaml:"key"`
	KeyMapping string  `yaml:"key_mapping"`
	HasHeader  bool    `yaml:"has_header"`
	Fields     []Field `yaml:"fields"`
}

// GetTypeFromFilepath takes the last element of the path and removes the
// extension
func GetTypeFromFilepath(filePath string) string {
	return (strings.TrimSuffix(
		filepath.Base(filePath),
		filepath.Ext(filePath),
	))
}

//ReadOptions reads ParserOptions from a yaml file
func (po *ParserOptions) ReadOptions(filepath string) error {
	// var possibleBehaviors = []string{"fatal", "filter", "error", "ignore"}
	// var possibleParser = PossibleParsers()
	// var possibleRegexp = sfregexp.PossibleRegexp()
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, po)
	if po.DataType == "" {
		// Default value for DataType is filename
		po.DataType = GetTypeFromFilepath(filepath)
	}
	if po.KeyMapping != "" {
		if po.Key != "" {
			return errors.New(`Cannot have both a "key_mapping" and a specified "key"`)
		}
		po.Key = "key"
	}
	if po.FileType == "xlsx" {
		po.Delimiter = ";"
	}

	// default field values
	for ind := range po.Fields {
		if po.Fields[ind].IfInvalid == "" {
			po.Fields[ind].IfInvalid = "fatal"
		}
		if po.Fields[ind].IfEmpty == "" {
			po.Fields[ind].IfEmpty = "fatal"
		}
		if po.Fields[ind].JSONName == "" {
			po.Fields[ind].JSONName = toSnakeCase(po.Fields[ind].GoName)
		}
		if po.Fields[ind].ValidityRegex == "" {
			po.Fields[ind].ValidityRegex = "nil"
		}
	}

	// Options
	for ind := range po.Fields {
		if (po.Fields[ind].TimeFormat == "" && po.Fields[ind].Parser == "time") &&
			(po.Fields[ind].TimeFormat != "" && po.Fields[ind].Parser != "time") {
			return errors.New("time_format option is only valid for time parsers")
		}
	}

	if err != nil {
		return err
	}
	return nil
}

// RegisteredParserOptions ...
func RegisteredParserOptions(folder string) (map[string]*ParserOptions, error) {
	pomap := make(map[string]*ParserOptions)

	// Lecture des noms de fichiers de ParserOptions
	dir, err := os.Open(folder)
	if err != nil {
		return nil, err
	}
	files, err := dir.Readdir(-1)
	dir.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != "yaml" {
			break
		}
		typename := GetTypeFromFilepath(file.Name())
		var po *ParserOptions
		err = po.ReadOptions(file.Name())
		if err != nil {
			return nil, err
		}
		pomap[typename] = po
	}
	return pomap, nil
}

// toSnakeCase converts a string to snake case.
func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
