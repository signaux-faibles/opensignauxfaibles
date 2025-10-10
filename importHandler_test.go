package main

// func Test_ImportBatch(t *testing.T) {

// 	err := ImportBatch(
// 		base.AdminBatch{},
// 		[]base.ParserType{},
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
// 	batch := base.MockBatch("filter", []base.BatchFile{base.NewBatchFile("this_file_does_not_exist")})

// 	err := ImportBatch(
// 		batch,
// 		[]base.ParserType{},
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
