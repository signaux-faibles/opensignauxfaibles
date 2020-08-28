package engine 

 var jsFunctions = map[string]map[string]string{
"common":{
"altaresToHuman": `function altaresToHuman(code) {
    "use strict";
    const codeLiquidation = [
        "PCL0108",
        "PCL010801",
        "PCL010802",
        "PCL030107",
        "PCL030307",
        "PCL030311",
        "PCL05010103",
        "PCL05010204",
        "PCL05010303",
        "PCL05010403",
        "PCL05010503",
        "PCL05010703",
        "PCL05011004",
        "PCL05011102",
        "PCL05011204",
        "PCL05011206",
        "PCL05011304",
        "PCL05011404",
        "PCL05011504",
        "PCL05011604",
        "PCL05011903",
        "PCL05012004",
        "PCL050204",
        "PCL0109",
        "PCL010901",
        "PCL030108",
        "PCL030308",
        "PCL05010104",
        "PCL05010205",
        "PCL05010304",
        "PCL05010404",
        "PCL05010504",
        "PCL05010803",
        "PCL05011005",
        "PCL05011103",
        "PCL05011205",
        "PCL05011207",
        "PCL05011305",
        "PCL05011405",
        "PCL05011505",
        "PCL05011904",
        "PCL05011605",
        "PCL05012005",
    ];
    const codePlanSauvegarde = [
        "PCL010601",
        "PCL0106",
        "PCL010602",
        "PCL030103",
        "PCL030303",
        "PCL03030301",
        "PCL05010101",
        "PCL05010202",
        "PCL05010301",
        "PCL05010401",
        "PCL05010501",
        "PCL05010506",
        "PCL05010701",
        "PCL05010705",
        "PCL05010801",
        "PCL05010805",
        "PCL05011002",
        "PCL05011202",
        "PCL05011302",
        "PCL05011402",
        "PCL05011502",
        "PCL05011602",
        "PCL05011901",
        "PCL0114",
        "PCL030110",
        "PCL030310",
    ];
    const codeRedressement = [
        "PCL0105",
        "PCL010501",
        "PCL010502",
        "PCL010503",
        "PCL030105",
        "PCL030305",
        "PCL05010102",
        "PCL05010203",
        "PCL05010302",
        "PCL05010402",
        "PCL05010502",
        "PCL05010702",
        "PCL05010706",
        "PCL05010802",
        "PCL05010806",
        "PCL05010901",
        "PCL05011003",
        "PCL05011101",
        "PCL05011203",
        "PCL05011303",
        "PCL05011403",
        "PCL05011503",
        "PCL05011603",
        "PCL05011902",
        "PCL05012003",
    ];
    const codeInBonis = [
        "PCL05",
        "PCL0501",
        "PCL050101",
        "PCL050102",
        "PCL050103",
        "PCL050104",
        "PCL050105",
        "PCL050106",
        "PCL050107",
        "PCL050108",
        "PCL050109",
        "PCL050110",
        "PCL050111",
        "PCL050112",
        "PCL050113",
        "PCL050114",
        "PCL050115",
        "PCL050116",
        "PCL050119",
        "PCL050120",
        "PCL050121",
        "PCL0503",
        "PCL050301",
        "PCL050302",
        "PCL0508",
        "PCL010504",
        "PCL010803",
        "PCL010902",
        "PCL050901",
        "PCL050902",
        "PCL050903",
        "PCL050904",
        "PCL0504",
        "PCL050303",
        "PCL050401",
        "PCL050402",
        "PCL050403",
        "PCL050404",
        "PCL050405",
        "PCL050406",
    ];
    const codeContinuation = ["PCL0202"];
    const codeSauvegarde = ["PCL0203", "PCL020301", "PCL0205", "PCL040408"];
    const codeCession = ["PCL0204", "PCL020401", "PCL020402", "PCL020403"];
    let res = null;
    if (codeLiquidation.includes(code))
        res = "liquidation";
    else if (codePlanSauvegarde.includes(code))
        res = "plan_sauvegarde";
    else if (codeRedressement.includes(code))
        res = "plan_redressement";
    else if (codeInBonis.includes(code))
        res = "in_bonis";
    else if (codeContinuation.includes(code))
        res = "continuation";
    else if (codeSauvegarde.includes(code))
        res = "sauvegarde";
    else if (codeCession.includes(code))
        res = "cession";
    return res;
}`,
"compareDebit": `function compareDebit(a, b) {
    "use strict";
    if (a.numero_historique < b.numero_historique)
        return -1;
    if (a.numero_historique > b.numero_historique)
        return 1;
    return 0;
}`,
"dateAddMonth": `function dateAddMonth(date, nbMonth) {
    "use strict";
    const result = new Date(date.getTime());
    result.setUTCMonth(result.getUTCMonth() + nbMonth);
    return result;
}`,
"flatten": `/**
 * Appelé par ` + "`" + `map()` + "`" + `, ` + "`" + `flatten()` + "`" + ` transforme les données importées (*Batches*)
 * d'une entreprise ou établissement afin de retourner un unique objet *plat*
 * contenant les valeurs finales de chaque type de données.
 *
 * Pour cela:
 * - il supprime les clés ` + "`" + `compact.delete` + "`" + ` des *Batches* en entrées;
 * - il agrège les propriétés apportées par chaque *Batch*, dans l'ordre chrono.
 */
function flatten(v, actual_batch) {
    "use strict";
    const res = Object.keys(v.batch || {})
        .sort()
        .filter((batch) => batch <= actual_batch)
        .reduce((m, batch) => {
        // Types intéressants = nouveaux types, ou types avec suppressions
        const delete_types = Object.keys((v.batch[batch].compact || {}).delete || {});
        const new_types = Object.keys(v.batch[batch]);
        const all_interesting_types = [
            ...new Set([...delete_types, ...new_types]),
        ];
        all_interesting_types.forEach((type) => {
            var _a, _b, _c;
            const typedData = m[type];
            if (typeof typedData === "object") {
                // On supprime les clés qu'il faut
                const keysToDelete = ((_c = (_b = (_a = v.batch[batch]) === null || _a === void 0 ? void 0 : _a.compact) === null || _b === void 0 ? void 0 : _b.delete) === null || _c === void 0 ? void 0 : _c[type]) || [];
                for (const hash of keysToDelete) {
                    delete typedData[hash];
                }
            }
            else {
                m[type] = {};
            }
            Object.assign(m[type], v.batch[batch][type]);
        });
        return m;
    }, { key: v.key, scope: v.scope });
    return res;
}`,
"forEachPopulatedProp": `// Appelle fct() pour chaque propriété définie (non undefined) de obj.
// Contrat: obj ne doit contenir que les clés définies dans son type.
function forEachPopulatedProp(obj, fct) {
    ;
    Object.keys(obj).forEach((key) => {
        if (typeof obj[key] !== "undefined")
            fct(key, obj[key]);
    });
}`,
"generatePeriodSerie": `function generatePeriodSerie(date_debut, date_fin) {
    "use strict";
    const date_next = new Date(date_debut.getTime());
    const serie = [];
    while (date_next.getTime() < date_fin.getTime()) {
        serie.push(new Date(date_next.getTime()));
        date_next.setUTCMonth(date_next.getUTCMonth() + 1);
    }
    return serie;
}`,
"omit": `// Fonction pour omettre des props, tout en retournant le bon type
function omit(object, ...propNames) {
    const result = Object.assign({}, object);
    for (const prop of propNames) {
        delete result[prop];
    }
    return result;
}`,
"procolToHuman": `function procolToHuman(action, stade) {
    "use strict";
    let res = null;
    if (action === "liquidation" && stade !== "abandon_procedure")
        res = "liquidation";
    else if (stade === "abandon_procedure" || stade === "fin_procedure")
        res = "in_bonis";
    else if (action === "redressement" && stade === "plan_continuation")
        res = "continuation";
    else if (action === "sauvegarde" && stade === "plan_continuation")
        res = "sauvegarde";
    else if (action === "sauvegarde")
        res = "plan_sauvegarde";
    else if (action === "redressement")
        res = "plan_redressement";
    return res;
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
"region": `function region(departement) {
    "use strict";
    const corr = {
        "01": "Auvergne-Rhône-Alpes",
        "03": "Auvergne-Rhône-Alpes",
        "07": "Auvergne-Rhône-Alpes",
        "15": "Auvergne-Rhône-Alpes",
        "26": "Auvergne-Rhône-Alpes",
        "38": "Auvergne-Rhône-Alpes",
        "42": "Auvergne-Rhône-Alpes",
        "43": "Auvergne-Rhône-Alpes",
        "63": "Auvergne-Rhône-Alpes",
        "69": "Auvergne-Rhône-Alpes",
        "73": "Auvergne-Rhône-Alpes",
        "74": "Auvergne-Rhône-Alpes",
        "02": "Hauts-de-France",
        "59": "Hauts-de-France",
        "60": "Hauts-de-France",
        "62": "Hauts-de-France",
        "80": "Hauts-de-France",
        "04": "Provence-Alpes-Côte d'Azur",
        "05": "Provence-Alpes-Côte d'Azur",
        "06": "Provence-Alpes-Côte d'Azur",
        "13": "Provence-Alpes-Côte d'Azur",
        "83": "Provence-Alpes-Côte d'Azur",
        "84": "Provence-Alpes-Côte d'Azur",
        "08": "Grand Est",
        "10": "Grand Est",
        "51": "Grand Est",
        "52": "Grand Est",
        "54": "Grand Est",
        "55": "Grand Est",
        "57": "Grand Est",
        "67": "Grand Est",
        "68": "Grand Est",
        "88": "Grand Est",
        "09": "Occitanie",
        "11": "Occitanie",
        "12": "Occitanie",
        "30": "Occitanie",
        "31": "Occitanie",
        "32": "Occitanie",
        "34": "Occitanie",
        "46": "Occitanie",
        "48": "Occitanie",
        "65": "Occitanie",
        "66": "Occitanie",
        "81": "Occitanie",
        "82": "Occitanie",
        "14": "Normandie",
        "27": "Normandie",
        "50": "Normandie",
        "61": "Normandie",
        "76": "Normandie",
        "18": "Centre-Val de Loire",
        "28": "Centre-Val de Loire",
        "36": "Centre-Val de Loire",
        "37": "Centre-Val de Loire",
        "41": "Centre-Val de Loire",
        "45": "Centre-Val de Loire",
        "16": "Nouvelle-Aquitaine",
        "17": "Nouvelle-Aquitaine",
        "19": "Nouvelle-Aquitaine",
        "23": "Nouvelle-Aquitaine",
        "24": "Nouvelle-Aquitaine",
        "33": "Nouvelle-Aquitaine",
        "40": "Nouvelle-Aquitaine",
        "47": "Nouvelle-Aquitaine",
        "64": "Nouvelle-Aquitaine",
        "79": "Nouvelle-Aquitaine",
        "86": "Nouvelle-Aquitaine",
        "87": "Nouvelle-Aquitaine",
        "20": "Corse",
        "21": "Bourgogne-Franche-Comté",
        "25": "Bourgogne-Franche-Comté",
        "39": "Bourgogne-Franche-Comté",
        "58": "Bourgogne-Franche-Comté",
        "70": "Bourgogne-Franche-Comté",
        "71": "Bourgogne-Franche-Comté",
        "89": "Bourgogne-Franche-Comté",
        "90": "Bourgogne-Franche-Comté",
        "22": "Bretagne",
        "29": "Bretagne",
        "35": "Bretagne",
        "56": "Bretagne",
        "44": "Pays de la Loire",
        "49": "Pays de la Loire",
        "53": "Pays de la Loire",
        "72": "Pays de la Loire",
        "85": "Pays de la Loire",
        "75": "Île-de-France",
        "77": "Île-de-France",
        "78": "Île-de-France",
        "91": "Île-de-France",
        "92": "Île-de-France",
        "93": "Île-de-France",
        "94": "Île-de-France",
        "95": "Île-de-France",
    };
    return corr[departement] || "";
}`,
},
"compact":{
"applyPatchesToBatch": `function applyPatchesToBatch(hashToAdd, hashToDelete, stockTypes, currentBatch) {
    var _a;
    // Application des suppressions
    stockTypes.forEach((type) => {
        const hashesToDelete = hashToDelete[type];
        if (hashesToDelete) {
            currentBatch.compact = currentBatch.compact || { delete: {} };
            currentBatch.compact.delete = currentBatch.compact.delete || {};
            currentBatch.compact.delete[type] = [...hashesToDelete];
        }
    });
    // Application des ajouts
    forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
        currentBatch[type] = [...hashesToAdd].reduce((typedBatchValues, hash) => {
            var _a;
            return (Object.assign(Object.assign({}, typedBatchValues), { [hash]: (_a = currentBatch[type]) === null || _a === void 0 ? void 0 : _a[hash] }));
        }, {});
    });
    // Retrait des propriété vides
    // - compact.delete vides
    const compactDelete = (_a = currentBatch.compact) === null || _a === void 0 ? void 0 : _a.delete;
    if (compactDelete) {
        forEachPopulatedProp(compactDelete, (type, keysToDelete) => {
            if (keysToDelete.length === 0) {
                delete compactDelete[type];
            }
        });
        if (Object.keys(compactDelete).length === 0) {
            delete currentBatch.compact;
        }
    }
    // - types vides
    forEachPopulatedProp(currentBatch, (type, typedBatchData) => {
        if (Object.keys(typedBatchData).length === 0) {
            delete currentBatch[type];
        }
    });
}`,
"applyPatchesToMemory": `function applyPatchesToMemory(hashToAdd, hashToDelete, memory) {
    // Prise en compte des suppressions de clés dans la mémoire
    forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
        hashesToDelete.forEach((hash) => {
            memory[type].delete(hash);
        });
    });
    // Prise en compte des ajouts de clés dans la mémoire
    forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
        hashesToAdd.forEach((hash) => {
            memory[type] = memory[type] || new Set();
            memory[type].add(hash);
        });
    });
}`,
"compactBatch": `/**
 * Appelée par reduce(), compactBatch() va générer un diff entre les
 * données de batch et les données précédentes fournies par memory.
 * Paramètres modifiés: currentBatch et memory.
 * Pré-requis: les batches précédents doivent avoir été compactés.
 */
function compactBatch(currentBatch, memory, fromBatchKey) {
    // Les types où il y a potentiellement des suppressions
    const stockTypes = completeTypes[fromBatchKey].filter((type) => (memory[type] || new Set()).size > 0);
    const { hashToAdd, hashToDelete } = listHashesToAddAndDelete(currentBatch, stockTypes, memory);
    fixRedundantPatches(hashToAdd, hashToDelete, memory);
    applyPatchesToMemory(hashToAdd, hashToDelete, memory);
    applyPatchesToBatch(hashToAdd, hashToDelete, stockTypes, currentBatch);
    return currentBatch;
}`,
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
                delete reporder[ro];
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
            siret,
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
    // Retourne les clés de obj, en respectant le type défini dans le type de obj.
    // Contrat: obj ne doit contenir que les clés définies dans son type.
    const typedObjectKeys = (obj) => Object.keys(obj);
    const currentState = batches.reduce((m, batch) => {
        //1. On supprime les clés de la mémoire
        if (batch.compact) {
            forEachPopulatedProp(batch.compact.delete, (type, keysToDelete) => {
                keysToDelete.forEach((key) => {
                    m[type].delete(key); // Should never fail or collection is corrupted
                });
            });
        }
        //2. On ajoute les nouvelles clés
        for (const type of typedObjectKeys(batch)) {
            if (type === "compact")
                continue;
            m[type] = m[type] || new Set();
            for (const key in batch[type]) {
                m[type].add(key);
            }
        }
        return m;
    }, {});
    return currentState;
}`,
"finalize": `// finalize permet de:
// - indiquer les établissements à inclure dans les calculs de variables
// (processus reduce.algo2)
// - intégrer les reporder pour permettre la reproductibilité de
// l'échantillonnage pour l'entraînement du modèle.
function finalize(k, companyDataValues) {
    "use strict";
    let o = Object.assign(Object.assign({}, companyDataValues), { index: { algo1: false, algo2: false } });
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
"fixRedundantPatches": `/**
 * Modification de hashToAdd et hashToDelete pour retirer les redondances.
 **/
function fixRedundantPatches(hashToAdd, hashToDelete, memory) {
    forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
        // Pour chaque cle supprimee: est-ce qu'elle est bien dans la
        // memoire ? sinon on la retire de la liste des clés supprimées (pas de
        // maj memoire)
        // -----------------------------------------------------------------------------------------------------------------
        hashToDelete[type] = new Set([...hashesToDelete].filter((hash) => {
            return (memory[type] || new Set()).has(hash);
        }));
        // Est-ce qu'elle a ete egalement ajoutee en même temps que
        // supprimée ? (par exemple remplacement d'un stock complet à
        // l'identique) Dans ce cas là, on retire cette clé des valeurs ajoutées
        // et supprimées
        // i.e. on herite de la memoire. (pas de maj de la memoire)
        // ------------------------------------------------------------------------------
        hashToDelete[type] = new Set([...(hashToDelete[type] || new Set())].filter((hash) => {
            const hashesToAdd = hashToAdd[type] || new Set();
            const also_added = hashesToAdd.has(hash);
            if (also_added) {
                hashesToAdd.delete(hash);
            }
            return !also_added;
        }));
    });
    forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
        // Pour chaque cle ajoutee: est-ce qu'elle est dans la memoire ? Si oui on filtre cette cle
        // i.e. on herite de la memoire. (pas de maj de la memoire)
        // ---------------------------------------------------------------------------------------------
        hashToAdd[type] = new Set([...hashesToAdd].filter((hash) => {
            return !(memory[type] || new Set()).has(hash);
        }));
    });
}`,
"listHashesToAddAndDelete": `/**
 * On recupère les clés ajoutées et les clés supprimées depuis currentBatch.
 * On ajoute aux clés supprimées les types stocks de la memoire.
 */
function listHashesToAddAndDelete(currentBatch, stockTypes, memory) {
    const hashToDelete = {};
    const hashToAdd = {};
    // Itération sur les types qui ont potentiellement subi des modifications
    // pour compléter hashToDelete et hashToAdd.
    // Les suppressions de types complets / stock sont gérés dans le bloc suivant.
    forEachPopulatedProp(currentBatch, (type) => {
        var _a;
        // Le type compact gère les clés supprimées
        // Ce type compact existe si le batch en cours a déjà été compacté.
        if (type === "compact") {
            const compactDelete = (_a = currentBatch.compact) === null || _a === void 0 ? void 0 : _a.delete;
            if (compactDelete) {
                forEachPopulatedProp(compactDelete, (deleteType, keysToDelete) => {
                    keysToDelete.forEach((hash) => {
                        ;
                        (hashToDelete[deleteType] =
                            hashToDelete[deleteType] || new Set()).add(hash);
                    });
                });
            }
        }
        else {
            for (const hash in currentBatch[type]) {
                ;
                (hashToAdd[type] = hashToAdd[type] || new Set()).add(hash);
            }
        }
    });
    stockTypes.forEach((type) => {
        hashToDelete[type] = new Set([
            ...(hashToDelete[type] || new Set()),
            ...memory[type],
        ]);
    });
    return {
        hashToAdd,
        hashToDelete,
    };
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
function reduce(key, values // chaque element contient plusieurs batches pour cette entreprise ou établissement
) {
    "use strict";
    // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
    // batchs. Sinon, juste fusion des attributs
    const auxBatchSet = new Set();
    const severalBatches = values.some((value) => {
        auxBatchSet.add(Object.keys(value.batch || {}));
        return auxBatchSet.size > 1;
    });
    // Fusion batch par batch des types de données sans se préoccuper des doublons.
    const naivelyMergedCompanyData = values.reduce((m, value) => {
        Object.keys(value.batch).forEach((batch) => {
            m.batch[batch] = Object.keys(value.batch[batch]).reduce((batchValues, type) => (Object.assign(Object.assign({}, batchValues), { [type]: value.batch[batch][type] })), m.batch[batch] || {});
        });
        return m;
    }, { key, scope: values[0].scope, batch: {} });
    // Cette fonction reduce() est appelée à deux moments:
    // 1. agregation par établissement d'objets ImportedData. Dans cet étape, on
    // ne travaille généralement que sur un seul batch.
    // 2. agregation de ces résultats au sein de RawData, en fusionnant avec les
    // données potentiellement présentes. Dans cette étape, on fusionne
    // généralement les données de plusieurs batches. (données historiques)
    if (!severalBatches)
        return naivelyMergedCompanyData;
    //////////////////////////////////////////////////
    // ETAPES DE LA FUSION AVEC DONNÉES HISTORIQUES //
    //////////////////////////////////////////////////
    // 0. On calcule la memoire au moment du batch à modifier
    const memoryBatches = Object.keys(naivelyMergedCompanyData.batch)
        .filter((batch) => batch < fromBatchKey)
        .sort()
        .reduce((m, batch) => {
        m.push(naivelyMergedCompanyData.batch[batch]);
        return m;
    }, []);
    // Memory conserve les données aplaties de tous les batches jusqu'à fromBatchKey
    // puis sera enrichie au fur et à mesure du traitement des batches suivants.
    const memory = f.currentState(memoryBatches);
    const reducedValue = {
        key: naivelyMergedCompanyData.key,
        scope: naivelyMergedCompanyData.scope,
        batch: {},
    };
    // Copie telle quelle des batches jusqu'à fromBatchKey.
    Object.keys(naivelyMergedCompanyData.batch)
        .filter((batch) => batch < fromBatchKey)
        .forEach((batch) => {
        reducedValue.batch[batch] = naivelyMergedCompanyData.batch[batch];
    });
    // On itère sur chaque batch à partir de fromBatchKey pour les compacter.
    // Il est possible qu'il y ait moins de batch en sortie que le nombre traité
    // dans la boucle, si ces batchs n'apportent aucune information nouvelle.
    batches
        .filter((batch) => batch >= fromBatchKey)
        .forEach((batch) => {
        const currentBatch = naivelyMergedCompanyData.batch[batch];
        const compactedBatch = compactBatch(currentBatch, memory, batch);
        if (Object.keys(compactedBatch).length > 0) {
            reducedValue.batch[batch] = compactedBatch;
        }
    });
    return reducedValue;
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
    return f
        .iterable(apconso)
        .sort((p1, p2) => (p1.periode < p2.periode ? 1 : -1));
}`,
"apdemande": `function apdemande(apdemande) {
    return f
        .iterable(apdemande)
        .sort((p1, p2) => p1.periode.start.getTime() < p2.periode.start.getTime() ? 1 : -1);
}`,
"bdf": `function bdf(hs) {
    "use strict";
    const bdf = {};
    // Déduplication par arrete_bilan_bdf
    f.iterable(hs)
        .filter((b) => b.arrete_bilan_bdf)
        .forEach((b) => {
        bdf[b.arrete_bilan_bdf.toISOString()] = b;
    });
    return f
        .iterable(bdf)
        .sort((a, b) => (a.annee_bdf < b.annee_bdf ? 1 : -1));
}`,
"compte": `function compte(compte) {
    const c = f.iterable(compte);
    return c.length > 0 ? c[c.length - 1] : undefined;
}`,
"cotisations": `function cotisations(vcotisation = {}) {

    const offset_cotisation = 0;
    const value_cotisation = {};
    // Répartition des cotisations sur toute la période qu'elle concerne
    Object.keys(vcotisation).forEach(function (h) {
        const cotisation = vcotisation[h];
        const periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end);
        periode_cotisation.forEach((date_cotisation) => {
            const date_offset = f.dateAddMonth(date_cotisation, offset_cotisation);
            value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat([cotisation.du / periode_cotisation.length]);
        });
    });
    const output_cotisation = [];
    serie_periode.forEach((p) => {
        output_cotisation.push((value_cotisation[p.getTime()] || []).reduce((m, c) => m + c, 0));
    });
    return output_cotisation;
}`,
"dateAddDay": `function dateAddDay(date, nbDays) {
    "use strict";
    const result = new Date(date.getTime());
    result.setDate(result.getDate() + nbDays);
    return result;
}`,
"dealWithProcols": `function dealWithProcols(data_source = {}, altar_or_procol) {

    return Object.keys(data_source)
        .reduce((events, hash) => {
        const the_event = data_source[hash];
        let etat = null;
        if (altar_or_procol === "altares")
            etat = f.altaresToHuman(the_event.code_evenement);
        else if (altar_or_procol === "procol")
            etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);
        if (etat !== null)
            events.push({ etat, date_procol: new Date(the_event.date_effet) });
        return events;
    }, [])
        .sort((a, b) => a.date_procol.getTime() - b.date_procol.getTime());
}`,
"debits": `function debits(vdebit = {}) {

    const last_treatment_day = 20;
    const ecn = Object.keys(vdebit).reduce((accu, h) => {
        const debit = vdebit[h];
        const start = debit.periode.start;
        const end = debit.periode.end;
        const num_ecn = debit.numero_ecart_negatif;
        const compte = debit.numero_compte;
        const key = start + "-" + end + "-" + num_ecn + "-" + compte;
        accu[key] = (accu[key] || []).concat([
            {
                hash: h,
                numero_historique: debit.numero_historique,
                date_traitement: debit.date_traitement,
            },
        ]);
        return accu;
    }, {});
    Object.keys(ecn).forEach((i) => {
        ecn[i].sort(f.compareDebit);
        const l = ecn[i].length;
        ecn[i].forEach((e, idx) => {
            if (idx <= l - 2) {
                vdebit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
            }
        });
    });
    const value_dette = {};
    Object.keys(vdebit).forEach(function (h) {
        const debit = vdebit[h];
        const debit_suivant = vdebit[debit.debit_suivant] || {
            date_traitement: date_fin,
        };
        //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
        const jour_traitement = debit.date_traitement.getUTCDate();
        const jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate();
        let date_traitement_debut;
        if (jour_traitement <= last_treatment_day) {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1));
        }
        let date_traitement_fin;
        if (jour_traitement_suivant <= last_treatment_day) {
            date_traitement_fin = new Date(Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_fin = new Date(Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + 1));
        }
        const periode_debut = date_traitement_debut;
        const periode_fin = date_traitement_fin;
        //generatePeriodSerie exlue la dernière période
        f.generatePeriodSerie(periode_debut, periode_fin).map((date) => {
            const time = date.getTime();
            value_dette[time] = (value_dette[time] || []).concat([
                {
                    periode: debit.periode.start,
                    part_ouvriere: debit.part_ouvriere,
                    part_patronale: debit.part_patronale,
                    montant_majorations: debit.montant_majorations || 0,
                },
            ]);
        });
    });
    return serie_periode.map((p) => (value_dette[p.getTime()] || []).reduce((m, c) => {
        m.part_ouvriere += c.part_ouvriere;
        m.part_patronale += c.part_patronale;
        m.montant_majorations += c.montant_majorations;
        return m;
    }, {
        part_ouvriere: 0,
        part_patronale: 0,
        montant_majorations: 0,
        periode: f.dateAddDay(f.dateAddMonth(p, 1), -1),
    }));
}`,
"delai": `function delai(delai) {
    return f.iterable(delai);
}`,
"diane": `function diane(hs) {
    "use strict";
    const diane = {};
    // Déduplication par arrete_bilan_diane
    f.iterable(hs)
        .filter((d) => d.arrete_bilan_diane)
        .forEach((d) => {
        diane[d.arrete_bilan_diane.toISOString()] = d;
    });
    return f
        .iterable(diane)
        .sort((a, b) => (a.exercice_diane < b.exercice_diane ? 1 : -1));
}`,
"effectifs": `function effectifs(effectif) {
    const mapEffectif = {};
    f.iterable(effectif).forEach((e) => {
        mapEffectif[e.periode.getTime()] =
            (mapEffectif[e.periode.getTime()] || 0) + e.effectif;
    });
    return serie_periode
        .map((p) => {
        return {
            periode: p,
            effectif: mapEffectif[p.getTime()] || -1,
        };
    })
        .filter((p) => p.effectif >= 0);
}`,
"finalize": `function finalize(_key, val) {
    return val;
}`,
"iterable": `function iterable(dict) {
    return typeof dict === "object" ? Object.keys(dict).map((h) => dict[h]) : [];
}`,
"joinUrssaf": `function joinUrssaf(effectif, debit) {
    const result = {
        effectif: [],
        part_patronale: [],
        part_ouvriere: [],
        montant_majorations: [],
    };
    debit.forEach((d, i) => {
        const e = effectif.filter((e) => serie_periode[i].getTime() === e.periode.getTime());
        if (e.length > 0) {
            result.effectif.push(e[0].effectif);
        }
        else {
            result.effectif.push(null);
        }
        result.part_patronale.push(d.part_patronale);
        result.part_ouvriere.push(d.part_ouvriere);
        result.montant_majorations.push(d.montant_majorations);
    });
    return result;
}`,
"map": `function map() {

    const value = f.flatten(this.value, actual_batch);
    if (this.value.scope === "etablissement") {
        const vcmde = {};
        vcmde.key = this.value.key;
        vcmde.batch = actual_batch;
        vcmde.sirene = f.sirene(f.iterable(value.sirene));
        vcmde.periodes = serie_periode;
        const effectif = f.effectifs(value.effectif);
        const debit = f.debits(value.debit);
        const join = f.joinUrssaf(effectif, debit);
        vcmde.debit_part_patronale = join.part_patronale;
        vcmde.debit_part_ouvriere = join.part_ouvriere;
        vcmde.debit_montant_majorations = join.montant_majorations;
        vcmde.effectif = join.effectif;
        vcmde.cotisation = f.cotisations(value.cotisation);
        vcmde.apconso = f.apconso(value.apconso);
        vcmde.apdemande = f.apdemande(value.apdemande);
        vcmde.delai = f.delai(value.delai);
        vcmde.compte = f.compte(value.compte);
        vcmde.procol = undefined; // Note: initialement, l'expression ci-dessous était affectée à vcmde.procol, puis écrasée plus bas. J'initialise quand même vcmde.procol ici pour ne pas faire échouer test-api.sh sur l'ordre des propriétés.
        const procol = [
            ...f.dealWithProcols(value.altares, "altares"),
            ...f.dealWithProcols(value.procol, "procol"),
        ];
        vcmde.last_procol = procol[procol.length - 1] || { etat: "in_bonis" };
        vcmde.idEntreprise = "entreprise_" + this.value.key.slice(0, 9);
        vcmde.procol = value.procol;
        emit("etablissement_" + this.value.key, vcmde);
    }
    else if (this.value.scope === "entreprise") {
        const v = {};
        const diane = f.diane(value.diane);
        const bdf = f.bdf(value.bdf);
        const sirene_ul = (value.sirene_ul || {})[Object.keys(value.sirene_ul || {})[0] || ""];
        const crp = value.crp;
        v.key = this.value.key;
        v.batch = actual_batch;
        if (diane.length > 0) {
            v.diane = diane;
        }
        if (bdf.length > 0) {
            v.bdf = bdf;
        }
        if (sirene_ul) {
            v.sirene_ul = sirene_ul;
        }
        if (crp) {
            v.crp = crp;
        }
        if (Object.keys(v) !== []) {
            emit("entreprise_" + this.value.key, v);
        }
    }
}`,
"reduce": `function reduce(_key, values) {
    return values.reduce((m, v) => {
        if (v.sirets) {
            // TODO: je n'ai pas trouvé d'affectation de valeur dans la propriété "sirets" => est-elle toujours d'actualité ?
            m.sirets = (m.sirets || []).concat(v.sirets);
            delete v.sirets;
        }
        Object.assign(m, v);
        return m;
    }, {});
}`,
"sirene": `// Cette fonction retourne les données sirene les plus récentes
function sirene(sireneArray) {
    return sireneArray[sireneArray.length - 1] || {}; // TODO: vérifier que sireneArray est bien classé dans l'ordre chronologique
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
"add": `function add(obj, output) {
    "use strict";
    Object.keys(output).forEach(function (periode) {
        if (periode in obj) {
            Object.assign(output[periode], obj[periode]);
        }
    });
}`,
"apart": `function apart(apconso, apdemande) {
    "use strict";
    const output_apart = {};
    // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
    const apart = Object.keys(apdemande).reduce((apart, hash) => {
        apart[apdemande[hash].id_demande.substring(0, 9)] = {
            demande: hash,
            consommation: [],
            periode_debut: new Date(0),
            periode_fin: new Date(0),
        };
        return apart;
    }, {});
    // on note le nombre d'heures demandées dans output_apart
    Object.keys(apdemande).forEach((hash) => {
        const periode_deb = apdemande[hash].periode.start;
        const periode_fin = apdemande[hash].periode.end;
        // Des periodes arrondies aux débuts de périodes
        // TODO meilleur arrondi
        const periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0));
        const periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        apart[apdemande[hash].id_demande.substring(0, 9)].periode_debut = periode_deb_floor;
        apart[apdemande[hash].id_demande.substring(0, 9)].periode_fin = periode_fin_ceil;
        const series = f.generatePeriodSerie(periode_deb_floor, periode_fin_ceil);
        series.forEach((date) => {
            const time = date.getTime();
            output_apart[time] = output_apart[time] || {};
            output_apart[time].apart_heures_autorisees = apdemande[hash].hta;
        });
    });
    // relier les consos faites aux demandes (hashs) dans apart
    Object.keys(apconso).forEach((hash) => {
        const valueap = apconso[hash];
        if (valueap.id_conso.substring(0, 9) in apart) {
            apart[valueap.id_conso.substring(0, 9)].consommation.push(hash);
        }
    });
    Object.keys(apart).forEach((k) => {
        if (apart[k].consommation.length > 0) {
            apart[k].consommation
                .sort((a, b) => apconso[a].periode.getTime() - apconso[b].periode.getTime())
                .forEach((h) => {
                const time = apconso[h].periode.getTime();
                output_apart[time] = output_apart[time] || {};
                output_apart[time].apart_heures_consommees =
                    (output_apart[time].apart_heures_consommees || 0) +
                        apconso[h].heure_consomme;
                output_apart[time].apart_motif_recours =
                    apdemande[apart[k].demande].motif_recours_se;
            });
            // Heures consommees cumulees sur la demande
            const series = f.generatePeriodSerie(apart[k].periode_debut, apart[k].periode_fin);
            series.reduce((accu, date) => {
                const time = date.getTime();
                //output_apart est déjà défini pour les heures autorisées
                accu = accu + (output_apart[time].apart_heures_consommees || 0);
                output_apart[time].apart_heures_consommees_cumulees = accu;
                return accu;
            }, 0);
        }
    });
    //Object.keys(output_apart).forEach(time => {
    //  if (output_effectif && time in output_effectif){
    //    output_apart[time].ratio_apart = (output_apart[time].apart_heures_consommees || 0) / (output_effectif[time].effectif * 157.67)
    //    //nbr approximatif d'heures ouvrées par mois
    //  }
    //})
    return output_apart;
}`,
"ccsf": `function ccsf(vCcsf, output_array) {
    "use strict";
    const ccsfHashes = Object.keys(vCcsf || {});
    output_array.forEach((val) => {
        const optccsf = ccsfHashes.reduce(function (accu, hash) {
            const ccsf = vCcsf[hash];
            if (ccsf.date_traitement.getTime() < val.periode.getTime() &&
                ccsf.date_traitement.getTime() > accu.date_traitement.getTime()) {
                return ccsf;
            }
            return accu;
        }, {
            date_traitement: new Date(0),
        });
        if (optccsf.date_traitement.getTime() !== 0) {
            val.date_ccsf = optccsf.date_traitement;
        }
    });
}`,
"cibleApprentissage": `function cibleApprentissage(output_indexed, n_months) {
    "use strict";
    // Mock two input instead of one for future modification
    const output_cotisation = output_indexed;
    const output_procol = output_indexed;
    // replace with const
    const all_keys = Object.keys(output_indexed);
    const merged_info = all_keys.reduce(function (m, k) {
        m[k] = {
            outcome: Boolean(output_procol[k].tag_failure || output_cotisation[k].tag_default),
        };
        return m;
    }, {});
    const output_outcome = f.lookAhead(merged_info, "outcome", n_months, true);
    const output_default = f.lookAhead(output_cotisation, "tag_default", n_months, true);
    const output_failure = f.lookAhead(output_procol, "tag_failure", n_months, true);
    const output_cible = all_keys.reduce(function (m, k) {
        const outputTimes = {};
        if (output_default[k])
            outputTimes.time_til_default = output_default[k].time_til_outcome;
        if (output_failure[k])
            outputTimes.time_til_failure = output_failure[k].time_til_outcome;
        return Object.assign(Object.assign({}, m), { [k]: Object.assign(Object.assign({}, output_outcome[k]), outputTimes) });
    }, {});
    return output_cible;
}`,
"compte": `function compte(compte) {
    "use strict";
    const output_compte = {};
    //  var offset_compte = 3
    Object.keys(compte).forEach((hash) => {
        const periode = compte[hash].periode.getTime().toString();
        output_compte[periode] = output_compte[periode] || {};
        output_compte[periode].compte_urssaf = compte[hash].numero_compte;
    });
    return output_compte;
}`,
"cotisation": `function cotisation(output_indexed) {
    "use strict";
    const sortieCotisation = {};

    const moyenne = (valeurs = []) => valeurs.some((val) => typeof val === "undefined")
        ? undefined
        : valeurs.reduce((p, c) => p + c, 0) / (valeurs.length || 1);
    // calcul de cotisation_moyenne sur 12 mois
    const futureArrays = {};
    Object.keys(output_indexed).forEach((periode) => {
        const input = output_indexed[periode];
        const périodeCourante = output_indexed[periode].periode;
        const douzeMoisÀVenir = f
            .generatePeriodSerie(périodeCourante, f.dateAddMonth(périodeCourante, 12))
            .map((periodeFuture) => ({ timestamp: periodeFuture.getTime() }))
            .filter(({ timestamp }) => timestamp in output_indexed);
        // Accumulation de cotisations sur les 12 mois à venir, pour calcul des moyennes
        douzeMoisÀVenir.forEach(({ timestamp }) => {
            const future = (futureArrays[timestamp] = futureArrays[timestamp] || {
                cotisations: [],
                montantsPP: [],
                montantsPO: [],
            });
            future.cotisations.push(input.cotisation);
            future.montantsPP.push(input.montant_part_patronale || 0);
            future.montantsPO.push(input.montant_part_ouvriere || 0);
        });
        // Calcul des cotisations moyennes à partir des valeurs accumulées ci-dessus
        const { cotisations, montantsPO, montantsPP } = futureArrays[periode];
        const out = (sortieCotisation[periode] = sortieCotisation[periode] || {});
        if (cotisations.length >= 12) {
            out.cotisation_moy12m = moyenne(cotisations);
        }
        if (typeof out.cotisation_moy12m === "undefined") {
            delete out.cotisation_moy12m;
        }
        else if (out.cotisation_moy12m > 0) {
            out.ratio_dette =
                ((input.montant_part_ouvriere || 0) +
                    (input.montant_part_patronale || 0)) /
                    out.cotisation_moy12m;
            if (!cotisations.includes(undefined) && !cotisations.includes(0)) {
                out.ratio_dette_moy12m = moyenne(montantsPO.map((_, i) => (montantsPO[i] + montantsPP[i]) / cotisations[i]));
            }
        }
        // Remplace dans cibleApprentissage
        //val.dette_any_12m = (val.montantsPA || []).reduce((p,c) => (c >=
        //100) || p, false) || (val.montantsPO || []).reduce((p, c) => (c >=
        //100) || p, false)
    });
    // Calcul des défauts URSSAF prolongés
    let counter = 0;
    Object.keys(sortieCotisation)
        .sort()
        .forEach((k) => {
        const { ratio_dette } = sortieCotisation[k];
        if (!ratio_dette)
            return;
        if (ratio_dette > 0.01) {
            sortieCotisation[k].tag_debit = true; // Survenance d'un débit d'au moins 1% des cotisations
        }
        if (ratio_dette > 1) {
            counter = counter + 1;
            if (counter >= 3)
                sortieCotisation[k].tag_default = true;
        }
        else
            counter = 0;
    });
    return sortieCotisation;
}`,
"cotisationsdettes": `/**
 * Calcule les variables liées aux cotisations sociales et dettes sur ces
 * cotisations.
 */
function cotisationsdettes(vCotisation, vDebit, periodes, finPériode // correspond à la variable globale date_fin
) {
    "use strict";

    // Tous les débits traitées après ce jour du mois sont reportées à la période suivante
    // Permet de s'aligner avec le calendrier de fourniture des données
    const lastAccountedDay = 20;
    const sortieCotisationsDettes = {};
    const value_cotisation = {};
    // Répartition des cotisations sur toute la période qu'elle concerne
    Object.keys(vCotisation).forEach(function (h) {
        const cotisation = vCotisation[h];
        const periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end);
        periode_cotisation.forEach((date_cotisation) => {
            value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat([cotisation.du / periode_cotisation.length]);
        });
    });
    // relier les débits
    // ecn: ecart negatif
    // map les débits: clé fabriquée maison => [{hash, numero_historique, date_traitement}, ...]
    // Pour un même compte, les débits avec le même num_ecn (chaque émission de facture) sont donc regroupés
    const ecn = Object.keys(vDebit).reduce((accu, h) => {
        //pour chaque debit
        const debit = vDebit[h];
        const start = debit.periode.start;
        const end = debit.periode.end;
        const num_ecn = debit.numero_ecart_negatif;
        const compte = debit.numero_compte;
        const key = start + "-" + end + "-" + num_ecn + "-" + compte;
        accu[key] = (accu[key] || []).concat([
            {
                hash: h,
                numero_historique: debit.numero_historique,
                date_traitement: debit.date_traitement,
            },
        ]);
        return accu;
    }, {});
    // Pour chaque numero_ecn, on trie et on chaîne les débits avec debit_suivant
    Object.keys(ecn).forEach((i) => {
        ecn[i].sort(f.compareDebit);
        const l = ecn[i].length;
        ecn[i].forEach((e, idx) => {
            if (idx <= l - 2) {
                vDebit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
            }
        });
    });
    const value_dette = {};
    // Pour chaque objet debit:
    // debit_traitement_debut => periode de traitement du débit
    // debit_traitement_fin => periode de traitement du debit suivant, ou bien finPériode
    // Entre ces deux dates, c'est cet objet qui est le plus à jour.
    Object.keys(vDebit).forEach(function (h) {
        const debit = vDebit[h];
        const debit_suivant = vDebit[debit.debit_suivant] || {
            date_traitement: finPériode,
        };
        //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
        const jour_traitement = debit.date_traitement.getUTCDate();
        const jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate();
        let date_traitement_debut;
        if (jour_traitement <= lastAccountedDay) {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1));
        }
        let date_traitement_fin;
        if (jour_traitement_suivant <= lastAccountedDay) {
            date_traitement_fin = new Date(Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_fin = new Date(Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + 1));
        }
        const periode_debut = date_traitement_debut;
        const periode_fin = date_traitement_fin;
        //f.generatePeriodSerie exlue la dernière période
        f.generatePeriodSerie(periode_debut, periode_fin).map((date) => {
            const time = date.getTime();
            value_dette[time] = (value_dette[time] || []).concat([
                {
                    periode: debit.periode.start,
                    part_ouvriere: debit.part_ouvriere,
                    part_patronale: debit.part_patronale,
                },
            ]);
        });
    });
    // TODO faire numero de compte ailleurs
    // Array des numeros de compte
    //var numeros_compte = Array.from(new Set(
    //  Object.keys(vCotisation).map(function (h) {
    //    return(vCotisation[h].numero_compte)
    //  })
    //))
    periodes.forEach(function (time) {
        sortieCotisationsDettes[time] = sortieCotisationsDettes[time] || {};
        let val = sortieCotisationsDettes[time];
        //output_cotisationsdettes[time].numero_compte_urssaf = numeros_compte
        if (time in value_cotisation) {
            // somme de toutes les cotisations dues pour une periode donnée
            val.cotisation = value_cotisation[time].reduce((a, cot) => a + cot, 0);
        }
        // somme de tous les débits (part ouvriere, part patronale)
        const montant_dette = (value_dette[time] || []).reduce(function (m, dette) {
            m.montant_part_ouvriere += dette.part_ouvriere;
            m.montant_part_patronale += dette.part_patronale;
            return m;
        }, {
            montant_part_ouvriere: 0,
            montant_part_patronale: 0,
        });
        val = Object.assign(val, montant_dette);
        const futureTimestamps = [1, 2, 3, 6, 12] // Penser à mettre à jour le type CotisationsDettesPassees pour tout changement
            .map((offset) => ({
            offset,
            timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
        }))
            .filter(({ timestamp }) => periodes.includes(timestamp));
        futureTimestamps.forEach(({ offset, timestamp }) => {
            sortieCotisationsDettes[timestamp] = Object.assign(Object.assign({}, sortieCotisationsDettes[timestamp]), { ["montant_part_ouvriere_past_" + offset]: val.montant_part_ouvriere, ["montant_part_patronale_past_" + offset]: val.montant_part_patronale });
        });
        if (val.montant_part_ouvriere + val.montant_part_patronale > 0) {
            const futureTimestamps = [0, 1, 2, 3, 4, 5]
                .map((offset) => ({
                timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
            }))
                .filter(({ timestamp }) => periodes.includes(timestamp));
            futureTimestamps.forEach(({ timestamp }) => {
                sortieCotisationsDettes[timestamp] = Object.assign(Object.assign({}, sortieCotisationsDettes[timestamp]), { interessante_urssaf: false });
            });
        }
    });
    return sortieCotisationsDettes;
}`,
"dealWithProcols": `function dealWithProcols(data_source, altar_or_procol, output_indexed) {
    "use strict";

    const codes = Object.keys(data_source)
        .reduce((events, hash) => {
        const the_event = data_source[hash];
        let etat = null;
        if (altar_or_procol === "altares")
            etat = f.altaresToHuman(the_event.code_evenement);
        else if (altar_or_procol === "procol")
            etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);
        if (etat !== null)
            events.push({
                etat,
                date_proc_col: new Date(the_event.date_effet),
            });
        return events;
    }, [])
        .sort((a, b) => {
        return a.date_proc_col.getTime() - b.date_proc_col.getTime();
    });
    codes.forEach((event) => {
        const periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0));
        const time_til_last = Object.keys(output_indexed).filter((val) => {
            return val >= periode_effet.toString();
        });
        time_til_last.forEach((time) => {
            if (time in output_indexed) {
                output_indexed[time].etat_proc_collective = event.etat;
                output_indexed[time].date_proc_collective = event.date_proc_col;
                if (event.etat !== "in_bonis")
                    output_indexed[time].tag_failure = true;
            }
        });
    });
}`,
"defaillances": `function defaillances(altares, procol, output_indexed) {
    "use strict";
    f.dealWithProcols(altares, "altares", output_indexed);
    f.dealWithProcols(procol, "procol", output_indexed);
}`,
"delais": `/**
 * Calcule pour chaque période le nombre de jours restants du délai accordé et
 * un indicateur de la déviation par rapport à un remboursement linéaire du
 * montant couvert par le délai. Un "délai" étant une demande accordée de délai
 * de paiement des cotisations sociales, pour un certain montant
 * (delai_montant_echeancier) et pendant une certaine période
 * (delai_nb_jours_total).
 * Contrat: cette fonction ne devrait être appelée que s'il y a eu au moins une
 * demande de délai.
 */
function delais(vDelai, debitParPériode, intervalleTraitement) {
    "use strict";
    const donnéesDélaiParPériode = {};
    Object.keys(vDelai).forEach(function (hash) {
        const delai = vDelai[hash];
        if (delai.duree_delai <= 0) {
            return;
        }
        // On arrondit les dates au premier jour du mois.
        const date_creation = new Date(Date.UTC(delai.date_creation.getUTCFullYear(), delai.date_creation.getUTCMonth(), 1, 0, 0, 0, 0));
        const date_echeance = new Date(Date.UTC(delai.date_echeance.getUTCFullYear(), delai.date_echeance.getUTCMonth(), 1, 0, 0, 0, 0));
        // Création d'un tableau de timestamps à raison de 1 par mois.
        f.generatePeriodSerie(date_creation, date_echeance)
            .filter((date) => date >= intervalleTraitement.premièreDate &&
            date <= intervalleTraitement.dernièreDate)
            .map(function (debutDeMois) {
            const time = debutDeMois.getTime();
            const remainingDays = nbDays(debutDeMois, delai.date_echeance);
            const inputAtTime = debitParPériode[time];
            const outputAtTime = {
                delai_nb_jours_restants: remainingDays,
                delai_nb_jours_total: delai.duree_delai,
                delai_montant_echeancier: delai.montant_echeancier,
            };
            if (typeof (inputAtTime === null || inputAtTime === void 0 ? void 0 : inputAtTime.montant_part_patronale) !== "undefined" &&
                typeof (inputAtTime === null || inputAtTime === void 0 ? void 0 : inputAtTime.montant_part_ouvriere) !== "undefined") {
                const detteActuelle = inputAtTime.montant_part_patronale +
                    inputAtTime.montant_part_ouvriere;
                const detteHypothétiqueRemboursementLinéaire = (delai.montant_echeancier * remainingDays) / delai.duree_delai;
                outputAtTime.delai_deviation_remboursement =
                    (detteActuelle - detteHypothétiqueRemboursementLinéaire) /
                        delai.montant_echeancier;
            }
            donnéesDélaiParPériode[time] = outputAtTime;
        });
    });
    return donnéesDélaiParPériode;
}`,
"detteFiscale": `function detteFiscale(diane) {
    "use strict";
    var _a, _b;
    const ratio = ((_a = diane["dette_fiscale_et_sociale"]) !== null && _a !== void 0 ? _a : NaN) /
        ((_b = diane["valeur_ajoutee"]) !== null && _b !== void 0 ? _b : NaN);
    return isNaN(ratio) ? null : ratio * 100;
}`,
"effectifs": `function effectifs(effobj, periodes, propertyName) {
    "use strict";
    const output_effectif = {};
    // Construction d'une map[time] = effectif à cette periode
    const map_effectif = Object.keys(effobj).reduce((m, hash) => {
        const effectif = effobj[hash];
        if (effectif === null) {
            return m;
        }
        const effectifTime = effectif.periode.getTime();
        m[effectifTime] = (m[effectifTime] || 0) + effectif.effectif;
        return m;
    }, {});
    //ne reporter que si le dernier est disponible
    // 1- quelle periode doit être disponible
    const last_period = new Date(periodes[periodes.length - 1]);
    const last_period_offset = f.dateAddMonth(last_period, offset_effectif + 1);
    // 2- Cette période est-elle disponible ?
    const available = last_period_offset.getTime() in map_effectif;
    //pour chaque periode (elles sont triees dans l'ordre croissant)
    periodes.reduce((accu, time) => {
        // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
        output_effectif[time] = output_effectif[time] || {};
        output_effectif[time][propertyName] =
            map_effectif[time] || (available ? accu : null);
        // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
        accu = map_effectif[time] || accu;
        Object.assign(output_effectif[time], {
            [propertyName + "_reporte"]: map_effectif[time] ? 0 : 1,
        });
        return accu;
    }, null);
    Object.keys(map_effectif).forEach((time) => {
        const periode = new Date(parseInt(time));
        const past_month_offsets = [6, 12, 18, 24]; // Note: à garder en synchro avec la définition du type PastPropertyName
        past_month_offsets.forEach((lookback) => {
            // On ajoute un offset pour partir de la dernière période où l'effectif est connu
            const time_past_lookback = f.dateAddMonth(periode, lookback - offset_effectif - 1);
            output_effectif[time_past_lookback.getTime()] =
                output_effectif[time_past_lookback.getTime()] || {};
            Object.assign(output_effectif[time_past_lookback.getTime()], {
                [propertyName + "_past_" + lookback]: map_effectif[time],
            });
        });
    });
    // On supprime les effectifs 'null'
    Object.keys(output_effectif).forEach((k) => {
        if (output_effectif[k].effectif === null &&
            output_effectif[k].effectif_ent === null) {
            delete output_effectif[k];
        }
    });
    return output_effectif;
}
/* TODO: appliquer même logique d'itération sur futureTimestamps que dans cotisationsdettes.ts */`,
"entr_bdf": `function entr_bdf(donnéesBdf, periodes) {
    "use strict";

    const outputBdf = {};
    for (const p of periodes) {
        outputBdf[p] = {};
    }
    for (const hash of Object.keys(donnéesBdf)) {
        const entréeBdf = donnéesBdf[hash];
        const periode_arrete_bilan = new Date(Date.UTC(entréeBdf.arrete_bilan_bdf.getUTCFullYear(), entréeBdf.arrete_bilan_bdf.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7);
        const series = f.generatePeriodSerie(periode_dispo, f.dateAddMonth(periode_dispo, 13));
        for (const periode of series) {
            const outputInPeriod = (outputBdf[periode.getTime()] =
                outputBdf[periode.getTime()] || {});
            const periodData = f.omit(entréeBdf, "raison_sociale", "secteur", "siren");
            // TODO: Éviter d'ajouter des données en dehors de ` + "`" + `periodes` + "`" + `, sans fausser le calcul des données passées (plus bas)
            Object.assign(outputInPeriod, periodData);
            if (outputInPeriod.annee_bdf) {
                outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1;
            }
            const pastData = f.omit(periodData, "arrete_bilan_bdf", "exercice_bdf");
            for (const prop of Object.keys(pastData)) {
                const past_year_offset = [1, 2];
                for (const offset of past_year_offset) {
                    const periode_offset = f.dateAddMonth(periode, 12 * offset);
                    const outputInPast = outputBdf[periode_offset.getTime()];
                    if (outputInPast) {
                        Object.assign(outputInPast, {
                            [prop + "_past_" + offset]: donnéesBdf[hash][prop],
                        });
                    }
                }
            }
        }
    }
    return outputBdf;
}`,
"entr_diane": `function entr_diane(donnéesDiane, output_indexed, periodes) {

    for (const hash of Object.keys(donnéesDiane)) {
        if (!donnéesDiane[hash].arrete_bilan_diane)
            continue;
        //donnéesDiane[hash].arrete_bilan_diane = new Date(Date.UTC(donnéesDiane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
        const periode_arrete_bilan = new Date(Date.UTC(donnéesDiane[hash].arrete_bilan_diane.getUTCFullYear(), donnéesDiane[hash].arrete_bilan_diane.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7); // 01/08 pour un bilan le 31/12, donc algo qui tourne en 01/09
        const series = f.generatePeriodSerie(periode_dispo, f.dateAddMonth(periode_dispo, 14) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (14+7)
        );
        for (const periode of series) {
            const rest = f.omit(donnéesDiane[hash], "marquee", "nom_entreprise", "numero_siren", "statut_juridique", "procedure_collective");
            if (periodes.includes(periode.getTime())) {
                Object.assign(output_indexed[periode.getTime()], rest);
            }
            for (const ratio of Object.keys(rest)) {
                if (donnéesDiane[hash][ratio] === null) {
                    if (periodes.includes(periode.getTime())) {
                        delete output_indexed[periode.getTime()][ratio];
                    }
                    continue;
                }
                // Passé
                const past_year_offset = [1, 2];
                for (const offset of past_year_offset) {
                    const periode_offset = f.dateAddMonth(periode, 12 * offset);
                    const variable_name = ratio + "_past_" + offset;
                    if (periode_offset.getTime() in output_indexed &&
                        ratio !== "arrete_bilan_diane" &&
                        ratio !== "exercice_diane") {
                        output_indexed[periode_offset.getTime()][variable_name] =
                            donnéesDiane[hash][ratio];
                    }
                }
            }
        }
        for (const periode of series) {
            if (periodes.includes(periode.getTime())) {
                // Recalcul BdF si ratios bdf sont absents
                const inputInPeriod = output_indexed[periode.getTime()];
                const outputInPeriod = output_indexed[periode.getTime()];
                if (!("poids_frng" in inputInPeriod)) {
                    const poids = f.poidsFrng(donnéesDiane[hash]);
                    if (poids !== null)
                        outputInPeriod.poids_frng = poids;
                }
                if (!("dette_fiscale" in inputInPeriod)) {
                    const dette = f.detteFiscale(donnéesDiane[hash]);
                    if (dette !== null)
                        outputInPeriod.dette_fiscale = dette;
                }
                if (!("frais_financier" in inputInPeriod)) {
                    const frais = f.fraisFinancier(donnéesDiane[hash]);
                    if (frais !== null)
                        outputInPeriod.frais_financier = frais;
                }
                // TODO: mettre en commun population des champs _past_ avec bdf ?
                const bdf_vars = [
                    "taux_marge",
                    "poids_frng",
                    "dette_fiscale",
                    "financier_court_terme",
                    "frais_financier",
                ];
                const past_year_offset = [1, 2];
                bdf_vars.forEach((k) => {
                    if (k in outputInPeriod) {
                        past_year_offset.forEach((offset) => {
                            const periode_offset = f.dateAddMonth(periode, 12 * offset);
                            const variable_name = k + "_past_" + offset;
                            if (periodes.includes(periode_offset.getTime())) {
                                output_indexed[periode_offset.getTime()][variable_name] =
                                    outputInPeriod[k];
                            }
                        });
                    }
                });
            }
        }
    }
    return output_indexed;
}`,
"entr_sirene": `function entr_sirene(sirene_ul, sériePériode) {
    "use strict";
    const retourEntrSirene = {};
    const sireneHashes = Object.keys(sirene_ul || {});
    sériePériode.forEach((période) => {
        if (sireneHashes.length !== 0) {
            const val = {};
            const sirene = sirene_ul[sireneHashes[sireneHashes.length - 1]];
            val.raison_sociale = f.raison_sociale(sirene.raison_sociale, sirene.nom_unite_legale, sirene.nom_usage_unite_legale, sirene.prenom1_unite_legale, sirene.prenom2_unite_legale, sirene.prenom3_unite_legale, sirene.prenom4_unite_legale);
            val.statut_juridique = sirene.statut_juridique || null;
            val.date_creation_entreprise = sirene.date_creation
                ? sirene.date_creation.getFullYear()
                : null;
            if (val.date_creation_entreprise &&
                sirene.date_creation &&
                sirene.date_creation >= new Date("1901/01/01")) {
                val.age_entreprise =
                    période.getFullYear() - val.date_creation_entreprise;
            }
            retourEntrSirene[période.getTime()] = val;
        }
    });
    return retourEntrSirene;
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
    // extraction de l'entreprise et des établissements depuis v
    const etab = f.omit(v, "entreprise");
    const entr = Object.assign({}, v.entreprise);
    const output = Object.keys(etab).map((siret) => {
        const { effectif } = etab[siret];
        if (effectif) {
            entr.effectif_entreprise = entr.effectif_entreprise || 0 + effectif;
        }
        const { apart_heures_consommees } = etab[siret];
        if (apart_heures_consommees) {
            entr.apart_entreprise =
                (entr.apart_entreprise || 0) + apart_heures_consommees;
        }
        if (etab[siret].montant_part_patronale ||
            etab[siret].montant_part_ouvriere) {
            entr.debit_entreprise =
                (entr.debit_entreprise || 0) +
                    (etab[siret].montant_part_patronale || 0) +
                    (etab[siret].montant_part_ouvriere || 0);
        }
        return Object.assign(Object.assign(Object.assign({}, etab[siret]), entr), { nbr_etablissements_connus: Object.keys(etab).length });
    });
    // NON: Pour l'instant, filtrage a posteriori
    // output = output.filter(siret_data => {
    //   return(siret_data.effectif) // Only keep if there is known effectif
    // })
    if (output.length > 0 && output.length <= 1500) {
        if (bsonsize(output) + bsonsize({ _id: k }) < maxBsonSize) {
            return output;
        }
        else {
            print("Warning: my name is " +
                JSON.stringify(k, null, 2) +
                " and I died in reduce.algo2/finalize.js");
            return { incomplete: true };
        }
    }
    else {
        return []; // ajouté pour résoudre erreur TS7030 (Not all code paths return a value)
    }
}`,
"financierCourtTerme": `function financierCourtTerme(diane) {
    "use strict";
    var _a, _b;
    const ratio = ((_a = diane["concours_bancaire_courant"]) !== null && _a !== void 0 ? _a : NaN) / ((_b = diane["ca"]) !== null && _b !== void 0 ? _b : NaN);
    return isNaN(ratio) ? null : ratio * 100;
}`,
"fraisFinancier": `function fraisFinancier(diane) {
    "use strict";
    var _a, _b, _c, _d, _e, _f;
    const ratio = ((_a = diane["interets"]) !== null && _a !== void 0 ? _a : NaN) /
        (((_b = diane["excedent_brut_d_exploitation"]) !== null && _b !== void 0 ? _b : NaN) +
            ((_c = diane["produits_financiers"]) !== null && _c !== void 0 ? _c : NaN) +
            ((_d = diane["produit_exceptionnel"]) !== null && _d !== void 0 ? _d : NaN) -
            ((_e = diane["charge_exceptionnelle"]) !== null && _e !== void 0 ? _e : NaN) -
            ((_f = diane["charges_financieres"]) !== null && _f !== void 0 ? _f : NaN));
    return isNaN(ratio) ? null : ratio * 100;
}`,
"interim": `function interim(interim, output_indexed) {
    "use strict";
    const output_effectif = output_indexed;
    // let periodes = Object.keys(output_indexed)
    // output_indexed devra être remplacé par output_effectif, et ne contenir que les données d'effectif.
    // periodes sera passé en argument.
    const output_interim = {};
    //  var offset_interim = 3
    Object.keys(interim).forEach((hash) => {
        const one_interim = interim[hash];
        const periode = one_interim.periode.getTime();
        // var periode_d = new Date(parseInt(interimTime))
        // var time_offset = f.dateAddMonth(time_d, -offset_interim)
        if (periode in output_effectif) {
            output_interim[periode] = output_interim[periode] || {};
            const { effectif } = output_effectif[periode];
            if (effectif) {
                output_interim[periode].interim_proportion = one_interim.etp / effectif;
            }
        }
        const past_month_offsets = [6, 12, 18, 24]; // En cas de changement, penser à mettre à jour le type SortieInterim
        past_month_offsets.forEach((offset) => {
            const time_past_offset = f.dateAddMonth(one_interim.periode, offset);
            if (periode in output_effectif &&
                time_past_offset.getTime() in output_effectif) {
                output_interim[time_past_offset.getTime()] =
                    output_interim[time_past_offset.getTime()] || {};
                const val_offset = output_interim[time_past_offset.getTime()];
                const { effectif } = output_effectif[periode];
                if (effectif) {
                    Object.assign(val_offset, {
                        [` + "`" + `interim_ratio_past_${offset}` + "`" + `]: one_interim.etp / effectif,
                    });
                }
            }
        });
    });
    return output_interim;
}`,
"lookAhead": `function lookAhead(data, attr_name, n_months, past) {
    "use strict";
    // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
    // demander: que va-t-il se passer) ou dans le future (past = false on
    // pourra se demander que s'est-il passé
    const chronologic = (a, b) => (a > b ? 1 : -1);
    const reverse = (a, b) => (b > a ? 1 : -1);
    let counter = -1;
    const output = Object.keys(data)
        .sort(past ? reverse : chronologic)
        .reduce(function (m, period) {
        // Si on a déjà détecté quelque chose, on compte le nombre de périodes
        if (counter >= 0)
            counter = counter + 1;
        if (data[period][attr_name]) {
            // si l'évènement se produit on retombe à 0
            counter = 0;
        }
        if (counter >= 0) {
            // l'évènement s'est produit
            m[period] = m[period] || {};
            m[period].time_til_outcome = counter;
            if (m[period].time_til_outcome <= n_months) {
                m[period].outcome = true;
            }
            else {
                m[period].outcome = false;
            }
        }
        return m;
    }, {});
    return output;
}`,
"map": `/**
 * ` + "`" + `map()` + "`" + ` est appelée pour chaque entreprise/établissement.
 *
 * Une entreprise/établissement est rattachée à des données de plusieurs types,
 * groupées par *Batch* (groupements de fichiers de données importés).
 *
 * Pour chaque période d'un entreprise/établissement, un objet contenant toutes
 * les données agrégées est émis (par appel à ` + "`" + `emit()` + "`" + `), à destination de
 * ` + "`" + `reduce()` + "`" + `, puis de ` + "`" + `finalize()` + "`" + `.
 */
function map() {
    "use strict";

    const v = f.flatten(this.value, actual_batch);
    if (v.scope === "etablissement") {
        const [output_array, // DonnéesAgrégées[] dans l'ordre chronologique
        output_indexed,] = f.outputs(v, serie_periode);
        // Les periodes qui nous interessent, triées
        const periodes = Object.keys(output_indexed)
            .sort()
            .map((timestamp) => parseInt(timestamp));
        if (includes["apart"] || includes["all"]) {
            if (v.apconso && v.apdemande) {
                const output_apart = f.apart(v.apconso, v.apdemande);
                Object.keys(output_apart).forEach((periode) => {
                    const data = {
                        [this._id]: Object.assign(Object.assign({}, output_apart[periode]), { siret: this._id }),
                    };
                    emit({
                        batch: actual_batch,
                        siren: this._id.substring(0, 9),
                        periode: new Date(Number(periode)),
                        type: "apart",
                    }, data);
                });
            }
        }
        if (includes["all"]) {
            if (v.compte) {
                const output_compte = f.compte(v.compte);
                f.add(output_compte, output_indexed);
            }
            if (v.effectif) {
                const output_effectif = f.effectifs(v.effectif, periodes, "effectif");
                f.add(output_effectif, output_indexed);
            }
            if (v.interim) {
                const output_interim = f.interim(v.interim, output_indexed);
                f.add(output_interim, output_indexed);
            }
            if (v.reporder) {
                const output_repeatable = f.repeatable(v.reporder);
                f.add(output_repeatable, output_indexed);
            }
            let output_cotisationsdettes;
            if (v.cotisation && v.debit) {
                output_cotisationsdettes = f.cotisationsdettes(v.cotisation, v.debit, periodes, date_fin);
                f.add(output_cotisationsdettes, output_indexed);
            }
            if (v.delai) {
                const output_delai = f.delais(v.delai, output_cotisationsdettes || {}, {
                    premièreDate: serie_periode[0],
                    dernièreDate: serie_periode[serie_periode.length - 1],
                });
                f.add(output_delai, output_indexed);
            }
            v.altares = v.altares || {};
            v.procol = v.procol || {};
            if (v.altares) {
                f.defaillances(v.altares, v.procol, output_indexed);
            }
            if (v.ccsf) {
                f.ccsf(v.ccsf, output_array);
            }
            if (v.sirene) {
                f.sirene(v.sirene, output_array);
            }
            f.populateNafAndApe(output_indexed, naf);
            const output_cotisation = f.cotisation(output_indexed);
            f.add(output_cotisation, output_indexed);
            const output_cible = f.cibleApprentissage(output_indexed, 18);
            f.add(output_cible, output_indexed);
            output_array.forEach((val) => {
                const data = {
                    [this._id]: val,
                };
                emit({
                    batch: actual_batch,
                    siren: this._id.substring(0, 9),
                    periode: val.periode,
                    type: "other",
                }, data);
            });
        }
    }
    if (v.scope === "entreprise") {
        if (includes["all"]) {
            const output_indexed = {};
            for (const periode of serie_periode) {
                output_indexed[periode.getTime()] = {
                    siren: v.key,
                    periode,
                    exercice_bdf: 0,
                    arrete_bilan_bdf: new Date(0),
                    exercice_diane: 0,
                    arrete_bilan_diane: new Date(0),
                };
            }
            if (v.sirene_ul) {
                const outputEntrSirene = f.entr_sirene(v.sirene_ul, serie_periode);
                f.add(outputEntrSirene, output_indexed);
            }
            const periodes = serie_periode.map((date) => date.getTime());
            if (v.effectif_ent) {
                const output_effectif_ent = f.effectifs(v.effectif_ent, periodes, "effectif_ent");
                f.add(output_effectif_ent, output_indexed);
            }
            v.bdf = v.bdf || {};
            v.diane = v.diane || {};
            if (v.bdf) {
                const outputBdf = f.entr_bdf(v.bdf, periodes);
                f.add(outputBdf, output_indexed);
            }
            if (v.diane) {
                const outputDiane = f.entr_diane(v.diane, output_indexed, periodes);
                f.add(outputDiane, output_indexed);
            }
            serie_periode.forEach((date) => {
                const periode = output_indexed[date.getTime()];
                if ((periode.arrete_bilan_bdf || new Date(0)).getTime() === 0 &&
                    (periode.arrete_bilan_diane || new Date(0)).getTime() === 0) {
                    delete output_indexed[date.getTime()];
                }
                if ((periode.arrete_bilan_bdf || new Date(0)).getTime() === 0) {
                    delete periode.arrete_bilan_bdf;
                }
                if ((periode.arrete_bilan_diane || new Date(0)).getTime() === 0) {
                    delete periode.arrete_bilan_diane;
                }
                emit({
                    batch: actual_batch,
                    siren: this._id.substring(0, 9),
                    periode: periode.periode,
                    type: "other",
                }, {
                    entreprise: periode,
                });
            });
        }
    }
}`,
"nbDays": `const nbDays = (firstDate, secondDate) => {
    const oneDay = 24 * 60 * 60 * 1000; // hours*minutes*seconds*milliseconds
    return Math.round(Math.abs((firstDate.getTime() - secondDate.getTime()) / oneDay));
};`,
"outputs": `/**
 * Appelé par ` + "`" + `map()` + "`" + ` pour chaque entreprise/établissement, ` + "`" + `outputs()` + "`" + ` retourne
 * un tableau contenant un objet de base par période, ainsi qu'une version
 * indexée par période de ce tableau, afin de faciliter l'agrégation progressive
 * de données dans ces structures par ` + "`" + `map()` + "`" + `.
 */
function outputs(v, serie_periode) {
    "use strict";
    const output_array = serie_periode.map(function (e) {
        return {
            siret: v.key,
            periode: e,
            effectif: null,
            etat_proc_collective: "in_bonis",
            interessante_urssaf: true,
            outcome: false,
        };
    });
    const output_indexed = output_array.reduce(function (periodes, val) {
        periodes[val.periode.getTime()] = val;
        return periodes;
    }, {});
    return [output_array, output_indexed];
}`,
"poidsFrng": `function poidsFrng(diane) {
    "use strict";
    return typeof diane["couverture_ca_fdr"] === "number"
        ? (diane["couverture_ca_fdr"] / 360) * 100
        : null;
}`,
"populateNafAndApe": `function populateNafAndApe(output_indexed, naf) {
    "use strict";
    Object.keys(output_indexed).forEach((k) => {
        const code_ape = output_indexed[k].code_ape;
        if (code_ape) {
            const code_naf = naf.n5to1[code_ape];
            output_indexed[k].code_naf = code_naf;
            output_indexed[k].libelle_naf = naf.n1[code_naf];
            const code_ape_niveau2 = code_ape.substring(0, 2);
            output_indexed[k].code_ape_niveau2 = code_ape_niveau2;
            const code_ape_niveau3 = code_ape.substring(0, 3);
            output_indexed[k].code_ape_niveau3 = code_ape_niveau3;
            const code_ape_niveau4 = code_ape.substring(0, 4);
            output_indexed[k].code_ape_niveau4 = code_ape_niveau4;
            output_indexed[k].libelle_ape2 = naf.n2[code_ape_niveau2];
            output_indexed[k].libelle_ape3 = naf.n3[code_ape_niveau3];
            output_indexed[k].libelle_ape4 = naf.n4[code_ape_niveau4];
            output_indexed[k].libelle_ape5 = naf.n5[code_ape];
        }
    });
}`,
"reduce": `function reduce(_key, values) {
    "use strict";
    return values.reduce((val, accu) => {
        return Object.assign(accu, val);
    }, {});
}`,
"repeatable": `function repeatable(rep) {
    "use strict";
    const output_repeatable = {};
    Object.keys(rep).forEach((hash) => {
        const one_rep = rep[hash];
        const periode = one_rep.periode.getTime();
        output_repeatable[periode] = output_repeatable[periode] || {};
        output_repeatable[periode].random_order = one_rep.random_order;
    });
    return output_repeatable;
}`,
"sirene": `function sirene(vSirene, output_array) {
    "use strict";
    const sireneHashes = Object.keys(vSirene || {});
    output_array.forEach((val) => {
        // geolocalisation
        if (sireneHashes.length !== 0) {
            const sirene = vSirene[sireneHashes[sireneHashes.length - 1]];
            val.siren = val.siret.substring(0, 9);
            val.latitude = sirene.lattitude || null;
            val.longitude = sirene.longitude || null;
            val.departement = sirene.departement || null;
            if (val.departement) {
                val.region = f.region(val.departement);
            }
            const regexp_naf = /^[0-9]{4}[A-Z]$/;
            if (sirene.ape && sirene.ape.match(regexp_naf)) {
                val.code_ape = sirene.ape;
            }
            val.raison_sociale = sirene.raison_sociale || null;
            // val.activite_saisonniere = sirene.activite_saisoniere || null
            // val.productif = sirene.productif || null
            // val.tranche_ca = sirene.tranche_ca || null
            // val.indice_monoactivite = sirene.indice_monoactivite || null
            val.date_creation_etablissement = sirene.date_creation
                ? sirene.date_creation.getFullYear()
                : null;
            if (val.date_creation_etablissement) {
                val.age =
                    sirene.date_creation && sirene.date_creation >= new Date("1901/01/01")
                        ? val.periode.getFullYear() - val.date_creation_etablissement
                        : null;
            }
        }
    });
}`,
"tauxMarge": `function tauxMarge(diane) {
    "use strict";
    var _a, _b;
    const ratio = ((_a = diane["excedent_brut_d_exploitation"]) !== null && _a !== void 0 ? _a : NaN) /
        ((_b = diane["valeur_ajoutee"]) !== null && _b !== void 0 ? _b : NaN);
    return isNaN(ratio) ? null : ratio * 100;
}`,
},
}
