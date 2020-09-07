  
  const mapFunction = function() {
    const invalidEntries = [];
    for (const batchKey of Object.keys(this.value.batch)) {
      const batchEntries = this.value.batch[batchKey];
      for (const dataType of Object.keys(batchEntries)) {
        const dataEntries = batchEntries[dataType];
        for (const dataHash of Object.keys(dataEntries)) {
          const dataEntry = dataEntries[dataHash];
          if (dataType === "delai") {
            if (!(
              dataEntry.date_creation &&
              dataEntry.date_creation.getTime &&
              dataEntry.date_echeance &&
              dataEntry.date_echeance.getTime &&
              typeof dataEntry.duree_delai === "number" &&
              dataEntry.duree_delai > 0 &&
              typeof dataEntry.montant_echeancier === "number" &&
              dataEntry.montant_echeancier > 0.01
            )) {
              emit(this._id, { batchKey, dataHash, dataEntry }); // return invalid entry
            }
          }
        }
      }
    }
  };
  const reduceFunction = function(id, invalidEntries) {
    return {invalidEntries};
  };

  print("invalid records:");
  printjson(db.RawData.mapReduce(
    mapFunction,
    reduceFunction, {
      limit: 100000,
      allowDiskUse: true,
      out: { inline: 1 }
    }
  ));
  print("done.");
