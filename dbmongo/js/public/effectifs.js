"use strict";

function effectifs(v) {
  const  { serie_periode } = this; // parameters passed from Go Pipeline as global variables  

  var mapEffectif = {}
  f.iterable(v.effectif).forEach(e => {
    mapEffectif[e.periode.getTime()] = (mapEffectif[e.periode.getTime()] || 0) + e.effectif
  })
  return serie_periode.map(p => {
    return {
      periode: p,
      effectif: mapEffectif[p.getTime()] || null
    }
  }).filter(p => p.effectif)
}