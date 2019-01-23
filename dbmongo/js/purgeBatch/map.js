function map() {
  delete this.value.batch[currentBatch]
  emit(this._id, this.value) 
}