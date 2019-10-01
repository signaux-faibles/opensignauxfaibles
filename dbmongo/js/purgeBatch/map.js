function map() {
  if (this.value.batch[currentBatch]){
    delete this.value.batch[currentBatch]
  }
  
  emit(this._id, this.value)
}
