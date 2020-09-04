  
  const mapFunction = function() {
    const dataType = "delai";
    const invalidEntries = [];
    for (const batchKey of Object.keys(this.value.batch)) {
      const batchEntries = this.value.batch[batchKey];
      const entriesToValidate = batchEntries[dataType];
      if (entriesToValidate) {
        for (const dataHash of Object.keys(entriesToValidate)) {
          const dataEntry = entriesToValidate[dataHash];
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
  };
  const reduceFunction = function(id, invalidEntries) {
    return {invalidEntries};
  };

  print("invalid records:");
  printjson(db.RawData.mapReduce(
    mapFunction,
    reduceFunction, {
      limit: 1000000,
      out: { inline: 1 }
    }
  ));
  print("done.");
