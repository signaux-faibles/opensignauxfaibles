function map() {
  "use strict";
  const batches = Object.keys(this.value.batch)
  batches.filter((key) => key >= fromBatchKey).forEach((key) => {
    delete this.value.batch[key]
  })
  // With a merge output, sending a new object, even empty, is compulsory
  emit(this._id, this.value)
}
