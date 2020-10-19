#!/bin/bash

# This helper generates JS commands for mongo shell to add sample data from objects.js.

function raw_data_documents {
  node --print -e '
    const { makeObjects } = require("./dbmongo/js/test/data/objects.js");
    makeObjects.toString()
      .replace("ISODate => ([", "[")
      .replace("])", "]");
  '
}

echo "
  db.Admin.insertOne({
    '_id' : {
        'key' : '1901',
        'type' : 'batch'
    },
    'param' : {
        'date_debut' : ISODate('2014-01-01T00:00:00.000+0000'),
        'date_fin' : ISODate('2016-01-01T00:00:00.000+0000'),
        'date_fin_effectif' : ISODate('2016-03-01T00:00:00.000+0000')
    },
    'name' : 'TestData'
  })

  db.RawData.insertMany(
    $(raw_data_documents)
  )
"
