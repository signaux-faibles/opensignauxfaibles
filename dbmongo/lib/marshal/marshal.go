package marshal

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
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

// GenericMarshal marshals a file depending on ParserOptions
func GenericMarshal(
	po *ParserOptions,
	cache engine.Cache,
	batch *engine.AdminBatch,
) (chan engine.Tuple, chan engine.Event, error) {

	object := Object{}
	object.Data = make(bson.M)

	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	go func() {
		defer close(eventChannel)
		defer close(outputChannel)

		for _, path := range batch.Files[po.DataType] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports,
			)

			// get current file name
			fullPath := viper.GetString("APP_DATA") + path
			openAndReadFile(po, fullPath, &tracker, outputChannel, eventChannel, cache, batch)
		}
	}()
	return outputChannel, eventChannel, nil
}

func openAndReadFile(
	po *ParserOptions,
	fullPath string,
	tracker *gournal.Tracker,
	outputChannel chan engine.Tuple,
	eventChannel chan engine.Event,
	cache engine.Cache,
	batch *engine.AdminBatch,
) {

	event := engine.Event{
		Code:    engine.Code("parser" + po.StructName),
		Channel: eventChannel,
	}

	file, err := os.Open(fullPath)
	defer file.Close()

	if err != nil {
		tracker.Error(engine.NewCriticError(err, "fatal"))
		event.CriticalReport("fatalError", *tracker)
		return
	} else {
		event.Info(fullPath + ": ouverture")
	}

	event.Info(fullPath + ": ouverture")

	readFile(po, file, tracker, event, outputChannel, cache, batch)

	event.InfoReport("abstract", *tracker)
}

func readFile(
	po *ParserOptions,
	file io.Reader,
	tracker *gournal.Tracker,
	event engine.Event,
	outputChannel chan engine.Tuple,
	cache engine.Cache,
	batch *engine.AdminBatch,
) {

	var marshallingMap = make(map[string]int)
	for _, field := range po.Fields {
		if field.CSVName != "" {
			marshallingMap[field.CSVName] = field.CSVCol
		}
	}

	reader := csv.NewReader(file)
	reader.Comma = []rune(po.Delimiter)[1]
	reader.LazyQuotes = true

	if po.HasHeader {
		row, err := reader.Read()
		if err != nil {
			tracker.Error(engine.NewCriticError(err, "fatal"))
			event.CriticalReport("fatalError", *tracker)
			return
		}
		err = CheckMarshallingMap(row, marshallingMap)
		if err != nil {
			tracker.Error(engine.NewCriticError(err, "fatal"))
		}
	}

	for {
		if tracker.Count%10000 == 0 && engine.ShouldBreak(*tracker, engine.MaxParsingErrors) {
			break
		}
		tracker.Next()

		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			(*tracker).Error(err)
			event.CriticalReport("fatalError", *tracker)
			break
		}

		obj := readOneLine(po, row, marshallingMap, tracker)

		if po.KeyMapping == "urssaf" {
			// mapping urssaf accounts to siret
			obj.Data["key"], err = GetSiret(
				obj.Data["NumeroCompte"].(string),
				obj.Data["Period"].(*time.Time),
				cache,
				batch,
			)

			if err != nil {
				tracker.Error(engine.NewMappingError(err, "error"))
				continue
			}
		}

		filtered, err := IsFiltered(obj.Key(), cache, batch)
		if err != nil {
			tracker.Error(engine.NewFilterError(err, "fatal"))
			break
		}
		if filtered {
			tracker.Error(engine.NewFilterError(errors.New("Row filtered from input filter file"), "filter"))
		}

		// Immediate fail if any fatal error
		anyFatal := false
		for _, e := range tracker.ErrorsInCurrentCycle() {
			switch c := e.(type) {
			case engine.CriticityError:
				if c.Criticity() == "fatal" {
					anyFatal = true
				}
			default:
			}
		}
		if anyFatal {
			event.CriticalReport("fatalError", *tracker)
			break
		}
		if !tracker.HasErrorInCurrentCycle() {
			outputChannel <- obj
		}
	}
}

//readOneLine reads one line of an input file
func readOneLine(po *ParserOptions, row []string, marshallingMap map[string]int, tracker *gournal.Tracker) Object {

	var cerr *engine.CriticError
	var err error
	var ok bool
	var data bson.M
	for _, field := range po.Fields {
		data[field.JSONName], err = ParserDict[field.Parser](
			row[marshallingMap[field.CSVName]],
			field.IfEmpty,
			field.IfInvalid,
			sfregexp.RegexpDict[field.ValidityRegex],
			field.TimeFormat,
		)

		cerr, ok = err.(*engine.CriticError)
		if !ok {
			tracker.Error(err)
		} else {
			tracker.Error(engine.NewParseError(cerr, data[field.JSONName].(string)))
		}

		if field.Mapping != "" {
			data[field.JSONName] = MappingDict[field.Mapping](data)
		}
	}
	return Object{Data: data, key: data[po.Key].(string), datatype: po.DataType, scope: po.Scope}
}
