package engine

import (
	"encoding/csv"
	"log/slog"
	"os"

	"github.com/globalsign/mgo/bson"

	"opensignauxfaibles/lib/marshal"
)

var csvFiles = map[string]*os.File{}

func InsertIntoCSVs() chan *Value {
	importing.Add(1)
	source := make(chan *Value, 10)
	defer closeCSVs()
	go func(chan *Value) {
		defer importing.Done()
		buffer := make(map[string]*Value)
		i := 0
		insertObjectsIntoImportedData := func() {
			for _, v := range buffer {
				writeBatchesToCSV(v.Value.Batch)
			}
			buffer = make(map[string]*Value)
			i = 0
		}

		for value := range source {
			if i >= 100 {
				insertObjectsIntoImportedData()
			}
			if knownValue, ok := buffer[value.Value.Key]; ok {
				newValue, _ := (*knownValue).Merge(*value)
				buffer[value.Value.Key] = &newValue
			} else {
				value.ID = bson.NewObjectId()
				buffer[value.Value.Key] = value
				i++
			}
		}
		// le canal a été fermé => importer les données restantes avant de rendre la main
		insertObjectsIntoImportedData()
	}(source)

	return source
}

func closeCSVs() {
	for _, file := range csvFiles {
		err := file.Close()
		slog.Error(
			"erreur pendant la fermeture du fichier",
			slog.Any("error", err),
			slog.String("filename", file.Name()),
		)
	}
}

func writeBatchesToCSV(batchs map[string]Batch) {
	for _, v := range batchs {
		writeBatchToCSV(v)
	}
}

func writeBatchToCSV(batch Batch) {
	for _, tuples := range batch {
		writeLinesToCSV(tuples)
	}
}

func writeLinesToCSV(tuples map[string]marshal.Tuple) {
	for _, tuple := range tuples {
		logger := slog.Default().With(slog.Any("tuple", tuple))
		csvWriter := openFile(tuple)
		err := csvWriter.Write(tuple.Values())
		if err != nil {
			logger.Error("erreur pendant l'écriture du tuple en csvWriter")
		}
		csvWriter.Flush()
	}
}

func openFile(tuple marshal.Tuple) *csv.Writer {
	logger := slog.Default().With(slog.Any("tuple", tuple))
	file, found := csvFiles[tuple.Type()]
	if found {
		return csv.NewWriter(file)
	}
	var err error
	filename := string(tuple.Type()) + ".csv"
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	csvFiles[tuple.Type()] = file
	if err != nil {
		logger.Error(
			"erreur pendant l'ouverture du fichier",
			slog.String("filename", filename),
			slog.Any("error", err),
		)
		panic(err)
	}
	writer := csv.NewWriter(file)
	logger.Warn("write headers")
	err = writer.Write(tuple.Headers())
	logger.Error(
		"erreur pendant l'écriture des headers'",
		slog.String("filename", filename),
		slog.Any("error", err),
	)
	writer.Flush()
	return writer
}
