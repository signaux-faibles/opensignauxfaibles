package engine 

 var jsFunctions = map[string]map[string]string{
"common":{
"altaresToHuman": `function altaresToHuman (code) {
  "use strict";
  var codeLiquidation = ['PCL0108', 'PCL010801','PCL010802','PCL030107','PCL030307','PCL030311','PCL05010103','PCL05010204','PCL05010303','PCL05010403','PCL05010503','PCL05010703','PCL05011004','PCL05011102','PCL05011204','PCL05011206','PCL05011304','PCL05011404','PCL05011504','PCL05011604','PCL05011903','PCL05012004','PCL050204','PCL0109','PCL010901','PCL030108','PCL030308','PCL05010104','PCL05010205','PCL05010304','PCL05010404','PCL05010504','PCL05010803','PCL05011005','PCL05011103','PCL05011205','PCL05011207','PCL05011305','PCL05011405','PCL05011505','PCL05011904','PCL05011605','PCL05012005'];
  var codePlanSauvegarde = ['PCL010601','PCL0106','PCL010602','PCL030103','PCL030303','PCL03030301','PCL05010101','PCL05010202','PCL05010301','PCL05010401','PCL05010501','PCL05010506','PCL05010701','PCL05010705','PCL05010801','PCL05010805','PCL05011002','PCL05011202','PCL05011302','PCL05011402','PCL05011502','PCL05011602','PCL05011901','PCL0114','PCL030110','PCL030310'];
  var codeRedressement = ['PCL0105','PCL010501','PCL010502','PCL010503','PCL030105','PCL030305','PCL05010102','PCL05010203','PCL05010302','PCL05010402','PCL05010502','PCL05010702','PCL05010706','PCL05010802','PCL05010806','PCL05010901','PCL05011003','PCL05011101','PCL05011203','PCL05011303','PCL05011403','PCL05011503','PCL05011603','PCL05011902','PCL05012003'];
  var codeInBonis = ['PCL05','PCL0501','PCL050101','PCL050102','PCL050103','PCL050104','PCL050105','PCL050106','PCL050107','PCL050108','PCL050109','PCL050110','PCL050111','PCL050112','PCL050113','PCL050114','PCL050115','PCL050116','PCL050119','PCL050120','PCL050121','PCL0503','PCL050301','PCL050302','PCL0508','PCL010504','PCL010803','PCL010902','PCL050901','PCL050902','PCL050903','PCL050904','PCL0504','PCL050303','PCL050401','PCL050402','PCL050403','PCL050404','PCL050405','PCL050406'];
  var codeContinuation = ['PCL0202'];
  var codeSauvegarde = ['PCL0203','PCL020301','PCL0205','PCL040408'];
  var codeCession = ['PCL0204','PCL020401','PCL020402','PCL020403'];
  var res = null;
  if (codeLiquidation.includes(code)) 
    res = 'liquidation';
  else if (codePlanSauvegarde.includes(code))
    res = 'plan_sauvegarde';
  else if (codeRedressement.includes(code))
    res = 'plan_redressement';
  else if (codeInBonis.includes(code))
    res = 'in_bonis';
  else if (codeContinuation.includes(code))
    res = 'continuation';
  else if (codeSauvegarde.includes(code))
    res = 'sauvegarde';
  else if (codeCession.includes(code))
    res = 'cession';
  return res;
}`,
"generatePeriodSerie": `function generatePeriodSerie (date_debut, date_fin) {
  "use strict";
  var date_next = new Date(date_debut.getTime())
  var serie = []
  while (date_next.getTime() < date_fin.getTime()) {
    serie.push(new Date(date_next.getTime()))
    date_next.setUTCMonth(date_next.getUTCMonth() + 1)
  }
  return serie
}`,
"raison_sociale": `function raison_sociale /*eslint-disable-line @typescript-eslint/no-unused-vars */(denomination_unite_legale, nom_unite_legale, nom_usage_unite_legale, prenom1_unite_legale, prenom2_unite_legale, prenom3_unite_legale, prenom4_unite_legale) {
    "use strict";
    const nomUsageUniteLegale = nom_usage_unite_legale
        ? nom_usage_unite_legale + "/"
        : "";
    const raison_sociale = denomination_unite_legale ||
        (nom_unite_legale +
            "*" +
            nomUsageUniteLegale +
            prenom1_unite_legale +
            " " +
            (prenom2_unite_legale || "") +
            " " +
            (prenom3_unite_legale || "") +
            " " +
            (prenom4_unite_legale || "") +
            " ").trim() + "/";
    return raison_sociale;
}`,
"region": `function region(departement){
  "use strict";
  var reg = ""
  switch (departement){
    case "01":
    case "03":
    case "07":
    case "15":
    case "26":
    case "38":
    case "42":
    case "43":
    case "63":
    case "69":
    case "69":
    case "73":
    case "74":
      reg = "Auvergne-Rhône-Alpes"
      break
    case "02":
    case "59":
    case "60":
    case "62":
    case "80":
      reg = "Hauts-de-France"
      break
    case "04":
    case "05":
    case "06":
    case "13":
    case "83":
    case "84":
      reg = "Provence-Alpes-Côte d'Azur"
      break
    case "08":
    case "10":
    case "51":
    case "52":
    case "54":
    case "55":
    case "57":
    case "67":
    case "68":
    case "88":
      reg = "Grand Est"
      break
    case "09":
    case "11":
    case "12":
    case "30":
    case "31":
    case "32":
    case "34":
    case "46":
    case "48":
    case "65":
    case "66":
    case "81":
    case "82":
      reg = "Occitanie"
      break
    case "14":
    case "27":
    case "50":
    case "61":
    case "76":
      reg = "Normandie"
      break
    case "18":
    case "28":
    case "36":
    case "37":
    case "41":
    case "45":
      reg = "Centre-Val de Loire"
      break
    case "16":
    case "17":
    case "19":
    case "23":
    case "24":
    case "33":
    case "40":
    case "47":
    case "64":
    case "79":
    case "86":
    case "87":
      reg = "Nouvelle-Aquitaine"
      break
    case "20":
      reg = "Corse"
      break
    case "21":
    case "25":
    case "39":
    case "58":
    case "70":
    case "71":
    case "89":
    case "90":
      reg = "Bourgogne-Franche-Comté"
      break
    case "22":
    case "29":
    case "35":
    case "56":
      reg = "Bretagne"
      break
    case "44":
    case "49":
    case "53":
    case "72":
    case "85":
      reg = "Pays de la Loire"
      break
    case "75":
    case "77":
    case "78":
    case "91":
    case "92":
    case "93":
    case "94":
    case "95":
      reg = "Île-de-France"
      break
  }
  return(reg)
}`,
"setBatchValueForType": `// Cette fonction TypeScript permet de vérifier que seuls les types reconnus
// peuvent être intégrés dans un BatchValue de destination.
// Ex: setBatchValueForType(batchValue, "pouet", {}) cause une erreur ts(2345).
function setBatchValueForType(batchValue, typeName, updatedValues) {
    batchValue[typeName] = updatedValues;
}`,
},
"compact":{
"complete_reporder": `// complete_reporder ajoute une propriété "reporder" pour chaque couple
// SIRET+période, afin d'assurer la reproductibilité de l'échantillonage.
function complete_reporder(siret, object) {
    "use strict";
    const batches = Object.keys(object.batch);
    batches.sort();
    const missing = {};
    serie_periode.forEach((p) => {
        missing[p.getTime()] = true;
    });
    batches.forEach((batch) => {
        const reporder = object.batch[batch].reporder || {};
        Object.keys(reporder).forEach((ro) => {
            if (!missing[reporder[ro].periode.getTime()]) {
                delete object.batch[batch].reporder[ro];
            }
            else {
                missing[reporder[ro].periode.getTime()] = false;
            }
        });
    });
    const lastBatch = batches[batches.length - 1];
    serie_periode
        .filter((p) => missing[p.getTime()])
        .forEach((p) => {
        const reporder_obj = object.batch[lastBatch].reporder || {};
        reporder_obj[p.toString()] = {
            random_order: Math.random(),
            periode: p,
            siret: siret,
        };
        object.batch[lastBatch].reporder = reporder_obj;
    });
    return object;
}`,
"currentState": `// currentState() agrège un ensemble de batch, en tenant compte des suppressions
// pour renvoyer le dernier état connu des données.
// Note: similaire à flatten() de reduce.algo2.
function currentState(batches) {
    "use strict";
    const currentState = batches.reduce((m, batch) => {
        //1. On supprime les clés de la mémoire
        Object.keys((batch.compact || { delete: [] }).delete).forEach((type) => {
            batch.compact.delete[type].forEach((key) => {
                m[type].delete(key); // Should never fail or collection is corrupted
            });
        });
        //2. On ajoute les nouvelles clés
        Object.keys(batch)
            .filter((type) => type !== "compact")
            .forEach((type) => {
            m[type] = m[type] || new Set();
            Object.keys(batch[type]).forEach((key) => {
                m[type].add(key);
            });
        });
        return m;
    }, {});
    return currentState;
}`,
"finalize": `// finalize permet de:
// - indiquer les établissements à inclure dans les calculs de variables
// (processus reduce.algo2)
// - intégrer les reporder pour permettre la reproductibilité de
// l'échantillonnage pour l'entraînement du modèle.
function finalize(k, o) {
    "use strict";
    o.index = { algo1: false, algo2: false };
    if (o.scope === "entreprise") {
        o.index.algo1 = true;
        o.index.algo2 = true;
    }
    else {
        // Est-ce que l'un des batchs a un effectif ?
        const batches = Object.keys(o.batch);
        batches.some((batch) => {
            const hasEffectif = Object.keys(o.batch[batch].effectif || {}).length > 0;
            o.index.algo1 = hasEffectif;
            o.index.algo2 = hasEffectif;
            return hasEffectif;
        });
        // Complete reporder if missing
        // TODO: do not complete if all indexes are false.
        o = f.complete_reporder(k, o);
    }
    return o;
}`,
"map": `function map() {
    "use strict";
    if (typeof this.value !== "object") {
        throw new Error("this.value should be a valid object, in compact::map()");
    }
    emit(this.value.key, this.value);
}`,
"reduce": `// Entrée: données d'entreprises venant de ImportedData, regroupées par entreprise ou établissement.
// Sortie: un objet fusionné par entreprise ou établissement, contenant les données historiques et les données importées, à destination de la collection RawData.
// Opérations: retrait des données doublons et application des corrections de données éventuelles.
function reduce(key, values) {
    "use strict";
    // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
    // batchs. Sinon, juste fusion des attributs
    const auxBatchSet = new Set();
    const severalBatches = values.some((value) => {
        auxBatchSet.add(Object.keys(value.batch || {}));
        return auxBatchSet.size > 1;
    });
    //fusion des attributs dans values
    const reduced_value = values.reduce((m, value) => {
        Object.keys(value.batch).forEach((batch) => {
            m.batch[batch] = m.batch[batch] || {};
            Object.keys(value.batch[batch]).forEach((type) => {
                const updatedValues = Object.assign(Object.assign({}, m.batch[batch][type]), value.batch[batch][type]);
                setBatchValueForType(m.batch[batch], type, updatedValues);
            });
        });
        return m;
    }, { key: key, scope: values[0].scope, batch: {} });
    // Cette fonction reduce() est appelée à deux moments:
    // 1. agregation par établissement d'objets ImportedData. Dans cet étape, on
    // ne travaille généralement que sur un seul batch.
    // 2. agregation de ces résultats au sein de RawData, en fusionnant avec les
    // données potentiellement présentes. Dans cette étape, on fusionne
    // généralement les données de plusieurs batches. (données historiques)
    if (!severalBatches)
        return reduced_value;
    //////////////////////////////////////////////////
    // ETAPES DE LA FUSION AVEC DONNÉES HISTORIQUES //
    //////////////////////////////////////////////////
    // 0. On calcule la memoire au moment du batch à modifier
    const memory_batches = Object.keys(reduced_value.batch)
        .filter((batch) => batch < batchKey)
        .sort()
        .reduce((m, batch) => {
        m.push(reduced_value.batch[batch]);
        return m;
    }, []);
    const memory = f.currentState(memory_batches);
    // Pour tous les batchs à modifier, c'est-à-dire le batch ajouté et tous les
    // suivants.
    const modified_batches = batches.filter((batch) => batch >= batchKey);
    modified_batches.forEach((batch) => {
        reduced_value.batch[batch] = reduced_value.batch[batch] || {};
        // Les types où il y  a potentiellement des suppressions
        let stock_types = completeTypes[batch].filter((type) => (memory[type] || new Set()).size > 0);
        // Les types qui ont bougé dans le batch en cours
        let new_types = Object.keys(reduced_value.batch[batch]);
        // On dedoublonne au besoin
        let all_interesting_types = [...new Set([...stock_types, ...new_types])];
        // Filtrage selon les types effectivement importés
        if (types.length > 0) {
            stock_types = stock_types.filter((type) => types.includes(type));
            new_types = new_types.filter((type) => types.includes(type));
            all_interesting_types = all_interesting_types.filter((type) => types.includes(type));
        }
        // 1. On recupère les cles ajoutes et les cles supprimes
        // -----------------------------------------------------
        const hashToDelete = {};
        const hashToAdd = {};
        all_interesting_types.forEach((type) => {
            // Le type compact gère les clés supprimées
            if (type === "compact") {
                if (reduced_value.batch[batch].compact.delete) {
                    Object.keys(reduced_value.batch[batch].compact.delete).forEach((delete_type) => {
                        reduced_value.batch[batch].compact.delete[delete_type].forEach((hash) => {
                            hashToDelete[delete_type] =
                                hashToDelete[delete_type] || new Set();
                            hashToDelete[delete_type].add(hash);
                        });
                    });
                }
            }
            else {
                Object.keys(reduced_value.batch[batch][type] || {}).forEach((hash) => {
                    hashToAdd[type] = hashToAdd[type] || new Set();
                    hashToAdd[type].add(hash);
                });
            }
        });
        //
        // 2. On ajoute aux cles supprimees les types stocks de la memoire.
        // ----------------------------------------------------------------
        stock_types.forEach((type) => {
            hashToDelete[type] = new Set([
                ...(hashToDelete[type] || new Set()),
                ...memory[type],
            ]);
        });
        Object.keys(hashToDelete).forEach((type) => {
            // 3.a Pour chaque cle supprimee: est-ce qu'elle est bien dans la
            // memoire ? sinon on la retire de la liste des clés supprimées (pas de
            // maj memoire)
            // -----------------------------------------------------------------------------------------------------------------
            hashToDelete[type] = new Set([...hashToDelete[type]].filter((hash) => {
                return (memory[type] || new Set()).has(hash);
            }));
            // 3.b Est-ce qu'elle a ete egalement ajoutee en même temps que
            // supprimée ? (par exemple remplacement d'un stock complet à
            // l'identique) Dans ce cas là, on retire cette clé des valeurs ajoutées
            // et supprimées
            // i.e. on herite de la memoire. (pas de maj de la memoire)
            // ------------------------------------------------------------------------------
            hashToDelete[type] = new Set([...hashToDelete[type]].filter((hash) => {
                const also_added = (hashToAdd[type] || new Set()).has(hash);
                if (also_added) {
                    hashToAdd[type].delete(hash);
                }
                return !also_added;
            }));
            // 3.c On retire les cles restantes de la memoire.
            // --------------------------------------------------
            hashToDelete[type].forEach((hash) => {
                memory[type].delete(hash);
            });
        });
        Object.keys(hashToAdd).forEach((type) => {
            // 4.a Pour chaque cle ajoutee: est-ce qu'elle est dans la memoire ? Si oui on filtre cette cle
            // i.e. on herite de la memoire. (pas de maj de la memoire)
            // ---------------------------------------------------------------------------------------------
            hashToAdd[type] = new Set([...hashToAdd[type]].filter((hash) => {
                return !(memory[type] || new Set()).has(hash);
            }));
            // 4.b Pour chaque cle ajoutee restante: on ajoute à la memoire.
            // -------------------------------------------------------------
            hashToAdd[type].forEach((hash) => {
                memory[type] = memory[type] || new Set();
                memory[type].add(hash);
            });
        });
        // 5. On met à jour reduced_value
        // -------------------------------
        stock_types.forEach((type) => {
            if (hashToDelete[type]) {
                reduced_value.batch[batch].compact =
                    reduced_value.batch[batch].compact || {};
                reduced_value.batch[batch].compact.delete =
                    reduced_value.batch[batch].compact.delete || {};
                reduced_value.batch[batch].compact.delete[type] = [
                    ...hashToDelete[type],
                ];
            }
        });
        new_types.forEach((type) => {
            if (hashToAdd[type] && type !== "compact") {
                const hashedValues = reduced_value.batch[batch][type];
                const updatedValues = Object.keys(hashedValues || {})
                    .filter((hash) => {
                    return hashToAdd[type].has(hash);
                })
                    .reduce((m, hash) => {
                    m[hash] = hashedValues[hash];
                    return m;
                }, {});
                setBatchValueForType(reduced_value.batch[batch], type, updatedValues);
            }
        });
        // 6. nettoyage
        // ------------
        if (reduced_value.batch[batch]) {
            //types vides
            Object.keys(reduced_value.batch[batch]).forEach((type) => {
                if (Object.keys(reduced_value.batch[batch][type]).length === 0) {
                    delete reduced_value.batch[batch][type];
                }
            });
            //hash à supprimer vides (compact.delete)
            if (reduced_value.batch[batch].compact &&
                reduced_value.batch[batch].compact.delete) {
                Object.keys(reduced_value.batch[batch].compact.delete).forEach((type) => {
                    if (reduced_value.batch[batch].compact.delete[type].length === 0) {
                        delete reduced_value.batch[batch].compact.delete[type];
                    }
                });
                if (Object.keys(reduced_value.batch[batch].compact.delete).length === 0) {
                    delete reduced_value.batch[batch].compact;
                }
            }
            //batchs vides
            if (Object.keys(reduced_value.batch[batch]).length === 0) {
                delete reduced_value.batch[batch];
            }
        }
    });
    return reduced_value;
}`,
},
"crossComputation":{
"apart": `{
  "$set": {
    "value.ratio_apart": {
      "$divide": [
        "$value.apart_heures_consommees",
        {
          "$multiply": [
            "$value.effectif",
            157.67
          ]
        }
      ]
    }
  }
}`,
},
"migrations":{
"agg_change_index_Features": `// db.getCollection("Features").aggregate(
// 	[
// 		// -- Stage 1 --
// 		{
// 			$project: {
// 			    "_id": {
// 			        "batch": "$info.batch",
// 			        "siret": "$value.siret",
// 			        "periode": "$info.periode"
// 			    },
// 			    "value": "$value"
// 			}
// 		},

// 		// -- Stage 2 --
// 		{
// 			$out: "Features"
// 		},
// 	]
// );


db.getCollection("Features").dropIndex({
    "info.batch" : 1,
    "value.random_order" : -1,
    "info.periode" : 1,
    "value.effectif" : 1,
    "info.siren" : 1
})



db.getCollection("Features").createIndex({
    "_id.batch" : 1,
    "value.random_order" : -1,
    "_id.periode" : 1,
    "value.effectif" : 1,
    "_id.siret" : 1
})`,
},
"public":{
"apconso": `function apconso(apconso) {
  "use strict";
  return f.iterable(apconso).sort((p1, p2) => p1.periode < p2.periode)
}`,
"apdemande": `function apdemande(apdemande) {
  "use strict";
  return f.iterable(apdemande).sort((p1, p2) => p1.periode < p2.periode)
}`,
"bdf": `function bdf(hs) {
  "use strict";
  return f.iterable(hs).sort((a, b) => a.annee_bdf < b.annee_bdf)
}`,
"compareDebit": `function compareDebit (a,b) {
  "use strict";
  if (a.numero_historique < b.numero_historique) return -1
  if (a.numero_historique > b.numero_historique) return 1
  return 0
}`,
"compte": `function compte(compte) {
  "use strict";
  const c = f.iterable(compte)
  return (c.length>0)?c[c.length-1]:undefined
}`,
"cotisations": `function cotisations(vcotisation) {
  "use strict";
  var offset_cotisation = 0 
  var value_cotisation = {}
  
  // Répartition des cotisations sur toute la période qu'elle concerne
  vcotisation = vcotisation || {}
  Object.keys(vcotisation).forEach(function (h) {
    var cotisation = vcotisation[h]
    var periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = f.dateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })

  var output_cotisation = []

  serie_periode.forEach(p => {
    output_cotisation.push(
      (value_cotisation[p.getTime()] || []) 
        .reduce((m,c) => m+c, 0)
    )
  })

  return(output_cotisation)
}`,
"dateAddDay": `function dateAddDay(date, nbMonth) {
  "use strict";
  var result = new Date(date.getTime())
  result.setDate( result.getDate() + nbMonth );
  return result
}`,
"dateAddMonth": `function dateAddMonth(date, nbMonth) {
  "use strict";
  var result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}`,
"dealWithProcols": `function dealWithProcols(data_source, altar_or_procol, output_indexed){
  "use strict";
  return Object.keys(data_source || {}).reduce((events,hash) => {
    var the_event = data_source[hash]

    let etat = {}
    if (altar_or_procol == "altares")
      etat = f.altaresToHuman(the_event.code_evenement);
    else if (altar_or_procol == "procol")
      etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);

    if (etat != null)
      events.push({"etat": etat, "date_procol": new Date(the_event.date_effet)})

    return(events)
  },[]).sort(
    (a,b) => {return(a.date_procol.getTime() > b.date_procol.getTime())}
  )
}`,
"debits": `function debits(vdebit) {
  "use strict";

  const last_treatment_day = 20
  vdebit = vdebit || {}
  var ecn = Object.keys(vdebit).reduce((accu, h) => {
      let debit = vdebit[h]
      var start = debit.periode.start
      var end = debit.periode.end
      var num_ecn = debit.numero_ecart_negatif
      var compte = debit.numero_compte
      var key = start + "-" + end + "-" + num_ecn + "-" + compte
      accu[key] = (accu[key] || []).concat([{
          "hash": h,
          "numero_historique": debit.numero_historique,
          "date_traitement": debit.date_traitement
      }]) 
      return accu
  }, {})

  Object.keys(ecn).forEach(i => {
      ecn[i].sort(f.compareDebit)
      var l = ecn[i].length
      ecn[i].forEach((e, idx) => {
          if (idx <= l - 2) {
              vdebit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
          }
      })
  })

  var value_dette = {}

  Object.keys(vdebit).forEach(function (h) {
    var debit = vdebit[h]

    var debit_suivant = (vdebit[debit.debit_suivant] || {"date_traitement" : date_fin})
    
    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante. 
    let jour_traitement = debit.date_traitement.getUTCDate() 
    let jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate()
    let date_traitement_debut
    if (jour_traitement <= last_treatment_day){
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1)
      )
    }

    let date_traitement_fin
    if (jour_traitement_suivant <= last_treatment_day) {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + 1)
      )
    }

    let periode_debut = date_traitement_debut
    let periode_fin = date_traitement_fin

    //generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(periode_debut, periode_fin).map(date => {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })
  })    

  const output_dette = []
  serie_periode.forEach(p => {
    output_dette.push(
      (value_dette[p.getTime()] || [])
        .reduce((m,c) => {
          return {
            part_ouvriere: m.part_ouvriere + c.part_ouvriere,
            part_patronale: m.part_patronale + c.part_patronale,
            periode: f.dateAddDay(f.dateAddMonth(p,1),-1) }
          }, {part_ouvriere: 0, part_patronale: 0})
    )
  })

  return(output_dette)
}`,
"delai": `function delai(delai) {
  "use strict";
  return f.iterable(delai)
}`,
"diane": `function diane(hs) {
  "use strict";
 return f.iterable(hs).sort((a, b) => a.exercice_diane < b.exercice_diane)
}`,
"effectifs": `function effectifs(v) {
  "use strict";
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
}`,
"finalize": `function finalize(_, v) {
  "use strict";
  return v
}`,
"flatten": `function flatten(v, actual_batch) {
  "use strict";
  var res = Object.keys(v.batch || {})
    .sort()
    .filter(batch => batch <= actual_batch)
    .reduce((m, batch) => {

      // Types intéressants = nouveaux types, ou types avec suppressions
      var delete_types = Object.keys((v.batch[batch].compact || {}).delete || {})
      var new_types =  Object.keys(v.batch[batch])
      var all_interesting_types = [...new Set([...delete_types, ...new_types])]

      all_interesting_types.forEach(type => {
        m[type] = (m[type] || {})
        // On supprime les clés qu'il faut
        if (v.batch[batch] && v.batch[batch].compact && v.batch[batch].compact.delete &&
          v.batch[batch].compact.delete[type] && v.batch[batch].compact.delete[type] != {}) {

          v.batch[batch].compact.delete[type].forEach(hash => {
            delete m[type][hash]
          })
        }
        Object.assign(m[type], v.batch[batch][type])
      })
      return m
    }, { "key": v.key, scope: v.scope })

  return(res)
}`,
"idEntreprise": `function idEntreprise(idEtablissement) {
  "use strict";
  return {
    scope: 'entreprise',
    key: idEtablissement.slice(0,9),
    batch: actual_batch
  }
}`,
"iterable": `function iterable(dict) {
  "use strict";
  try {
    return Object.keys(dict).map(h => {
      return dict[h]
    })
  } catch(error) {
    return []
  }
}`,
"map": `function map() {
  "use strict";
  var value = f.flatten(this.value, actual_batch)

  if (this.value.scope=="etablissement") {
    let vcmde = {}
    vcmde.key = this.value.key
    vcmde.batch = actual_batch
    vcmde.effectif = f.effectifs(value)
    vcmde.dernier_effectif = vcmde.effectif[vcmde.effectif.length - 1]
    vcmde.sirene = f.sirene(f.iterable(value.sirene))
    vcmde.cotisation = f.cotisations(value.cotisation)
    vcmde.debit = f.debits(value.debit)
    vcmde.apconso = f.apconso(value.apconso)
    vcmde.apdemande = f.apconso(value.apdemande)
    vcmde.delai = f.delai(value.delai)
    vcmde.compte = f.compte(value.compte)
    vcmde.procol = f.dealWithProcols(value.altares, "altares",  null).concat(f.dealWithProcols(value.procol, "procol",  null))
    vcmde.last_procol = vcmde.procol[vcmde.procol.length - 1] || {"etat": "in_bonis"}
    vcmde.idEntreprise = "entreprise_" + this.value.key.slice(0,9)
    vcmde.procol = value.procol

    emit("etablissement_" + this.value.key, vcmde)
  }
  else if (this.value.scope == "entreprise") {
    let v = {}
    let diane = f.diane(value.diane)
    let bdf = f.bdf(value.bdf)
    let sirene_ul = (value.sirene_ul || {})[Object.keys(value.sirene_ul || {})[0] || ""]
    let crp = value.crp
    v.key = this.value.key
    v.batch = actual_batch
    
    if (diane.length > 0) {
      v.diane = diane
    }
    if (bdf.length > 0) {
      v.bdf = bdf
    }
    if (sirene_ul) {
      v.sirene_ul = sirene_ul
    }
    if (crp) {
      v.crp = crp
    }
    if (Object.keys(v) != []) {
      emit("entreprise_" + this.value.key, v)
    }
  }
}`,
"procolToHuman": `function procolToHuman (action, stade) {
  "use strict";
  var res = null;
  if (action == "liquidation" && stade != "abandon_procedure") 
    res = 'liquidation';
  else if (stade == "abandon_procedure" || stade == "fin_procedure")
    res = 'in_bonis';
  else if (action == "redressement" && stade == "plan_continuation")
    res = 'continuation';
  else if (action == "sauvegarde" && stade == "plan_continuation")
    res = 'sauvegarde';
  else if (action == "sauvegarde")
    res = 'plan_sauvegarde';
  else if (action == "redressement")
    res = 'plan_redressement';
  return res;
}`,
"reduce": `function reduce(key, values) {
  "use strict";
  if (key.scope="entreprise") {
    values = values.reduce((m, v) => {
      if (v.sirets) {
        m.sirets = (m.sirets || []).concat(v.sirets)
        delete v.sirets
      }
      Object.assign(m, v)
      return m
    }, {})
  }
  return values
}`,
"sirene": `function sirene(sireneArray) {
  "use strict";
  return sireneArray.reduce((accu, k) => {
    return k
  }, {})
}`,
},
"purgeBatch":{
"finalize": `function finalize(k, o) {
    "use strict";
    return o
}`,
"map": `function map() {
  "use strict";
  if (this.value.batch[currentBatch]){
    delete this.value.batch[currentBatch]
  }
  // With a merge at the end, sending a new object, even empty, is compulsary
    emit(this._id, this.value)
}`,
"reduce": `function reduce(key, values) {
    "use strict";
    return values
}`,
},
"reduce.algo2":{
"add": `function add(obj, output){
  "use strict";
  Object.keys(output).forEach(function(periode) {
    if (periode in obj){
      Object.assign(output[periode], obj[periode])
    } else {
      // throw new EvalError(
      //   "Attention, l'objet à fusionner ne possède pas les mêmes périodes que l'objet dans lequel il est fusionné"
      // )
    }
  })
}`,
"apart": `function apart (apconso, apdemande) {
  "use strict";

  var output_apart = {}

  // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
  var apart = Object.keys(apdemande).reduce((apart, hash) => {
    apart[apdemande[hash].id_demande.substring(0, 9)] = {
      "demande": hash,
      "consommation": [],
      "periode_debut": 0,
      "periode_fin": 0
    }
    return apart
  }, {})

  // on note le nombre d'heures demandées dans output_apart
  Object.keys(apdemande).forEach(hash => {
    var periode_deb = apdemande[hash].periode.start
    var periode_fin = apdemande[hash].periode.end

    // Des periodes arrondies aux débuts de périodes
    // TODO meilleur arrondi
    var periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0))
    var periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0))
    apart[apdemande[hash].id_demande.substring(0, 9)].periode_debut = periode_deb_floor
    apart[apdemande[hash].id_demande.substring(0, 9)].periode_fin = periode_fin_ceil

    var series = f.generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach( date => {
      let time = date.getTime()
      output_apart[time] = output_apart[time] || {}
      output_apart[time].apart_heures_autorisees = apdemande[hash].hta
    })
  })

  // relier les consos faites aux demandes (hashs) dans apart
  Object.keys(apconso).forEach(hash => {
    var valueap = apconso[hash]
    if (valueap.id_conso.substring(0, 9) in apart) {
      apart[valueap.id_conso.substring(0, 9)].consommation.push(hash)
    }
  })

  Object.keys(apart).forEach(k => {
    if (apart[k].consommation.length > 0) {
      apart[k].consommation.sort(
        (a,b) => (apconso[a].periode.getTime() >= apconso[b].periode.getTime())
      ).forEach( (h) => {
        var time = apconso[h].periode.getTime()
        output_apart[time] = output_apart[time] || {}
        output_apart[time].apart_heures_consommees = (output_apart[time].apart_heures_consommees || 0) + apconso[h].heure_consomme
        output_apart[time].apart_motif_recours = apdemande[apart[k].demande].motif_recours_se
      })

      // Heures consommees cumulees sur la demande
      let series = f.generatePeriodSerie(apart[k].periode_debut, apart[k].periode_fin)
      series.reduce( (accu, date) => {
        let time = date.getTime()

        //output_apart est déjà défini pour les heures autorisées
        accu = accu + (output_apart[time].apart_heures_consommees || 0)
        output_apart[time].apart_heures_consommees_cumulees = accu
        return(accu)
      }, 0)
    }
  })

  //Object.keys(output_apart).forEach(time => {
  //  if (output_effectif && time in output_effectif){
  //    output_apart[time].ratio_apart = (output_apart[time].apart_heures_consommees || 0) / (output_effectif[time].effectif * 157.67)
  //    //nbr approximatif d'heures ouvrées par mois
  //  }
  //})
  return(output_apart)
}`,
"ccsf": `function ccsf(v, output_array){
  "use strict";

  var ccsfHashes = Object.keys(v.ccsf || {})

  output_array.forEach(val => {
    var optccsf = ccsfHashes.reduce( function (accu, hash) {
      let ccsf = v.ccsf[hash]
      if (ccsf.date_traitement.getTime() < val.periode.getTime() && ccsf.date_traitement.getTime() > accu.date_traitement.getTime()) {
        let accu = ccsf
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
}`,
"cibleApprentissage": `function cibleApprentissage(output_indexed, n_months) {
  "use strict";

  // Mock two input instead of one for future modification
  var output_cotisation = output_indexed
  var output_procol = output_indexed
  // replace with const
  var all_keys = Object.keys(output_indexed)
  //

  var merged_info = all_keys.reduce(function(m, k) {
    m[k] = {outcome: Boolean(
      output_procol[k].tag_failure || output_cotisation[k].tag_default
    )}
    return m
  }, {})

  var output_outcome = f.lookAhead(merged_info, "outcome", n_months, true)
  var output_default = f.lookAhead(output_cotisation, "tag_default", n_months, true)
  var output_failure = f.lookAhead(output_procol, "tag_failure", n_months, true)

  var output_cible = all_keys.reduce(function(m, k) {
    m[k] = {}

    if (output_outcome[k])
      m[k] = output_outcome[k]
    if (output_default[k])
      m[k].time_til_default = output_default[k].time_til_outcome
    if (output_failure[k])
      m[k].time_til_failure = output_failure[k].time_til_outcome
    return m
  }, {})

  return output_cible
}`,
"compareDebit": `function compareDebit (a,b) {
  "use strict";
  if (a.numero_historique < b.numero_historique) return -1
  if (a.numero_historique > b.numero_historique) return 1
  return 0
}`,
"compte": `function compte (v, periodes) {
  "use strict";
  let output_compte = {}

  //  var offset_compte = 3
  Object.keys(v.compte).forEach(hash =>{
    var periode = v.compte[hash].periode.getTime()

    output_compte[periode] =  output_compte[periode] || {}
    output_compte[periode].compte_urssaf =  v.compte[hash].numero_compte
  })

  return output_compte
}`,
"cotisation": `function cotisation(output_indexed, output_array) {
  "use strict";
  // calcul de cotisation_moyenne sur 12 mois
  Object.keys(output_indexed).forEach(k => {
    let periode_courante = output_indexed[k].periode
    let periode_12_mois = f.dateAddMonth(periode_courante, 12)
    let series = f.generatePeriodSerie(periode_courante, periode_12_mois)
    series.forEach(periode => {
      if (periode.getTime() in output_indexed){
        if ("cotisation" in output_indexed[periode_courante.getTime()])
          output_indexed[periode.getTime()].cotisation_array = (output_indexed[periode.getTime()].cotisation_array || []).concat(output_indexed[periode_courante.getTime()].cotisation)

        output_indexed[periode.getTime()].montant_pp_array =
          (output_indexed[periode.getTime()].montant_pp_array || []).concat( output_indexed[periode_courante.getTime()].montant_part_patronale)
        output_indexed[periode.getTime()].montant_po_array =
          (output_indexed[periode.getTime()].montant_po_array || []).concat( output_indexed[periode_courante.getTime()].montant_part_ouvriere)
      }
    })
  })

  output_array.forEach(val => {
    val.cotisation_array = (val.cotisation_array || [] )
    val.cotisation_moy12m = val.cotisation_array.reduce( (p, c) => p + c, 0) / (val.cotisation_array.length || 1)
    if (val.cotisation_moy12m > 0) {
      val.ratio_dette = (val.montant_part_ouvriere + val.montant_part_patronale) / val.cotisation_moy12m
      let pp_average = (val.montant_pp_array || []).reduce((p, c) => p + c, 0) / (val.montant_pp_array.length || 1)
      let po_average =  (val.montant_po_array || []).reduce((p, c) => p + c, 0) / (val.montant_po_array.length || 1)
      val.ratio_dette_moy12m = (po_average + pp_average) / val.cotisation_moy12m
    }
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montant_pp_array || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montant_po_array || []).reduce((p, c) => (c >=
    //100) || p, false)
    delete val.cotisation_array
    delete val.montant_pp_array
    delete val.montant_po_array
  })

  // Calcul des défauts URSSAF prolongés
  var counter = 0
  Object.keys(output_indexed).sort().forEach(k => {
    if (output_indexed[k].ratio_dette > 0.01){
      output_indexed[k].tag_debit = true // Survenance d'un débit d'au moins 1% des cotisations
    }
    if (output_indexed[k].ratio_dette > 1){
      counter = counter + 1
      if (counter >= 3)
        output_indexed[k].tag_default = true
    } else
      counter = 0
  })
}`,
"cotisationsdettes": `function cotisationsdettes(v, periodes) {
  "use strict";

  // Tous les débits traitées après ce jour du mois sont reportées à la période suivante
  // Permet de s'aligner avec le calendrier de fourniture des données
  const last_treatment_day = 20

  var output_cotisationsdettes = {}

  // TODO Cotisations avec un mois de retard ? Bizarre, plus maintenant que l'export se fait le 20
  // var offset_cotisation = 1
  const offset_cotisation = 0
  var value_cotisation = {}

  // Répartition des cotisations sur toute la période qu'elle concerne
  Object.keys(v.cotisation).forEach(function (h) {
    var cotisation = v.cotisation[h]
    var periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = f.dateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })



  // relier les débits
  // ecn: ecart negatif
  // map les débits: clé fabriquée maison => [{hash, numero_historique, date_traitement}, ...]
  // Pour un même compte, les débits avec le même num_ecn (chaque émission de facture) sont donc regroupés
  var ecn = Object.keys(v.debit).reduce((accu, h) => {
      //pour chaque debit
      let debit = v.debit[h]

      var start = debit.periode.start
      var end = debit.periode.end
      var num_ecn = debit.numero_ecart_negatif
      var compte = debit.numero_compte
      var key = start + "-" + end + "-" + num_ecn + "-" + compte
      accu[key] = (accu[key] || []).concat([{
          "hash": h,
          "numero_historique": debit.numero_historique,
          "date_traitement": debit.date_traitement
      }])
      return accu
  }, {})

  // Pour chaque numero_ecn, on trie et on chaîne les débits avec debit_suivant
  Object.keys(ecn).forEach(i => {
      ecn[i].sort(f.compareDebit)
      var l = ecn[i].length
      ecn[i].forEach((e, idx) => {
          if (idx <= l - 2) {
              v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
          }
      })
  })

  var value_dette = {}
  // Pour chaque objet debit:
  // debit_traitement_debut => periode de traitement du débit
  // debit_traitement_fin => periode de traitement du debit suivant, ou bien date_fin
  // Entre ces deux dates, c'est cet objet qui est le plus à jour.
  Object.keys(v.debit).forEach(function (h) {
    var debit = v.debit[h]

    var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})

    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
    let jour_traitement = debit.date_traitement.getUTCDate()
    let jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate()
    let date_traitement_debut
    if (jour_traitement <= last_treatment_day){
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1)
      )
    }

    let date_traitement_fin
    if (jour_traitement_suivant <= last_treatment_day) {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + 1)
      )
    }

    let periode_debut = date_traitement_debut
    let periode_fin = date_traitement_fin

    //f.generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(periode_debut, periode_fin).map(date => {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })
  })

  // TODO faire numero de compte ailleurs
  // Array des numeros de compte
  //var numeros_compte = Array.from(new Set(
  //  Object.keys(v.cotisation).map(function (h) {
  //    return(v.cotisation[h].numero_compte)
  //  })
  //))

  periodes.forEach(function (time) {
    output_cotisationsdettes[time] = output_cotisationsdettes[time] || {}
    var val = output_cotisationsdettes[time]
  //output_cotisationsdettes[time].numero_compte_urssaf = numeros_compte
    if (time in value_cotisation){
      // somme de toutes les cotisations dues pour une periode donnée
      val.cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
    }

    // somme de tous les débits (part ouvriere, part patronale, montant_majorations)
    let montant_dette = (value_dette[time] || []).reduce(function (m, dette) {
      m.montant_part_ouvriere += dette.part_ouvriere
      m.montant_part_patronale += dette.part_patronale
      m.montant_majorations += dette.montant_majorations
      return m
    }, {"montant_part_ouvriere": 0, "montant_part_patronale": 0, "montant_majorations": 0})
    val = Object.assign(val, montant_dette)


    let past_month_offsets = [1,2,3,6,12]
    let time_d = new Date(parseInt(time))

    past_month_offsets.forEach(offset => {
      let time_offset = f.dateAddMonth(time_d, offset)
      let variable_name_part_ouvriere = "montant_part_ouvriere_past_" + offset
      let variable_name_part_patronale = "montant_part_patronale_past_" + offset
      output_cotisationsdettes[time_offset.getTime()] = output_cotisationsdettes[time_offset.getTime()] || {}
      let val_offset = output_cotisationsdettes[time_offset.getTime()]
      val_offset[variable_name_part_ouvriere] = val.montant_part_ouvriere
      val_offset[variable_name_part_patronale] = val.montant_part_patronale
    })

    let future_month_offsets = [0, 1, 2, 3, 4, 5]
    if (val.montant_part_ouvriere + val.montant_part_patronale > 0){
      future_month_offsets.forEach(offset => {
        let time_offset = f.dateAddMonth(time_d, offset)
        output_cotisationsdettes[time_offset.getTime()] = output_cotisationsdettes[time_offset.getTime()] || {}
        output_cotisationsdettes[time_offset.getTime()].interessante_urssaf = false
      })
    }
  })

  return(output_cotisationsdettes)
}`,
"dateAddMonth": `function dateAddMonth(date, nbMonth) {
  "use strict";
  var result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}`,
"dealWithProcols": `function dealWithProcols(data_source, altar_or_procol, output_indexed){
  "use strict";
  var codes  =  Object.keys(data_source).reduce((events, hash) => {
    var the_event = data_source[hash]

    if (altar_or_procol == "altares")
      var etat = f.altaresToHuman(the_event.code_evenement);
    else if (altar_or_procol == "procol")
      var etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);

    if (etat != null)
      events.push({"etat": etat, "date_proc_col": new Date(the_event.date_effet)})

    return(events)
  },[]).sort(
    (a,b) => {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
  )

  codes.forEach(
    event => {
      let periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
      var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})

      time_til_last.forEach(time => {
        if (time in output_indexed) {
          output_indexed[time].etat_proc_collective = event.etat
          output_indexed[time].date_proc_collective = event.date_proc_col
          if (event.etat != "in_bonis")
            output_indexed[time].tag_failure = true
        }
      })
    }
  )
}`,
"defaillances": `function defaillances (v, output_indexed) {
  "use strict";
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
  
  
  `,
"delais": `function delais(v, output_indexed) {
    "use strict";
    Object.keys(v.delai).map(function (hash) {
        const delai = v.delai[hash];
        // On arrondit les dates au premier jour du mois.
        const date_creation = new Date(Date.UTC(delai.date_creation.getUTCFullYear(), delai.date_creation.getUTCMonth(), 1, 0, 0, 0, 0));
        const date_echeance = new Date(Date.UTC(delai.date_echeance.getUTCFullYear(), delai.date_echeance.getUTCMonth(), 1, 0, 0, 0, 0));
        // Création d'un tableau de timestamps à raison de 1 par mois.
        const pastYearTimes = f
            .generatePeriodSerie(date_creation, date_echeance)
            .map(function (date) {
            return date.getTime();
        });
        pastYearTimes.map(function (time) {
            if (time in output_indexed) {
                const remaining_months = date_echeance.getUTCMonth() -
                    new Date(time).getUTCMonth() +
                    12 *
                        (date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear());
                output_indexed[time].delai = remaining_months;
                output_indexed[time].duree_delai = delai.duree_delai;
                output_indexed[time].montant_echeancier = delai.montant_echeancier;
                if (delai.duree_delai > 0) {
                    output_indexed[time].ratio_dette_delai =
                        (output_indexed[time].montant_part_patronale +
                            output_indexed[time].montant_part_ouvriere -
                            (delai.montant_echeancier * remaining_months * 30) /
                                delai.duree_delai) /
                            delai.montant_echeancier;
                }
            }
        });
    });
}`,
"detteFiscale": `function detteFiscale (diane){
  "use strict";
  if  (("dette_fiscale_et_sociale" in diane) && (diane["dette_fiscale_et_sociale"] !== null) &&
      ("valeur_ajoutee" in diane) && (diane["valeur_ajoutee"] !== null) &&
      (diane["valeur_ajoutee"] != 0)){
    return diane["dette_fiscale_et_sociale"]/ diane["valeur_ajoutee"] * 100
  } else {
    return null
  }
}`,
"effectifs": `function effectifs (effobj, periodes, effectif_name) {
  "use strict";

  let output_effectif = {}

  // Construction d'une map[time] = effectif à cette periode
  let map_effectif = Object.keys(effobj).reduce((m, hash) => {
    var effectif = effobj[hash]
    if (effectif == null) {
      return m
    }
    var effectifTime = effectif.periode.getTime()
    m[effectifTime] = (m[effectifTime] || 0) + effectif.effectif
    return m
  }, {})

  //ne reporter que si le dernier est disponible
  // 1- quelle periode doit être disponible
  var last_period = new Date(parseInt(periodes[periodes.length - 1]))
  var last_period_offset = f.dateAddMonth(last_period, offset_effectif + 1)
  // 2- Cette période est-elle disponible ?

  var available = map_effectif[last_period_offset.getTime()] ? 1 : 0


  //pour chaque periode (elles sont triees dans l'ordre croissant)
  periodes.reduce((accu, time) => {
    var periode = new Date(parseInt(time))
    // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
    output_effectif[time] = output_effectif[time] || {}
    output_effectif[time][effectif_name] = map_effectif[time] || (available ? accu : null)


    // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
    accu = map_effectif[time] || accu

    output_effectif[time][effectif_name + "_reporte"] = map_effectif[time] ? 0 : 1
    return(accu)
  }, null)

  Object.keys(map_effectif).forEach(time => {
    var periode = new Date(parseInt(time))
    var past_month_offsets = [6,12,18,24]
    past_month_offsets.forEach(lookback => {
      // On ajoute un offset pour partir de la dernière période où l'effectif est connu
      var time_past_lookback = f.dateAddMonth(periode, lookback - offset_effectif - 1)

      var variable_name_effectif = effectif_name + "_past_" + lookback
      output_effectif[time_past_lookback.getTime()] = output_effectif[time_past_lookback.getTime()] || {}
      output_effectif[time_past_lookback.getTime()][variable_name_effectif] = map_effectif[time]
    })
  })

  // On supprime les effectifs 'null'
  Object.keys(output_effectif).forEach(k => {
    if (output_effectif[k].effectif == null && output_effectif[k].effectif_ent == null) {
      delete output_effectif[k]
    }
  })
  return(output_effectif)
}`,
"finalize": `function finalize(k, v) {
  "use strict";
  const maxBsonSize = 16777216;

  // v de la forme
  // _id: {batch / siren / periode / type}
  // value: {siret1: {}, siret2: {}, "siren": {}}
  //
  ///
  ///////////////////////////////////////////////
  // consolidation a l'echelle de l'entreprise //
  ///////////////////////////////////////////////
  ///
  //

  let etablissements_connus = []
  let entreprise = (v.entreprise || {})

  Object.keys(v).forEach(siret =>{
    if (siret != "entreprise") {
      etablissements_connus[siret] = true
      if (v[siret].effectif){
        entreprise.effectif_entreprise = (entreprise.effectif_entreprise || 0) + v[siret].effectif // initialized to null
      }
      if (v[siret].apart_heures_consommees){
        entreprise.apart_entreprise = (entreprise.apart_entreprise || 0) + v[siret].apart_heures_consommees // initialized to 0
      }
      if (v[siret].montant_part_patronale || v[siret].montant_part_ouvriere){
        entreprise.debit_entreprise = (entreprise.debit_entreprise || 0) +
          (v[siret].montant_part_patronale || 0) +
          (v[siret].montant_part_ouvriere || 0)
      }
    }
  })


  Object.keys(v).forEach(siret =>{
    if (siret != "entreprise"){
      Object.assign(v[siret], entreprise)
    }
  })

  // une fois que les comptes sont faits...
  let output = []
  let nb_connus = Object.keys(etablissements_connus).length
  Object.keys(v).forEach(siret => {
    if (siret != "entreprise" && v[siret]) {
      v[siret].nbr_etablissements_connus = nb_connus
      output.push(v[siret])
    }
  })

  // NON: Pour l'instant, filtrage a posteriori
  // output = output.filter(siret_data => {
  //   return(siret_data.effectif) // Only keep if there is known effectif
  // })

  if (output.length > 0 && nb_connus <= 1500){
    if ((Object.bsonsize(output)  + Object.bsonsize({"_id": k})) < maxBsonSize){
      return output
    } else {
      print("Warning: my name is " + JSON.stringify(key, null, 2) + " and I died in reduce.algo2/finalize.js")
      return {"incomplete": true}
    }
  }
}`,
"financierCourtTerme": `function financierCourtTerme(diane) {
  "use strict";
  if  (("concours_bancaire_courant" in diane) && (diane["concours_bancaire_courant"] !== null) &&
    ("ca" in diane) && (diane["ca"] !== null) &&
    (diane["ca"] != 0)){
    return diane["concours_bancaire_courant"]/diane["ca"] * 100
  } else {
    return null
  }
}`,
"flatten": `function flatten(v, actual_batch) {
  "use strict";
  var res = Object.keys(v.batch || {})
    .sort()
    .filter(batch => batch <= actual_batch)
    .reduce((m, batch) => {

      // Types intéressants = nouveaux types, ou types avec suppressions
      var delete_types = Object.keys((v.batch[batch].compact || {}).delete || {})
      var new_types =  Object.keys(v.batch[batch])
      var all_interesting_types = [...new Set([...delete_types, ...new_types])]

      all_interesting_types.forEach(type => {
        m[type] = (m[type] || {})
        // On supprime les clés qu'il faut
        if (v.batch[batch] && v.batch[batch].compact && v.batch[batch].compact.delete &&
          v.batch[batch].compact.delete[type] && v.batch[batch].compact.delete[type] != {}) {

          v.batch[batch].compact.delete[type].forEach(hash => {
            delete m[type][hash]
          })
        }
        Object.assign(m[type], v.batch[batch][type])
      })
      return m
    }, { "key": v.key, scope: v.scope })

  return(res)
}`,
"fraisFinancier": `function fraisFinancier(diane) {
    "use strict";
    if ("interets" in diane &&
        diane["interets"] !== null &&
        "excedent_brut_d_exploitation" in diane &&
        diane["excedent_brut_d_exploitation"] !== null &&
        "produits_financiers" in diane &&
        diane["produits_financiers"] !== null &&
        "charges_financieres" in diane &&
        diane["charges_financieres"] !== null &&
        "charge_exceptionnelle" in diane &&
        diane["charge_exceptionnelle"] !== null &&
        "produit_exceptionnel" in diane &&
        diane["produit_exceptionnel"] !== null &&
        diane["excedent_brut_d_exploitation"] +
            diane["produits_financiers"] +
            diane["produit_exceptionnel"] -
            diane["charge_exceptionnelle"] -
            diane["charges_financieres"] !==
            0) {
        return ((diane["interets"] /
            (diane["excedent_brut_d_exploitation"] +
                diane["produits_financiers"] +
                diane["produit_exceptionnel"] -
                diane["charge_exceptionnelle"] -
                diane["charges_financieres"])) *
            100);
    }
    else {
        return null;
    }
}`,
"interim": `function interim (interim, output_indexed) {
  "use strict";
  let output_effectif = output_indexed
  // let periodes = Object.keys(output_indexed)
  // output_indexed devra être remplacé par output_effectif, et ne contenir que les données d'effectif.
  // periodes sera passé en argument.

  let output_interim = {}

  //  var offset_interim = 3

  Object.keys(interim).forEach(hash =>{
    var one_interim = interim[hash]
    var periode = one_interim.periode.getTime()
    // var periode_d = new Date(parseInt(interimTime))
    // var time_offset = f.dateAddMonth(time_d, -offset_interim)
    if (periode in output_effectif){
      output_interim[periode] = output_interim[periode] || {}
      output_interim[periode].interim_proportion = one_interim.etp / output_effectif[periode].effectif
    }

    var past_month_offsets = [6, 12, 18, 24]
    past_month_offsets.forEach(offset =>{
      var time_past_offset = f.dateAddMonth(one_interim.periode, offset)
      var variable_name_interim = "interim_ratio_past_" + offset
      if (periode in output_effectif && time_past_offset.getTime() in output_effectif){
        output_interim[time_past_offset.getTime()] =  output_interim[time_past_offset.getTime()] || {}
        var val_offset = output_interim[time_past_offset.getTime()]
        val_offset[variable_name_interim] = one_interim.etp  / output_effectif[periode].effectif
      }
    })
  })

  return output_interim
}`,
"lookAhead": `function lookAhead(data, attr_name, n_months, past) {
  "use strict";
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé
  var sorting_fun = function(a, b) { return(a >= b) }
  if (past) {
    sorting_fun = function(a, b) { return(a <= b) }
  }

  var counter = -1
  var output = Object.keys(data).sort(sorting_fun).reduce(function (m, period) {
    // Si on a déjà détecté quelque chose, on compte le nombre de périodes
    if (counter >= 0) counter = counter + 1

    if (data[period][attr_name]) {
      // si l'évènement se produit on retombe à 0
      counter = 0
    }

    if (counter >= 0) {
      // l'évènement s'est produit
      m[period] = m[period] || {}
      m[period].time_til_outcome = counter
      if (m[period].time_til_outcome <= n_months) {
        m[period].outcome = true
      } else {
        m[period].outcome = false
      }
    }
    return m
  }, {})

  return output
}`,
"map": `function map () {
  "use strict";
  let v = f.flatten(this.value, actual_batch)

  if (v.scope == "etablissement") {
    const o = f.outputs(v, serie_periode)
    let output_array = o[0] // [ OutputValue ] // in chronological order
    let output_indexed = o[1] // { Periode -> OutputValue } // OutputValue: cf outputs()

    // Les periodes qui nous interessent, triées
    var periodes = Object.keys(output_indexed).sort((a,b) => (a >= b))

    if (includes["apart"] || includes["all"]){
      if (v.apconso && v.apdemande) {
        let output_apart = f.apart(v.apconso, v.apdemande)
        Object.keys(output_apart).forEach(periode => {
          let data = {}
          data[this._id] = output_apart[periode]
          data[this._id].siret = this._id
          periode = new Date(Number(periode))
          emit(
            {
              'batch': actual_batch,
              'siren': this._id.substring(0, 9),
              'periode': periode,
              'type': 'apart'
            },
            data
          )
        })
      }
    }

    if (includes["all"]){

      if (v.compte) {
        var output_compte = f.compte(v)
        f.add(output_compte, output_indexed)
      }

      if (v.effectif) {
        var output_effectif = f.effectifs(v.effectif, periodes, "effectif")
        f.add(output_effectif, output_indexed)
      }

      if (v.interim){
        let output_interim = f.interim(v.interim, output_indexed)
        f.add(output_interim, output_indexed)
      }

      if (v.reporder){
        let output_repeatable = f.repeatable(v.reporder)
        f.add(output_repeatable, output_indexed)
      }

      if (v.delai) {f.delais(v, output_indexed)}

      v.altares = v.altares || {}
      v.procol = v.procol || {}

      if (v.altares) {
        f.defaillances(v, output_indexed)
      }

      if (v.cotisation && v.debit) {
        let output_cotisationsdettes = f.cotisationsdettes(v, periodes)
        f.add(output_cotisationsdettes, output_indexed)
      }

      if (v.ccsf) {f.ccsf(v, output_array)}
      if (v.sirene) {f.sirene(v, output_array)}

      f.populateNafAndApe(output_indexed, naf)

      f.cotisation(output_indexed, output_array)

      let output_cible = f.cibleApprentissage(output_indexed, 18)
      f.add(output_cible, output_indexed)
      output_array.forEach(val => {
        let data = {}
        data[this._id] = val
        emit(
          {
            'batch': actual_batch,
            'siren': this._id.substring(0, 9),
            'periode': val.periode,
            'type': 'other'
          },
          data
        )
      })
    }
  }

  if (v.scope == "entreprise") {

    if (includes["all"]){
      var output_array = serie_periode.map(function (e) {
        return {
          "siren": v.key,
          "periode": e,
          "exercice_bdf": 0,
          "arrete_bilan_bdf": new Date(0),
          "exercice_diane": 0,
          "arrete_bilan_diane": new Date(0)
        }
      })

      var output_indexed = output_array.reduce(function (periode, val) {
        periode[val.periode.getTime()] = val
        return periode
      }, {})

      if (v.sirene_ul) {f.sirene_ul(v, output_array)}

      var periodes = Object.keys(output_indexed).sort((a,b) => (a >= b))
      if (v.effectif_ent) {
        var output_effectif_ent = f.effectifs(v.effectif_ent, periodes, "effectif_ent")
        f.add(output_effectif_ent, output_indexed)
      }

      var output_indexed = output_array.reduce(function (periode, val) {
        periode[val.periode.getTime()] = val
        return periode
      }, {})

      v.bdf = (v.bdf || {})
      v.diane = (v.diane || {})

      Object.keys(v.bdf).forEach(hash => {
      }, {})

      v.bdf = (v.bdf || {})
      v.diane = (v.diane || {})

      Object.keys(v.bdf).forEach(hash => {
        let periode_arrete_bilan = new Date(Date.UTC(v.bdf[hash].arrete_bilan_bdf.getUTCFullYear(), v.bdf[hash].arrete_bilan_bdf.getUTCMonth() +1, 1, 0, 0, 0, 0))
        let periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7)
        let series = f.generatePeriodSerie(
          periode_dispo,
          f.dateAddMonth(periode_dispo, 13)
        )

        series.forEach(periode => {
          Object.keys(v.bdf[hash]).filter( k => {
            var omit = ["raison_sociale","secteur", "siren"]
            return (v.bdf[hash][k] != null &&  !(omit.includes(k)))
          }).forEach(k => {
            if (periode.getTime() in output_indexed){
              output_indexed[periode.getTime()][k] = v.bdf[hash][k]
              output_indexed[periode.getTime()].exercice_bdf = output_indexed[periode.getTime()].annee_bdf - 1
            }

            let past_year_offset = [1,2]
            past_year_offset.forEach( offset =>{
              let periode_offset = f.dateAddMonth(periode, 12* offset)
              let variable_name =  k + "_past_" + offset
              if (periode_offset.getTime() in output_indexed &&
                k != "arrete_bilan_bdf" &&
                k != "exercice_bdf"){
                output_indexed[periode_offset.getTime()][variable_name] = v.bdf[hash][k]
              }
            })
          })
        })
      })

      Object.keys(v.diane).filter(hash => v.diane[hash].arrete_bilan_diane).forEach(hash => {
        //v.diane[hash].arrete_bilan_diane = new Date(Date.UTC(v.diane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
        let periode_arrete_bilan = new Date(Date.UTC(v.diane[hash].arrete_bilan_diane.getUTCFullYear(), v.diane[hash].arrete_bilan_diane.getUTCMonth() +1, 1, 0, 0, 0, 0))
        let periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7) // 01/08 pour un bilan le 31/12, donc algo qui tourne en 01/09
        let series = f.generatePeriodSerie(
          periode_dispo,
          f.dateAddMonth(periode_dispo, 14) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (14+7)
        )

        series.forEach(periode => {
          Object.keys(v.diane[hash]).filter( k => {
            var omit = ["marquee", "nom_entreprise","numero_siren",
              "statut_juridique", "procedure_collective"]
            return (v.diane[hash][k] != null &&  !(omit.includes(k)))
          }).forEach(k => {
            if (periode.getTime() in output_indexed){
              output_indexed[periode.getTime()][k] = v.diane[hash][k]
            }

            // Passé

            let past_year_offset = [1,2]
            past_year_offset.forEach(offset =>{
              let periode_offset = f.dateAddMonth(periode, 12 * offset)
              let variable_name =  k + "_past_" + offset

              if (periode_offset.getTime() in output_indexed &&
                k != "arrete_bilan_diane" &&
                k != "exercice_diane"){
                output_indexed[periode_offset.getTime()][variable_name] = v.diane[hash][k]
              }
            })
          }
          )
        })

        series.forEach(periode => {
          if (periode.getTime() in output_indexed){
            // Recalcul BdF si ratios bdf sont absents
            if (!("poids_frng" in output_indexed[periode.getTime()]) && (f.poidsFrng(v.diane[hash]) !== null)){
              output_indexed[periode.getTime()].poids_frng = f.poidsFrng(v.diane[hash])
            }
            if (!("dette_fiscale" in output_indexed[periode.getTime()]) && (f.detteFiscale(v.diane[hash]) !== null)){
              output_indexed[periode.getTime()].dette_fiscale = f.detteFiscale(v.diane[hash])
            }
            if (!("frais_financier" in output_indexed[periode.getTime()]) && (f.fraisFinancier(v.diane[hash]) !== null)){
              output_indexed[periode.getTime()].frais_financier = f.fraisFinancier(v.diane[hash])
            }

            var bdf_vars = ["taux_marge", "poids_frng", "dette_fiscale", "financier_court_terme", "frais_financier"]
            let past_year_offset = [1,2]
            bdf_vars.forEach(k =>{
              if (k in output_indexed[periode.getTime()]){
                past_year_offset.forEach(offset =>{
                  let periode_offset = f.dateAddMonth(periode, 12 * offset)
                  let variable_name =  k + "_past_" + offset

                  if (periode_offset.getTime() in output_indexed){
                    output_indexed[periode_offset.getTime()][variable_name] = output_indexed[periode.getTime()][k]
                  }
                })
              }
            })
          }
        })
      })


      output_array.forEach((periode, index) => {
        if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0 && (periode.arrete_bilan_diane || new Date(0)).getTime() == 0) {
          delete output_array[index]
        }
        if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0){
          delete periode.arrete_bilan_bdf
        }
        if ((periode.arrete_bilan_diane||new Date(0)).getTime() == 0){
          delete periode.arrete_bilan_diane
        }

        emit(
          {
            'batch': actual_batch,
            'siren': this._id.substring(0, 9),
            'periode': periode.periode,
            'type': 'other'
          },
          {
            'entreprise': periode
          }
        )
      })
    }
  }
}`,
"outputs": `function outputs (v, serie_periode) {
  "use strict";
  var output_array = serie_periode.map(function (e) {
    return {
      "siret": v.key,
      "periode": e,
      "effectif": null,
      "etat_proc_collective": "in_bonis",
      "interessante_urssaf": true,
      "outcome": false
    }
  });

  var output_indexed = output_array.reduce(function (periodes, val) {
      periodes[val.periode.getTime()] = val
      return periodes
  }, {})

  return [output_array, output_indexed]
}`,
"poidsFrng": `function poidsFrng(diane){
  "use strict";
  if  (("couverture_ca_fdr" in diane) && (diane["couverture_ca_fdr"] !== null)){
    return diane["couverture_ca_fdr"]/360 * 100
  } else {
    return null
  }
}`,
"populateNafAndApe": `function populateNafAndApe(output_indexed, naf) {
  "use strict";
  Object.keys(output_indexed).forEach(k =>{
    if (("code_ape" in output_indexed[k]) && (output_indexed[k].code_ape !== null)){
      var code_ape = output_indexed[k].code_ape
      output_indexed[k].code_naf = naf.n5to1[code_ape]
      output_indexed[k].libelle_naf = naf.n1[output_indexed[k].code_naf]
      output_indexed[k].code_ape_niveau2 = code_ape.substring(0,2)
      output_indexed[k].code_ape_niveau3 = code_ape.substring(0,3)
      output_indexed[k].code_ape_niveau4 = code_ape.substring(0,4)
      output_indexed[k].libelle_ape2 = naf.n2[output_indexed[k].code_ape_niveau2]
      output_indexed[k].libelle_ape3 = naf.n3[output_indexed[k].code_ape_niveau3]
      output_indexed[k].libelle_ape4 = naf.n4[output_indexed[k].code_ape_niveau4]
      output_indexed[k].libelle_ape5 = naf.n5[code_ape]
    }
  })
}`,
"procolToHuman": `function procolToHuman (action, stade) {
  "use strict";
  var res = null;
  if (action == "liquidation" && stade != "abandon_procedure") 
    res = 'liquidation';
  else if (stade == "abandon_procedure" || stade == "fin_procedure")
    res = 'in_bonis';
  else if (action == "redressement" && stade == "plan_continuation")
    res = 'continuation';
  else if (action == "sauvegarde" && stade == "plan_continuation")
    res = 'sauvegarde';
  else if (action == "sauvegarde")
    res = 'plan_sauvegarde';
  else if (action == "redressement")
    res = 'plan_redressement';
  return res;
}`,
"reduce": `function reduce(key, values) {
  "use strict";
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {})
}`,
"repeatable": `function repeatable(rep){
  "use strict";
  let output_repeatable = {}
  Object.keys(rep).forEach(hash => {
    var one_rep = rep[hash]
    var periode = one_rep.periode.getTime()
    output_repeatable[periode] = output_repeatable[periode] || {}
    output_repeatable[periode].random_order = one_rep.random_order
  })

  return(output_repeatable)

}`,
"sirene": `function sirene (v, output_array) {
  "use strict";
  var sireneHashes = Object.keys(v.sirene || {})

  output_array.forEach(val => {
    // geolocalisation

    if (sireneHashes.length != 0) {
      var sirene = v.sirene[sireneHashes[sireneHashes.length - 1]]
      val.siren = val.siret.substring(0, 9)
      val.latitude = sirene.lattitude || null
      val.longitude = sirene.longitude || null
      val.departement = sirene.departement || null
      if (val.departement){
        val.region = f.region(val.departement)
      }
      var regexp_naf = /^[0-9]{4}[A-Z]$/
      if (sirene.ape && sirene.ape.match(regexp_naf)){
        val.code_ape  = sirene.ape
      }
      val.raison_sociale = sirene.raison_sociale || null
      // val.activite_saisonniere = sirene.activite_saisoniere || null
      // val.productif = sirene.productif || null
      // val.tranche_ca = sirene.tranche_ca || null
      // val.indice_monoactivite = sirene.indice_monoactivite || null
      val.date_creation_etablissement = sirene.date_creation ? sirene.date_creation.getFullYear() : null
      val.age = (sirene.date_creation && sirene.date_creation >= new Date("1901/01/01")) ? val.periode.getFullYear() - val.date_creation_etablissement : null
    }
  })
}`,
"sirene_ul": `function sirene_ul(v, output_array) {
  "use strict";
  var sireneHashes = Object.keys(v.sirene_ul || {})
  output_array.forEach(val => {
    if (sireneHashes.length != 0) {
      var sirene = v.sirene_ul[sireneHashes[sireneHashes.length - 1]]
      val.siren = val.siren
      val.raison_sociale = f.raison_sociale(
        sirene.raison_sociale,
        sirene.nom_unite_legale,
        sirene.nom_usage_unite_legale,
        sirene.prenom1_unite_legale,
        sirene.prenom2_unite_legale,
        sirene.prenom3_unite_legale,
        sirene.prenom4_unite_legale
      )
      val.statut_juridique = sirene.statut_juridique || null
      val.date_creation_entreprise = sirene.date_creation ? sirene.date_creation.getFullYear() : null
      val.age_entreprise = (sirene.date_creation && sirene.date_creation >= new Date("1901/01/01")) ? val.periode.getFullYear() - val.date_creation_entreprise : null
    }
  })
}`,
"tauxMarge": `function tauxMarge(diane) {
  "use strict";
  if  (("excedent_brut_d_exploitation" in diane) && (diane["excedent_brut_d_exploitation"] !== null) &&
    ("valeur_ajoutee" in diane) && (diane["valeur_ajoutee"] !== null) &&
    (diane["excedent_brut_d_exploitation"] != 0)){
    return diane["excedent_brut_d_exploitation"]/diane["valeur_ajoutee"] * 100
  } else {
    return null
  }
}`,
},
}
