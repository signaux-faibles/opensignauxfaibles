function map() {
  if (this.value.batch[currentBatch]){
    delete this.value.batch[currentBatch]
  }
  if (Object.keys(this.value.batch).length > 0){
    emit(this._id, this.value)
  }
}
