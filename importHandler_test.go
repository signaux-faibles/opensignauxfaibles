package main

// func Test_ImportBatch(t *testing.T) {

// 	err := ImportBatch(
// 		engine.AdminBatch{},
// 		[]engine.ParserType{},
// 		EmptyRegistry{},
// 		NoFilter,
// 		TestSinkFactory{},
// 		DiscardReportSink{},
// 	)

// 	if err == nil {
// 		t.Error("ImportBatch devrait nous empêcher d'importer sans filtre")
// 	}
// }

// func Test_ImportBatchWithUnreadableFilter(t *testing.T) {
// 	batch := engine.MockBatch("filter", []engine.BatchFile{engine.NewBatchFile("this_file_does_not_exist")})

// 	err := ImportBatch(
// 		batch,
// 		[]engine.ParserType{},
// 		// TODO check
// 		nil,
// 		NoFilter,
// 		TestSinkFactory{},
// 		DiscardReportSink{},
// 	)
// 	if err == nil {
// 		t.Error("ImportBatch devrait échouer en tentant d'ouvrir un fichier filtre illisible")
// 	}
// }
