function ccsf(v, output_array)
  var ccsfHashes = Object.keys(v.ccsf || {}) 

  output_array.forEach(val => {        
      var optccsf = ccsfHashes.reduce( 
          function (accu, hash) { 
              ccsf = v.ccsf[hash] 
              if (ccsf.date_traitement.getTime() < val.periode.getTime() && ccsf.date_traitement.getTime() > accu.date_traitement.getTime()) { 
                  accu = ccsf 
              } 
              return(accu)
          }, 
          { 
              date_traitement: new Date(0) 
          } 
      )         
      
      if (optccsf.date_traitement.getTime() != 0) { 
          val.date_ccsf = optccsf.date_traitement 
      } 
  })