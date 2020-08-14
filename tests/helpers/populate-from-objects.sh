#!/bin/bash

# This helper generates JS commands for mongo shell to add sample data from objects.js.

cat << CONTENTS
  db.Admin.insertOne({
    "_id" : {
        "key" : "1905",
        "type" : "batch"
    },
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })
CONTENTS

RAW_DATA_DOCUMENTS=$(node --print -e "require('./dbmongo/js/test/data/objects.js').makeObjects.toString().replace('ISODate => ([', '[').replace('])', ']')")
echo "db.RawData.insertMany(${RAW_DATA_DOCUMENTS})"
