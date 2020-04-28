// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

// creates all parsers which options files are in folderName
func main() {
	var folderName = "parserSpecs"
	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.Fatal("Finds could not be found: ", err.Error())
	}

	for _, file := range files {
		po := marshal.ParserOptions{}
		wd, _ := os.Getwd()
		fileName := filepath.Join(wd, folderName, file.Name())
		err := po.ReadOptions(fileName)
		if err != nil {
			log.Fatal("File could not be read: ", fileName, "\n", err.Error())
		}
		dirName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create directory: ", dirName, err.Error())
		}
		// GenerateParser(po, filepath.Join(dirName, dirName+"Parser.go"))
	}

}

//GenerateParser generates a parser from a parser option file
//func GenerateParser(po marshal.ParserOptions, outputFile string) {
//	out, err := os.Create(outputFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer out.Close()

//	temp, err := parserTemplate()
//	if err != nil {
//		log.Fatal("Parser could not be generated: " + err.Error())
//	}
//	temp.Execute(out, po)
//}

//func parserTemplate() (*template.Template, error) {
//	temp, err := template.New("Parser").Funcs(
//		template.FuncMap{
//			"getParserType": marshal.GetParserType,
//		}).Parse(`// Code generated by go generate; DO NOT EDIT.

//{{- /* Define a backtick variable */ -}}
//{{- $tick := "` + "`" + `" }}
//package {{ .DataType }}

//import (
//	"opensignauxfaibles/dbmongo/lib/engine"
//	"opensignauxfaibles/dbmongo/lib/marshal"
//	{{ if eq .FileType "csv" -}}
//	"encoding/csv"
//	"os"
//	{{ else if eq .FileType "xlsx" -}}
//	"github.com/tealeg/xlsx"
//	{{ end -}}
//	"errors"
//	"io"
//	"time"

//	"github.com/signaux-faibles/gournal"

//	"github.com/spf13/viper"
//)

//// {{ .StructName }} is an automatically generated parsing struct
//type {{ .StructName }} struct {
//{{- if .KeyMapping }}
//key string
//{{- end }}
//{{- range .Fields }}
//	{{ .GoName }} {{ getParserType .Parser }} {{ $tick -}} json:" {{- .JSONName -}} " bson:" {{- .JSONName -}} " {{- $tick }}
//{{- end }}
//}

//// Key de l'objet
//func (obj {{ .StructName -}} ) Key() string {
//	return obj. {{- .Key }}
//}

//// Scope de l'objet
//func (rep {{ .StructName -}} ) Scope() string {
//	return "{{ .Scope }}"
//}

//// Type de l'objet
//func (rep {{ .StructName -}} ) Type() string {
//	return "{{ .DataType }}"
//}

//// Parser fonction qui retourne data et journaux
//func Parser(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
//	outputChannel := make(chan engine.Tuple)
//	eventChannel := make(chan engine.Event)

//	go func() {
//		defer close(eventChannel)
//		defer close(outputChannel)
//		for _, path := range batch.Files[" {{- .DataType -}} "] {
//			tracker := gournal.NewTracker(
//				map[string]string{"path": path},
//				engine.TrackerReports,
//			)

//			// get current file name
//			fullPath := viper.GetString("APP_DATA") + path
//			openAndReadFile(fullPath, &tracker, outputChannel, eventChannel, cache, batch)
//		}
//	}()
//	return outputChannel, eventChannel
//}

//func openAndReadFile(
//  fullPath string,
//	tracker *gournal.Tracker,
//	outputChannel chan engine.Tuple,
//	eventChannel chan engine.Event,
//	cache engine.Cache,
//	batch *engine.AdminBatch,
//	){

//	event := engine.Event{
//		Code:    "parser {{- .StructName }}",
//		Channel: eventChannel,
//	}

//	{{ if eq .FileType "csv" -}}
//	file, err := os.Open(fullPath)
//	defer file.Close()
//	{{ else if eq .FileType "xlsx" -}}
//	xlFile, err := xlsx.OpenFile(fullPath)
//	{{- end }}
//	if err != nil {
//		tracker.Error(engine.NewCriticError(err, "fatal"))
//		event.CriticalReport("fatalError", *tracker)
//		return
//	}

//	{{- if eq .FileType "xlsx" }}
//	sheet := file.Sheets[0]
//	file := marshal.ReadXlsx(sheet.Rows[:])
//	{{- end }}

//	event.Info(fullPath + ": ouverture")
//	readFile(file, tracker, event, outputChannel, cache, batch)

//	event.InfoReport("abstract", *tracker)
//}

//func readFile(
//  file io.Reader,
//	tracker *gournal.Tracker,
//	event engine.Event,
//	outputChannel chan engine.Tuple,
//	cache engine.Cache,
//	batch *engine.AdminBatch,
//	){

//	var marshallingMap = make(map[string]int)
//	{{- range .Fields }}
//	{{- if .CSVName }}
//	marshallingMap["{{ .CSVName }}"] = {{ .CSVCol }}
//	{{- end }}
//	{{- end }}

//	reader := csv.NewReader(file)
//	reader.Comma = '{{ .Delimiter }}'
//	reader.LazyQuotes = true

//	{{ if .HasHeader -}}
//	row, err := reader.Read()
//	if err != nil {
//		tracker.Error(engine.NewCriticError(err, "fatal"))
//		event.CriticalReport("fatalError", *tracker)
//		return
//	}
//	err = marshal.CheckMarshallingMap(row, marshallingMap)
//	if err != nil {
//	  tracker.Error(engine.NewCriticError(err, "fatal"))
//		}
//	{{- end }}
//	for {
//		tracker.Next()
//		row, err := reader.Read()
//		if err == io.EOF {
//			break
//		} else if err != nil {
//			(*tracker).Error(err)
//			event.CriticalReport("fatalError", *tracker)
//			break
//		}

//		{{ .DataType }} := readOneLine(row, marshallingMap, tracker)

//    {{- if eq .KeyMapping "urssaf" }}
//    // mapping urssaf accounts to siret
//		{{ .DataType }}.key, err = marshal.GetSiret({{ .DataType }}.NumeroCompte, {{ .DataType }}.Period, cache, batch)
//		if err != nil {
//		  tracker.Error(engine.NewMappingError(err, "error"))
//			continue
//		}
//    {{- end }}

//    filtered, err := marshal.IsFiltered({{ .DataType }}.{{ .Key }}, cache, batch)
//		if err != nil {
//		  tracker.Error(engine.NewFilterError(err, "fatal"))
//		  break
//		}
//		if filtered {
//			tracker.Error(engine.NewFilterError(errors.New("Row filtered from input filter file"), "filter"))
//		}

//		// Immediate fail if any fatal error
//	  anyFatal := false
//	  for _, e := range tracker.ErrorsInCurrentCycle() {
//			switch c := e.(type){
//			case engine.CriticityError:
//			  if c.Criticity() == "fatal" {
//				  anyFatal = true
//					}
//		  default:
//			}
//		}
//		if anyFatal {
//			event.CriticalReport("fatalError", *tracker)
//			break
//		}
//		if !tracker.HasErrorInCurrentCycle() {
//			outputChannel <- {{ .DataType }}
//		}
//	}
//}

////readOneLine
//func readOneLine(row []string, marshallingMap map[string]int, tracker *gournal.Tracker) {{ .StructName }} {

//	var cerr *engine.CriticError
//	var ok bool
//	{{- range .Fields }}
//	{{ .GoName }} , err := marshal.{{ .Parser }}(
//		row[marshallingMap[" {{- .CSVName -}}"]],
//		" {{- .IfEmpty -}} ",
//		" {{- .IfInvalid -}} ",
//		{{ .ValidityRegex }},
//	)

//	cerr, ok = err.(*engine.CriticError)
//	if !ok {
//	   tracker.Error(err)
//	} else {
//		 tracker.Error(engine.NewParseError(cerr, "{{ .GoName -}}" ))
//	}

//	{{- if .Mapping }}
//	{{ .GoName }} = marshal.{{ .Mapping -}} ( {{- .GoName -}} )
//	{{- end }}
//  {{- end }}

//	obj := {{ .StructName }}{}
//	{{- range .Fields }}
//		obj.{{ .GoName }} = {{ .GoName }}
//	{{- end }}
//	return obj
//}`)
//	return temp, err
//}
