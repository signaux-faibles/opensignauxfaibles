function lookAhead(data, attr_name, n_months, past){
 
  // Est-ce que l'évènement se répercute dans le passé (past = true; on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false; on
  // pourra se demander que s'est-il passé

  var sorting_fun = (
    (a, b) => a>=b
  ) 
  if (past) {
    sorting_fun = (
      (a, b) => a<=b
    )
  }

  var output = {}
  
  let counter = -1
  // Object.keys(data) représente les periodes 
  Object.keys(data).sort(sorting_fun).forEach( k => {
    if (counter >= 0) counter = counter + 1  // Si on a détecter quelque chose, on ajoute un à chaque période. 
    
    if (data[k][attr_name]){ // l'évènement se produit 
      counter = 0   
    }
    if (counter >= 0){ // l'évènement s'est déjà produit
      output[k] = output[k] || {}
      output[k].time_til_outcome = counter
    }
  })

  Object.keys(output).forEach( k => {
    if (output[k].time_til_outcome <= n_months){
      output[k].outcome = true
    } else output[k].outcome = false
  })
  return (output)
}
