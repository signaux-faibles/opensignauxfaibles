db.getCollection("Features").aggregate(
	[
		// -- Stage 1 --
		{
			$project: {
			    "_id": {
			        "batch": "$info.batch",
			        "siret": "$value.siret",
			        "periode": "$info.periode"
			    },
			    "value": "$value"
			}
		},

		// -- Stage 2 --
		{
			$out: "Features"
		},
	]
);


db.getCollection("Features").dropIndex({
    "info.batch" : 1,
    "value.random_order" : -1,
    "info.periode" : 1,
    "value.effectif" : 1,
    "info.siren" : 1
})



db.getCollection("Features").createIndex({
    "_id.batch" : 1,
    "value.random_order" : -1,
    "_id.periode" : 1,
    "value.effectif" : 1,
    "_id.siren" : 1
})

