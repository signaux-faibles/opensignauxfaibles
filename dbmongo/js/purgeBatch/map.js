function map() {
  if (this.value.batch[currentBatch]){
    delete this.value.batch[currentBatch]
  }
  // With a merge at the end, sending a new object, even empty, is compulsary
    emit(this._id, this.value)
}
