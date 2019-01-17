function cibleApprentissage(output_indexed) {
  counter = -1
  Object.keys(output_indexed).sort((a,b)=> a<=b).forEach( k => {
    if (counter >=0) counter = counter + 1 
    if (output_indexed[k].tag_outcome == "default" || output_indexed[k].tag_outcome == "failure"){
      counter = 0 
    }
    if (counter >= 0){
      output_indexed[k].time_til_outcome = counter
    }
  })
}