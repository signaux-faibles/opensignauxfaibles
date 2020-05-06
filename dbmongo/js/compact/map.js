function map() {      
  "use strict";
  try{
    if (this.value != null) {
      emit(this.value.key, this.value) 
    }   
  } catch(error) {
    print(this.value.key)
  }
}
