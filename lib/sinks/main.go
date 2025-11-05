// Package sinks fournit les destinations de sortie pour les données traitées
// par le moteur d'import de Signaux Faibles.
//
// Le package définit l'interface DataSink qui permet d'envoyer les tuples
// produits par les parsers vers différentes destinations :
//   - PostgresSink : écrit les données dans une base PostgreSQL
//   - CsvSink : exporte les données au format CSV
//   - StdoutSink : affiche les données sur la sortie standard
//   - CompositeSink : combine plusieurs sinks pour écrire simultanément vers
//     plusieurs destinations
//
// Les sinks sont créés via des factories (SinkFactory) qui permettent
// d'instancier le sink approprié en fonction du type de parser.
//
// Exemple d'utilisation :
//
//	factory := NewPostgresSinkFactory(dbPool)
//	sink, err := factory.CreateSink(engine.Sirene)
//	if err != nil {
//	    return err
//	}
//	err = sink.ProcessOutput(ctx, tupleChannel)
//
// Pour combiner plusieurs sinks :
//
//	factory := Combine(
//	    NewPostgresSinkFactory(dbPool),
//	    NewCsvSinkFactory(outputDir),
//	)
package sinks
