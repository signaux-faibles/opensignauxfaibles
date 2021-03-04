package engine 

 var jsFunctions = map[string]map[string]string{
"common":{
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
        const dataInBatch = v.batch[batch];
        if (dataInBatch === undefined)
            return m;
        // Types intéressants = nouveaux types, ou types avec suppressions
        const delete_types = Object.keys((dataInBatch.compact || {}).delete || {});
        const new_types = Object.keys(dataInBatch);
        const all_interesting_types = [
            ...new Set([...delete_types, ...new_types]),
        ];
        all_interesting_types.forEach((type) => {
            var _a, _b;
            const typedData = m[type];
            if (typeof typedData === "object") {
                // On supprime les clés qu'il faut
                const keysToDelete = ((_b = (_a = dataInBatch.compact) === null || _a === void 0 ? void 0 : _a.delete) === null || _b === void 0 ? void 0 : _b[type]) || [];
                for (const hash of keysToDelete) {
                    delete typedData[hash];
                }
            }
            else {
                m[type] = {};
            }
            Object.assign(m[type], dataInBatch[type]);
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
"makePeriodeMap": `/**
 * makePeriodeMap() retourne une nouvelle instance de la classe ParPériode
 * (équivalente à Map<Timestamp, T>). Cette fonction a été fournie à défaut
 * d'être parvenu à inclure directement la classe ParPériode dans le scope
 * transmis à MongoDB depuis le traitement map-reduce lancé par le code Go.
 * @param arg (optionnel) - pour initialiser la Map avec un tableau d'entries.
 */
function makePeriodeMap(arg) {
    /**
     * IntMap est une ré-implémentation partielle de Map<Timestamp, Value>
     * utilisant un objet JavaScript pour indexer les entrées, et rendue
     * nécéssaire par le fait que la classe Map de MongoDB n'est pas standard.
     */
    class IntMap {
        constructor(entries) {
            this.data = {};
            if (entries) {
                for (const [key, value] of entries) {
                    this.data[key] = value;
                }
            }
        }
        has(key) {
            return Object.prototype.hasOwnProperty.call(this.data, key);
        }
        get(key) {
            return this.data[key];
        }
        set(key, value) {
            this.data[key] = value;
            return this;
        }
        get size() {
            return Object.keys(this.data).length;
        }
        clear() {
            this.data = {};
        }
        delete(key) {
            const exists = this.has(key);
            delete this.data[key];
            return exists;
        }
        *keys() {
            for (const k of Object.keys(this.data)) {
                yield parseInt(k);
            }
        }
        *values() {
            for (const val of Object.values(this.data)) {
                yield val;
            }
        }
        *entries() {
            for (const [k, v] of Object.entries(this.data)) {
                yield [parseInt(k), v];
            }
        }
        forEach(callbackfn, thisArg) {
            for (const [key, value] of this.entries()) {
                callbackfn.call(thisArg, value, key, this);
            }
        }
        [Symbol.iterator]() {
            return this.entries();
        }
        get [Symbol.toStringTag]() {
            return "IntMap";
        }
    }
    /**
     * Cette classe étend Map<Timestamp, T> pour valider les dates passées
     * en tant que clés et supporter diverses représentations de ces dates
     * (ex: instance Date, timestamp numérique ou chaine de caractères), tout en
     * evitant que des chaines de caractères arbitaires y soient passées.
     */
    class ParPériodeImpl extends IntMap {
        /** Extraie le timestamp d'une date, quelque soit sa représentation. */
        getNumericValue(période) {
            if (typeof période === "number")
                return période;
            if (typeof période === "string")
                return parseInt(période);
            if (période instanceof Date)
                return période.getTime();
            throw new TypeError("type non supporté: " + typeof période);
        }
        /** Vérifie que le timestamp retourné par getNumericValue est valide. */
        getTimestamp(période) {
            const timestamp = this.getNumericValue(période);
            if (isNaN(timestamp) || new Date(timestamp).getTime() !== timestamp) {
                throw new RangeError("valeur invalide: " + période);
            }
            return timestamp;
        }
        /** @throws TypeError ou RangeError si la période n'est pas valide. */
        has(période) {
            return super.has(this.getTimestamp(période));
        }
        /** @throws TypeError ou RangeError si la période n'est pas valide. */
        get(période) {
            return super.get(this.getTimestamp(période));
        }
        /** @throws TypeError ou RangeError si la période n'est pas valide. */
        set(période, val) {
            const timestamp = this.getTimestamp(période);
            super.set(timestamp, val);
            return this;
        }
        /** @throws TypeError ou RangeError si la période n'est pas valide. */
        assign(période, val) {
            var _a;
            const timestamp = this.getTimestamp(période);
            const current = (_a = super.get(timestamp)) !== null && _a !== void 0 ? _a : {};
            super.set(timestamp, Object.assign(current, val));
            return this;
        }
    }
    return new ParPériodeImpl(arg);
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
    f.forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
        currentBatch[type] = [...hashesToAdd].reduce((typedBatchValues, hash) => {
            var _a;
            return (Object.assign(Object.assign({}, typedBatchValues), { [hash]: (_a = currentBatch[type]) === null || _a === void 0 ? void 0 : _a[hash] }));
        }, {});
    });
    // Retrait des propriété vides
    // - compact.delete vides
    const compactDelete = (_a = currentBatch.compact) === null || _a === void 0 ? void 0 : _a.delete;
    if (compactDelete) {
        f.forEachPopulatedProp(compactDelete, (type, keysToDelete) => {
            if (keysToDelete.length === 0) {
                delete compactDelete[type];
            }
        });
        if (Object.keys(compactDelete).length === 0) {
            delete currentBatch.compact;
        }
    }
    // - types vides
    f.forEachPopulatedProp(currentBatch, (type, typedBatchData) => {
        if (Object.keys(typedBatchData).length === 0) {
            delete currentBatch[type];
        }
    });
}`,
"applyPatchesToMemory": `function applyPatchesToMemory(hashToAdd, hashToDelete, memory) {
    // Prise en compte des suppressions de clés dans la mémoire
    f.forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
        hashesToDelete.forEach((hash) => {
            var _a;
            (_a = memory[type]) === null || _a === void 0 ? void 0 : _a.delete(hash);
        });
    });
    // Prise en compte des ajouts de clés dans la mémoire
    f.forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
        hashesToAdd.forEach((hash) => {
            var _a;
            memory[type] = memory[type] || new Set();
            (_a = memory[type]) === null || _a === void 0 ? void 0 : _a.add(hash);
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
    var _a;
    // Les types où il y a potentiellement des suppressions
    const stockTypes = ((_a = completeTypes[fromBatchKey]) !== null && _a !== void 0 ? _a : []).filter((type) => (memory[type] || new Set()).size > 0);
    const { hashToAdd, hashToDelete } = f.listHashesToAddAndDelete(currentBatch, stockTypes, memory);
    f.fixRedundantPatches(hashToAdd, hashToDelete, memory);
    f.applyPatchesToMemory(hashToAdd, hashToDelete, memory);
    f.applyPatchesToBatch(hashToAdd, hashToDelete, stockTypes, currentBatch);
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
        var _a;
        const reporder = (_a = object.batch[batch]) === null || _a === void 0 ? void 0 : _a.reporder;
        if (reporder === undefined)
            return;
        Object.keys(reporder).forEach((ro) => {
            var _a;
            const periode = (_a = reporder[ro]) === null || _a === void 0 ? void 0 : _a.periode;
            if (periode === undefined)
                return;
            if (!missing[periode.getTime()]) {
                delete reporder[ro];
            }
            else {
                missing[periode.getTime()] = false;
            }
        });
    });
    const lastBatch = batches[batches.length - 1];
    if (lastBatch === undefined)
        throw new Error("the last batch should not be undefined");
    serie_periode
        .filter((p) => missing[p.getTime()])
        .forEach((p) => {
        var _a;
        const dataInLastBatch = object.batch[lastBatch];
        if (dataInLastBatch === undefined)
            return;
        const reporder_obj = (_a = dataInLastBatch.reporder) !== null && _a !== void 0 ? _a : {};
        reporder_obj[p.toString()] = {
            random_order: Math.random(),
            periode: p,
            siret,
        };
        dataInLastBatch.reporder = reporder_obj;
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
        var _a;
        //1. On supprime les clés de la mémoire
        if (batch.compact) {
            f.forEachPopulatedProp(batch.compact.delete, (type, keysToDelete) => {
                keysToDelete.forEach((key) => {
                    var _a;
                    (_a = m[type]) === null || _a === void 0 ? void 0 : _a.delete(key); // Should never fail or collection is corrupted
                });
            });
        }
        //2. On ajoute les nouvelles clés
        for (const type of typedObjectKeys(batch)) {
            if (type === "compact")
                continue;
            m[type] = m[type] || new Set();
            for (const key in batch[type]) {
                (_a = m[type]) === null || _a === void 0 ? void 0 : _a.add(key);
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
    let o = Object.assign(Object.assign({}, companyDataValues), { index: { algo2: false } });
    if (o.scope === "entreprise") {
        o.index.algo2 = true;
    }
    else {
        // Est-ce que l'un des batchs a un effectif ?
        const batches = Object.keys(o.batch);
        batches.some((batch) => {
            var _a;
            const hasEffectif = Object.keys(((_a = o.batch[batch]) === null || _a === void 0 ? void 0 : _a.effectif) || {}).length > 0;
            o.index.algo2 = hasEffectif;
            return hasEffectif;
        });
        // Complete reporder if missing
        o = f.complete_reporder(k, o);
    }
    return o;
}`,
"fixRedundantPatches": `/**
 * Modification de hashToAdd et hashToDelete pour retirer les redondances.
 **/
function fixRedundantPatches(hashToAdd, hashToDelete, memory) {
    f.forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
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
    f.forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
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
    f.forEachPopulatedProp(currentBatch, (type) => {
        var _a;
        // Le type compact gère les clés supprimées
        // Ce type compact existe si le batch en cours a déjà été compacté.
        if (type === "compact") {
            const compactDelete = (_a = currentBatch.compact) === null || _a === void 0 ? void 0 : _a.delete;
            if (compactDelete) {
                f.forEachPopulatedProp(compactDelete, (deleteType, keysToDelete) => {
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
            ...(memory[type] || new Set()),
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
    if (values.length === 0)
        throw new Error(` + "`" + `reduce: values of key ${key} should contain at least one item` + "`" + `);
    const firstValue = values[0];
    if (firstValue === undefined)
        throw new Error(` + "`" + `reduce: values of key ${key} should contain at least one item` + "`" + `);
    // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
    // batchs. Sinon, juste fusion des attributs
    const auxBatchSet = new Set();
    const severalBatches = values.some((value) => {
        Object.keys(value.batch || {}).forEach((batch) => auxBatchSet.add(batch));
        return auxBatchSet.size > 1;
    });
    // Fusion batch par batch des types de données sans se préoccuper des doublons.
    const naivelyMergedCompanyData = values.reduce((m, value) => {
        Object.keys(value.batch).forEach((batch) => {
            var _a;
            const dataInBatch = (_a = value.batch[batch]) !== null && _a !== void 0 ? _a : {};
            m.batch[batch] = Object.keys(dataInBatch).reduce((batchValues, type) => (Object.assign(Object.assign({}, batchValues), { [type]: Object.assign(Object.assign({}, batchValues[type]), dataInBatch[type]) })), m.batch[batch] || {});
        });
        return m;
    }, { key, scope: firstValue.scope, batch: {} });
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
        const dataInBatch = naivelyMergedCompanyData.batch[batch];
        if (dataInBatch !== undefined)
            m.push(dataInBatch);
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
        const mergedBatch = naivelyMergedCompanyData.batch[batch];
        if (mergedBatch !== undefined) {
            reducedValue.batch[batch] = mergedBatch;
        }
    });
    // On itère sur chaque batch à partir de fromBatchKey pour les compacter.
    // Il est possible qu'il y ait moins de batch en sortie que le nombre traité
    // dans la boucle, si ces batchs n'apportent aucune information nouvelle.
    batches
        .filter((batch) => batch >= fromBatchKey)
        .forEach((batch) => {
        const currentBatch = naivelyMergedCompanyData.batch[batch];
        if (currentBatch !== undefined) {
            const compactedBatch = f.compactBatch(currentBatch, memory, batch);
            if (Object.keys(compactedBatch).length > 0) {
                reducedValue.batch[batch] = compactedBatch;
            }
        }
    });
    return reducedValue;
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
    return Object.values(apconso !== null && apconso !== void 0 ? apconso : {}).sort((p1, p2) => p1.periode < p2.periode ? 1 : -1);
}`,
"apdemande": `function apdemande(apdemande) {
    return Object.values(apdemande !== null && apdemande !== void 0 ? apdemande : {}).sort((p1, p2) => p1.periode.start.getTime() < p2.periode.start.getTime() ? 1 : -1);
}`,
"bdf": `function bdf(hs) {
    "use strict";
    const bdf = {};
    // Déduplication par arrete_bilan_bdf
    Object.values(hs !== null && hs !== void 0 ? hs : {})
        .filter((b) => b.arrete_bilan_bdf)
        .forEach((b) => {
        bdf[b.arrete_bilan_bdf.toISOString()] = b;
    });
    return Object.values(bdf !== null && bdf !== void 0 ? bdf : {}).sort((a, b) => a.annee_bdf < b.annee_bdf ? 1 : -1);
}`,
"compte": `function compte(compte) {
    const c = Object.values(compte !== null && compte !== void 0 ? compte : {});
    return c.length > 0 ? c[c.length - 1] : undefined;
}`,
"cotisations": `function cotisations(vcotisation = {}) {
    const offset_cotisation = 0;
    const value_cotisation = {};
    // Répartition des cotisations sur toute la période qu'elle concerne
    for (const cotisation of Object.values(vcotisation)) {
        const periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end);
        periode_cotisation.forEach((date_cotisation) => {
            const date_offset = f.dateAddMonth(date_cotisation, offset_cotisation);
            value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat([cotisation.du / periode_cotisation.length]);
        });
    }
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
"debits": `function debits(vdebit = {}) {
    var _a;
    const last_treatment_day = 20;
    const ecn = {};
    for (const [h, debit] of Object.entries(vdebit)) {
        const start = debit.periode.start;
        const end = debit.periode.end;
        const num_ecn = debit.numero_ecart_negatif;
        const compte = debit.numero_compte;
        const key = start + "-" + end + "-" + num_ecn + "-" + compte;
        ecn[key] = (ecn[key] || []).concat([
            {
                hash: h,
                numero_historique: debit.numero_historique,
                date_traitement: debit.date_traitement,
            },
        ]);
    }
    for (const ecnItem of Object.values(ecn)) {
        ecnItem.sort(f.compareDebit);
        const l = ecnItem.length;
        ecnItem.forEach((e, idx) => {
            if (idx <= l - 2) {
                const hashedDataInVDebit = vdebit[e === null || e === void 0 ? void 0 : e.hash];
                const next = ecnItem[idx + 1];
                if (hashedDataInVDebit !== undefined && next !== undefined) {
                    hashedDataInVDebit.debit_suivant = next.hash;
                }
            }
        });
    }
    const value_dette = {};
    for (const debit of Object.values(vdebit)) {
        const nextDate = (debit.debit_suivant && ((_a = vdebit[debit.debit_suivant]) === null || _a === void 0 ? void 0 : _a.date_traitement)) ||
            date_fin;
        //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
        const jour_traitement = debit.date_traitement.getUTCDate();
        const jour_traitement_suivant = nextDate.getUTCDate();
        let date_traitement_debut;
        if (jour_traitement <= last_treatment_day) {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1));
        }
        let date_traitement_fin;
        if (jour_traitement_suivant <= last_treatment_day) {
            date_traitement_fin = new Date(Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth()));
        }
        else {
            date_traitement_fin = new Date(Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth() + 1));
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
                    montant_majorations: /*debit.montant_majorations ||*/ 0,
                },
            ]);
        });
    }
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
    return Object.values(delai !== null && delai !== void 0 ? delai : {});
}`,
"diane": `function diane(hs) {
    "use strict";
    const diane = {};
    // Déduplication par arrete_bilan_diane
    Object.values(hs !== null && hs !== void 0 ? hs : {})
        .filter((d) => d.arrete_bilan_diane)
        .forEach((d) => {
        diane[d.arrete_bilan_diane.toISOString()] = d;
    });
    return Object.values(diane !== null && diane !== void 0 ? diane : {}).sort((a, b) => { var _a, _b; return ((_a = a.exercice_diane) !== null && _a !== void 0 ? _a : 0) < ((_b = b.exercice_diane) !== null && _b !== void 0 ? _b : 0) ? 1 : -1; });
}`,
"effectifs": `function effectifs(effectif) {
    const mapEffectif = f.makePeriodeMap();
    Object.values(effectif !== null && effectif !== void 0 ? effectif : {}).forEach((e) => {
        mapEffectif.set(e.periode, (mapEffectif.get(e.periode) || 0) + e.effectif);
    });
    return serie_periode
        .map((p) => {
        return {
            periode: p,
            effectif: mapEffectif.get(p) || -1,
        };
    })
        .filter((p) => p.effectif >= 0);
}`,
"finalize": `function finalize(_key, val) {
    return val;
}`,
"joinUrssaf": `function joinUrssaf(effectif, debit) {
    const result = {
        effectif: [],
        part_patronale: [],
        part_ouvriere: [],
        montant_majorations: [],
    };
    for (const [i, d] of debit.entries()) {
        const e = effectif.filter((e) => { var _a; return ((_a = serie_periode[i]) === null || _a === void 0 ? void 0 : _a.getTime()) === e.periode.getTime(); });
        if (e.length > 0 && e[0] !== undefined) {
            result.effectif.push(e[0].effectif);
        }
        else {
            result.effectif.push(null);
        }
        result.part_patronale.push(d.part_patronale);
        result.part_ouvriere.push(d.part_ouvriere);
        result.montant_majorations.push(d.montant_majorations);
    }
    return result;
}`,
"map": `function map() {
    var _a, _b, _c, _d, _e, _f;
    const value = f.flatten(this.value, actual_batch);
    if (this.value.scope === "etablissement") {
        const vcmde = {};
        vcmde.key = this.value.key;
        vcmde.batch = actual_batch;
        vcmde.sirene = f.sirene(Object.values((_a = value.sirene) !== null && _a !== void 0 ? _a : {}));
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
        vcmde.procol = Object.values((_b = value.procol) !== null && _b !== void 0 ? _b : {});
        emit("etablissement_" + this.value.key, vcmde);
    }
    else if (this.value.scope === "entreprise") {
        const v = {};
        const diane = f.diane(value.diane);
        const bdf = f.bdf(value.bdf);
        const sirene_ul = (_d = Object.values((_c = value.sirene_ul) !== null && _c !== void 0 ? _c : {})[0]) !== null && _d !== void 0 ? _d : null;
        const ellisphere = (_f = Object.values((_e = value.ellisphere) !== null && _e !== void 0 ? _e : {})[0]) !== null && _f !== void 0 ? _f : null;
        if (sirene_ul) {
            sirene_ul.raison_sociale = f.raison_sociale(sirene_ul.raison_sociale, sirene_ul.nom_unite_legale, sirene_ul.nom_usage_unite_legale, sirene_ul.prenom1_unite_legale, sirene_ul.prenom2_unite_legale, sirene_ul.prenom3_unite_legale, sirene_ul.prenom4_unite_legale);
        }
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
        if (ellisphere) {
            v.ellisphere = ellisphere;
        }
        if (value.paydex) {
            v.paydex = Object.values(value.paydex).sort((p1, p2) => p1.date_valeur.getTime() - p2.date_valeur.getTime());
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
    return sireneArray[sireneArray.length - 1] || {}; // TODO: vérifier que sireneArray est bien classé dans l'ordre chronologique -> c'est sûr qu'il ne l'est pas, vérifier que pour toute la base on a bien un objet sirene unique !
}`,
},
"purgeBatch":{
"finalize": `function finalize(k, o) {
    "use strict";
    return o
}`,
"map": `function map() {
  "use strict";
  const batches = Object.keys(this.value.batch)
  batches.filter((key) => key >= fromBatchKey).forEach((key) => {
    delete this.value.batch[key]
  })
  // With a merge output, sending a new object, even empty, is compulsory
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
    for (const période of output.keys()) {
        output.assign(période, obj.get(période));
    }
}`,
"apart.crossComputation.json": `{
  "$set": {
    "value.ratio_apart": {
      "$let": {
        "vars": {
          "nbHeuresOuvreesMoyParMois": 157.67
        },
        "in": {
          "$divide": [
            "$value.apart_heures_consommees",
            {
              "$multiply": ["$value.effectif", "$$nbHeuresOuvreesMoyParMois"]
            }
          ]
        }
      }
    }
  }
}`,
"apart": `function apart(apconso, apdemande) {
    "use strict";
    const output_apart = f.makePeriodeMap();
    // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
    const apart = {};
    for (const [hash, apdemandeEntry] of Object.entries(apdemande)) {
        apart[apdemandeEntry.id_demande.substring(0, 9)] = {
            demande: hash,
            consommation: [],
            periode_debut: new Date(0),
            periode_fin: new Date(0),
        };
    }
    // on note le nombre d'heures demandées dans output_apart
    for (const apdemandeEntry of Object.values(apdemande)) {
        const periode_deb = apdemandeEntry.periode.start;
        const periode_fin = apdemandeEntry.periode.end;
        // Des periodes arrondies aux débuts de périodes
        // TODO: arrondir au debut du mois le plus proche, au lieu de tronquer la date. (ex: cas du dernier jour d'un mois)
        const periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0));
        const periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        const apartForSiren = apart[apdemandeEntry.id_demande.substring(0, 9)];
        if (apartForSiren === undefined) {
            const error = (message) => {
                throw new Error(message);
            };
            error("siren should be included in apart");
        }
        else {
            apartForSiren.periode_debut = periode_deb_floor;
            apartForSiren.periode_fin = periode_fin_ceil;
        }
        f.generatePeriodSerie(periode_deb_floor, periode_fin_ceil).forEach((période) => {
            output_apart.assign(période, {
                apart_heures_autorisees: apdemandeEntry.hta,
            });
        });
    }
    // relier les consos faites aux demandes (hashs) dans apart
    for (const [hash, valueap] of Object.entries(apconso)) {
        const apartForSiren = apart[valueap.id_conso.substring(0, 9)];
        if (apartForSiren !== undefined) {
            apartForSiren.consommation.push(hash);
        }
    }
    for (const apartEntry of Object.values(apart)) {
        if (apartEntry.consommation.length > 0) {
            apartEntry.consommation
                .sort((a, b) => {
                var _a, _b, _c, _d;
                return ((_b = (_a = apconso[a]) === null || _a === void 0 ? void 0 : _a.periode) !== null && _b !== void 0 ? _b : new Date()).getTime() -
                    ((_d = (_c = apconso[b]) === null || _c === void 0 ? void 0 : _c.periode) !== null && _d !== void 0 ? _d : new Date()).getTime();
            } // TODO: use ` + "`" + `never` + "`" + ` type assertion here?
            )
                .forEach((h) => {
                var _a, _b, _c, _d, _e;
                const période = (_a = apconso[h]) === null || _a === void 0 ? void 0 : _a.periode;
                if (période === undefined) {
                    return;
                }
                const current = (_b = output_apart.get(période)) !== null && _b !== void 0 ? _b : {};
                const heureConso = (_c = apconso[h]) === null || _c === void 0 ? void 0 : _c.heure_consomme;
                if (heureConso !== undefined) {
                    current.apart_heures_consommees =
                        ((_d = current.apart_heures_consommees) !== null && _d !== void 0 ? _d : 0) + heureConso;
                }
                const motifRecours = (_e = apdemande[apartEntry.demande]) === null || _e === void 0 ? void 0 : _e.motif_recours_se;
                if (motifRecours !== undefined) {
                    current.apart_motif_recours = motifRecours;
                }
                output_apart.set(période, current);
            });
            // Heures consommees cumulees sur la demande
            f.generatePeriodSerie(apartEntry.periode_debut, apartEntry.periode_fin).reduce((accu, période) => {
                var _a;
                //output_apart est déjà défini pour les heures autorisées
                const { apart_heures_consommees } = (_a = output_apart.get(période)) !== null && _a !== void 0 ? _a : {};
                accu = accu + (apart_heures_consommees !== null && apart_heures_consommees !== void 0 ? apart_heures_consommees : 0);
                output_apart.assign(période, { apart_heures_consommees_cumulees: accu });
                return accu;
            }, 0);
        }
    }
    // Note: à la fin de l'opération map-reduce, sfdata va calculer la propriété
    // ratio_apart depuis apart.crossComputation.json.
    return output_apart;
}`,
"ccsf": `function ccsf(vCcsf, output_array) {
    "use strict";
    output_array.forEach((val) => {
        let optccsfDateTraitement = new Date(0);
        for (const ccsf of Object.values(vCcsf)) {
            if (ccsf.date_traitement.getTime() < val.periode.getTime() &&
                ccsf.date_traitement.getTime() > optccsfDateTraitement.getTime()) {
                optccsfDateTraitement = ccsf.date_traitement;
            }
        }
        if (optccsfDateTraitement.getTime() !== 0) {
            val.date_ccsf = optccsfDateTraitement;
        }
    });
}`,
"cibleApprentissage": `function cibleApprentissage(output_indexed, n_months) {
    "use strict";
    var _a, _b;
    // Mock two input instead of one for future modification
    const output_cotisation = output_indexed;
    const output_procol = output_indexed;
    // replace with const
    const périodes = [...output_indexed.keys()];
    const merged_info = f.makePeriodeMap();
    for (const période of périodes) {
        merged_info.set(période, {
            outcome: Boolean(((_a = output_procol.get(période)) === null || _a === void 0 ? void 0 : _a.tag_failure) || ((_b = output_cotisation.get(période)) === null || _b === void 0 ? void 0 : _b.tag_default)),
        });
    }
    const output_outcome = f.lookAhead(merged_info, "outcome", n_months, true);
    const output_default = f.lookAhead(output_cotisation, "tag_default", n_months, true);
    const output_failure = f.lookAhead(output_procol, "tag_failure", n_months, true);
    const output_cible = périodes.reduce(function (m, k) {
        const oDefault = output_default.get(k);
        const oFailure = output_failure.get(k);
        return m.set(k, Object.assign(Object.assign(Object.assign({}, output_outcome.get(k)), (oDefault && { time_til_default: oDefault.time_til_outcome })), (oFailure && { time_til_failure: oFailure.time_til_outcome })));
    }, f.makePeriodeMap());
    return output_cible;
}`,
"compte": `function compte(compte) {
    "use strict";
    const output_compte = f.makePeriodeMap();
    //  var offset_compte = 3
    for (const { periode, numero_compte } of Object.values(compte)) {
        output_compte.assign(periode, { compte_urssaf: numero_compte });
    }
    return output_compte;
}`,
"cotisation": `function cotisation(output_indexed) {
    "use strict";
    var _a, _b, _c;
    const sortieCotisation = f.makePeriodeMap();
    const moyenne = (valeurs = []) => valeurs.some((val) => typeof val === "undefined")
        ? undefined
        : valeurs.reduce((p, c) => p + c, 0) / (valeurs.length || 1);
    // calcul de cotisation_moyenne sur 12 mois
    const futureArrays = f.makePeriodeMap();
    for (const [période, input] of output_indexed.entries()) {
        const périodeCourante = (_a = output_indexed.get(période)) === null || _a === void 0 ? void 0 : _a.periode;
        if (périodeCourante === undefined)
            continue;
        const douzeMoisÀVenir = f
            .generatePeriodSerie(périodeCourante, f.dateAddMonth(périodeCourante, 12))
            .filter((périodeFuture) => output_indexed.has(périodeFuture));
        // Accumulation de cotisations sur les 12 mois à venir, pour calcul des moyennes
        douzeMoisÀVenir.forEach((périodeFuture) => {
            const future = futureArrays.get(périodeFuture) || {
                cotisations: [],
                montantsPP: [],
                montantsPO: [],
            };
            future.cotisations.push(input.cotisation);
            future.montantsPP.push(input.montant_part_patronale || 0);
            future.montantsPO.push(input.montant_part_ouvriere || 0);
            futureArrays.set(périodeFuture, future);
        });
        // Calcul des cotisations moyennes à partir des valeurs accumulées ci-dessus
        const { cotisations, montantsPO, montantsPP } = (_b = futureArrays.get(période)) !== null && _b !== void 0 ? _b : {};
        const out = (_c = sortieCotisation.get(période)) !== null && _c !== void 0 ? _c : {};
        if (cotisations && cotisations.length >= 12) {
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
            if (montantsPO &&
                montantsPP &&
                cotisations &&
                !cotisations.includes(undefined) &&
                !cotisations.includes(0)) {
                const detteVals = [];
                for (const [i, cotisation] of cotisations.entries()) {
                    const montPO = montantsPO[i];
                    const montPP = montantsPP[i];
                    if (cotisation !== undefined &&
                        montPO !== undefined &&
                        montPP !== undefined) {
                        detteVals.push((montPO + montPP) / cotisation);
                    }
                }
                out.ratio_dette_moy12m = moyenne(detteVals);
            }
        }
        sortieCotisation.set(période, out);
        // Remplace dans cibleApprentissage
        //val.dette_any_12m = (val.montantsPA || []).reduce((p,c) => (c >=
        //100) || p, false) || (val.montantsPO || []).reduce((p, c) => (c >=
        //100) || p, false)
    }
    // Calcul des défauts URSSAF prolongés
    let counter = 0;
    for (const cotis of sortieCotisation.values()) {
        if (!cotis.ratio_dette)
            continue;
        if (cotis.ratio_dette > 0.01) {
            cotis.tag_debit = true; // Survenance d'un débit d'au moins 1% des cotisations
        }
        if (cotis.ratio_dette > 1) {
            counter = counter + 1;
            if (counter >= 3)
                cotis.tag_default = true;
        }
        else
            counter = 0;
    }
    return sortieCotisation;
}`,
"cotisationsdettes": `/**
 * Calcule les variables liées aux cotisations sociales et dettes sur ces
 * cotisations.
 */
function cotisationsdettes(vCotisation, vDebit, periodes, finPériode // correspond à la variable globale date_fin
) {
    "use strict";
    var _a;
    // Tous les débits traitées après ce jour du mois sont reportées à la période suivante
    // Permet de s'aligner avec le calendrier de fourniture des données
    const lastAccountedDay = 20;
    const sortieCotisationsDettes = f.makePeriodeMap();
    const value_cotisation = f.makePeriodeMap();
    // Répartition des cotisations sur toute la période qu'elle concerne
    for (const cotisation of Object.values(vCotisation)) {
        const periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end);
        periode_cotisation.forEach((date_cotisation) => {
            value_cotisation.set(date_cotisation, (value_cotisation.get(date_cotisation) || []).concat([
                cotisation.du / periode_cotisation.length,
            ]));
        });
    }
    // relier les débits
    // ecn: ecart negatif
    // map les débits: clé fabriquée maison => [{hash, numero_historique, date_traitement}, ...]
    // Pour un même compte, les débits avec le même num_ecn (chaque émission de facture) sont donc regroupés
    const ecn = {};
    for (const [h, debit] of Object.entries(vDebit)) {
        const start = debit.periode.start;
        const end = debit.periode.end;
        const num_ecn = debit.numero_ecart_negatif;
        const compte = debit.numero_compte;
        const key = start + "-" + end + "-" + num_ecn + "-" + compte;
        ecn[key] = (ecn[key] || []).concat([
            {
                hash: h,
                numero_historique: debit.numero_historique,
                date_traitement: debit.date_traitement,
            },
        ]);
    }
    // Pour chaque numero_ecn, on trie et on chaîne les débits avec debit_suivant
    for (const ecnEntry of Object.values(ecn)) {
        ecnEntry.sort(f.compareDebit);
        const l = ecnEntry.length;
        ecnEntry
            .filter((_, idx) => idx <= l - 2)
            .forEach((e, idx) => {
            const vDebitForHash = vDebit[e.hash];
            const next = ((ecnEntry === null || ecnEntry === void 0 ? void 0 : ecnEntry[idx + 1]) || {}).hash;
            if (vDebitForHash && next !== undefined)
                vDebitForHash.debit_suivant = next;
        });
    }
    const value_dette = f.makePeriodeMap();
    // Pour chaque objet debit:
    // debit_traitement_debut => periode de traitement du débit
    // debit_traitement_fin => periode de traitement du debit suivant, ou bien finPériode
    // Entre ces deux dates, c'est cet objet qui est le plus à jour.
    for (const debit of Object.values(vDebit)) {
        const nextDate = (debit.debit_suivant && ((_a = vDebit[debit.debit_suivant]) === null || _a === void 0 ? void 0 : _a.date_traitement)) ||
            finPériode;
        //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
        const jour_traitement = debit.date_traitement.getUTCDate();
        const jour_traitement_suivant = nextDate.getUTCDate();
        let date_traitement_debut;
        if (jour_traitement <= lastAccountedDay) {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth()));
        }
        else {
            date_traitement_debut = new Date(Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1));
        }
        let date_traitement_fin;
        if (jour_traitement_suivant <= lastAccountedDay) {
            date_traitement_fin = new Date(Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth()));
        }
        else {
            date_traitement_fin = new Date(Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth() + 1));
        }
        //f.generatePeriodSerie exlue la dernière période
        f.generatePeriodSerie(date_traitement_debut, date_traitement_fin).forEach((date) => {
            var _a;
            value_dette.set(date, [
                ...((_a = value_dette.get(date)) !== null && _a !== void 0 ? _a : []),
                {
                    periode: debit.periode.start,
                    part_ouvriere: debit.part_ouvriere,
                    part_patronale: debit.part_patronale,
                },
            ]);
        });
    }
    // TODO faire numero de compte ailleurs
    // Array des numeros de compte
    //var numeros_compte = Array.from(new Set(
    //  Object.keys(vCotisation).map(function (h) {
    //    return(vCotisation[h].numero_compte)
    //  })
    //))
    periodes.forEach(function (time) {
        var _a;
        const val = (_a = sortieCotisationsDettes.get(time)) !== null && _a !== void 0 ? _a : {};
        //val.numero_compte_urssaf = numeros_compte
        const valueCotis = value_cotisation.get(time);
        if (valueCotis !== undefined) {
            // somme de toutes les cotisations dues pour une periode donnée
            val.cotisation = valueCotis.reduce((a, cot) => a + cot, 0);
        }
        // somme de tous les débits (part ouvriere, part patronale)
        val.montant_part_ouvriere = (value_dette.get(time) || []).reduce((acc, { part_ouvriere }) => acc + part_ouvriere, 0);
        val.montant_part_patronale = (value_dette.get(time) || []).reduce((acc, { part_patronale }) => acc + part_patronale, 0);
        sortieCotisationsDettes.set(time, val);
        const monthOffsets = [1, 2, 3, 6, 12];
        const futureTimestamps = monthOffsets
            .map((offset) => ({
            offset,
            timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
        }))
            .filter(({ timestamp }) => periodes.includes(timestamp));
        futureTimestamps.forEach(({ offset, timestamp }) => {
            var _a;
            sortieCotisationsDettes.set(timestamp, Object.assign(Object.assign({}, ((_a = sortieCotisationsDettes.get(timestamp)) !== null && _a !== void 0 ? _a : {})), { [` + "`" + `montant_part_ouvriere_past_${offset}` + "`" + `]: val.montant_part_ouvriere, [` + "`" + `montant_part_patronale_past_${offset}` + "`" + `]: val.montant_part_patronale }));
        });
        if (val.montant_part_ouvriere + val.montant_part_patronale > 0) {
            const futureTimestamps = [0, 1, 2, 3, 4, 5]
                .map((offset) => ({
                timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
            }))
                .filter(({ timestamp }) => periodes.includes(timestamp));
            futureTimestamps.forEach(({ timestamp }) => {
                var _a;
                sortieCotisationsDettes.set(timestamp, Object.assign(Object.assign({}, ((_a = sortieCotisationsDettes.get(timestamp)) !== null && _a !== void 0 ? _a : {})), { interessante_urssaf: false }));
            });
        }
    });
    return sortieCotisationsDettes;
}`,
"defaillances": `function defaillances(défaillances, output_indexed) {
    "use strict";
    const codes = Object.keys(défaillances)
        .reduce((events, hash) => {
        const the_event = défaillances[hash];
        let etat = null;
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
        const time_til_last = [...output_indexed.keys()].filter((période) => {
            return période >= periode_effet.getTime();
        });
        time_til_last.forEach((time) => {
            if (output_indexed.has(time)) {
                output_indexed.assign(time, Object.assign({ etat_proc_collective: event.etat, date_proc_collective: event.date_proc_col }, (event.etat !== "in_bonis" && { tag_failure: true })));
            }
        });
    });
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
    const donnéesDélaiParPériode = f.makePeriodeMap();
    Object.values(vDelai).forEach((delai) => {
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
            const remainingDays = f.nbDays(debutDeMois, delai.date_echeance);
            const inputAtTime = debitParPériode.get(debutDeMois);
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
            donnéesDélaiParPériode.set(debutDeMois, outputAtTime);
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
"effectifs": `function effectifs(entréeEffectif, periodes, clé) {
    "use strict";
    var _a;
    const sortieEffectif = f.makePeriodeMap();
    // Construction d'une map[time] = effectif à cette periode
    const mapEffectif = f.makePeriodeMap();
    Object.keys(entréeEffectif).forEach((hash) => {
        const effectif = entréeEffectif[hash];
        if (effectif !== null && effectif !== undefined) {
            mapEffectif.set(effectif.periode, effectif.effectif);
        }
    });
    // On reporte dans les dernières périodes le dernier effectif connu
    // Ne reporter que si le dernier effectif est disponible
    const dernièrePériodeAvecEffectifConnu = f.dateAddMonth(new Date(periodes[periodes.length - 1]), offset_effectif + 1);
    const effectifÀReporter = (_a = mapEffectif.get(dernièrePériodeAvecEffectifConnu)) !== null && _a !== void 0 ? _a : null;
    const makeReporteProp = (clé) => ` + "`" + `${clé}_reporte` + "`" + `;
    periodes.forEach((time) => {
        var _a;
        sortieEffectif.set(time, Object.assign(Object.assign({}, ((_a = sortieEffectif.get(time)) !== null && _a !== void 0 ? _a : {})), { [clé]: mapEffectif.get(time) || effectifÀReporter, [makeReporteProp(clé)]: mapEffectif.get(time) ? 0 : 1 }));
    });
    const makePastProp = (clé, offset) => ` + "`" + `${clé}_past_${offset}` + "`" + `;
    mapEffectif.forEach((effectifAtTime, time) => {
        const futureOffsets = [6, 12, 18, 24];
        const futureTimestamps = futureOffsets
            .map((offset) => ({
            offset,
            timestamp: f
                .dateAddMonth(new Date(time), offset - offset_effectif - 1)
                // TODO: réfléchir à si l'offset est nécessaire pour l'algo.
                // Ces valeurs permettent de calculer les dernières variations réelles
                // d'effectif sur la période donnée (par exemple: 6 mois),
                // en excluant les valeurs reportées qui
                // pourraient conduire à des variations = 0
                .getTime(),
        }))
            .filter(({ timestamp }) => periodes.includes(timestamp));
        futureTimestamps.forEach(({ offset, timestamp }) => {
            var _a;
            sortieEffectif.set(timestamp, Object.assign(Object.assign({}, ((_a = sortieEffectif.get(timestamp)) !== null && _a !== void 0 ? _a : {})), { [makePastProp(clé, offset)]: effectifAtTime }));
        });
    });
    return sortieEffectif;
}`,
"entr_bdf": `function entr_bdf(donnéesBdf, periodes) {
    "use strict";
    const outputBdf = f.makePeriodeMap(periodes.map((période) => [période, {}]));
    for (const entréeBdf of Object.values(donnéesBdf)) {
        const periode_arrete_bilan = new Date(Date.UTC(entréeBdf.arrete_bilan_bdf.getUTCFullYear(), entréeBdf.arrete_bilan_bdf.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7);
        const series = f.generatePeriodSerie(periode_dispo, f.dateAddMonth(periode_dispo, 13));
        for (const periode of series) {
            const outputInPeriod = outputBdf.get(periode) || {};
            const periodData = f.omit(entréeBdf, "raison_sociale", "secteur", "siren");
            // TODO: Éviter d'ajouter des données en dehors de ` + "`" + `periodes` + "`" + `, sans fausser le calcul des données passées (plus bas)
            Object.assign(outputInPeriod, periodData);
            if (outputInPeriod.annee_bdf) {
                outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1;
            }
            const pastData = f.omit(periodData, "arrete_bilan_bdf", "annee_bdf");
            const makePastProp = (prop, offset) => ` + "`" + `${prop}_past_${offset}` + "`" + `;
            for (const prop of Object.keys(pastData)) {
                const past_year_offset = [1, 2];
                for (const offset of past_year_offset) {
                    const periode_offset = f.dateAddMonth(periode, 12 * offset);
                    const outputInPast = outputBdf.get(periode_offset);
                    if (outputInPast) {
                        outputInPast[makePastProp(prop, offset)] = entréeBdf[prop];
                    }
                }
            }
            outputBdf.set(periode, outputInPeriod);
        }
    }
    return outputBdf;
}`,
"entr_diane": `function entr_diane(donnéesDiane, output_indexed, periodes) {
    for (const entréeDiane of Object.values(donnéesDiane)) {
        if (!entréeDiane.arrete_bilan_diane)
            continue;
        //entréeDiane.arrete_bilan_diane = new Date(Date.UTC(entréeDiane.exercice_diane, 11, 31, 0, 0, 0, 0))
        const periode_arrete_bilan = new Date(Date.UTC(entréeDiane.arrete_bilan_diane.getUTCFullYear(), entréeDiane.arrete_bilan_diane.getUTCMonth() + 1, 1, 0, 0, 0, 0));
        const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7); // 01/08 pour un bilan le 31/12, donc algo qui tourne en 01/09
        const series = f.generatePeriodSerie(periode_dispo, f.dateAddMonth(periode_dispo, 14) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (14+7)
        );
        for (const periode of series) {
            const rest = f.omit(entréeDiane, 
            // "marquee",
            "nom_entreprise", "numero_siren", "statut_juridique", "procedure_collective");
            const makePastProp = (prop, offset) => ` + "`" + `${prop}_past_${offset}` + "`" + `;
            if (periodes.includes(periode.getTime())) {
                output_indexed.assign(periode, rest);
            }
            for (const ratio of Object.keys(rest)) {
                if (entréeDiane[ratio] === null) {
                    const outputAtTime = output_indexed.get(periode);
                    if (outputAtTime !== undefined &&
                        periodes.includes(periode.getTime())) {
                        delete outputAtTime[ratio];
                    }
                    continue;
                }
                // Passé
                const past_year_offset = [1, 2];
                for (const offset of past_year_offset) {
                    const periode_offset = f.dateAddMonth(periode, 12 * offset);
                    const variable_name = makePastProp(ratio, offset);
                    const outputAtOffset = output_indexed.get(periode_offset);
                    if (outputAtOffset !== undefined &&
                        ratio !== "arrete_bilan_diane" &&
                        ratio !== "exercice_diane") {
                        outputAtOffset[variable_name] = entréeDiane[ratio];
                    }
                }
            }
        }
        for (const periode of series) {
            const inputInPeriod = output_indexed.get(periode);
            const outputInPeriod = output_indexed.get(periode);
            if (periodes.includes(periode.getTime()) &&
                inputInPeriod &&
                outputInPeriod) {
                // Recalcul BdF si ratios bdf sont absents
                if (!("poids_frng" in inputInPeriod)) {
                    const poids = f.poidsFrng(entréeDiane);
                    if (poids !== null)
                        outputInPeriod.poids_frng = poids;
                }
                if (!("dette_fiscale" in inputInPeriod)) {
                    const dette = f.detteFiscale(entréeDiane);
                    if (dette !== null)
                        outputInPeriod.dette_fiscale = dette;
                }
                if (!("frais_financier" in inputInPeriod)) {
                    const frais = f.fraisFinancier(entréeDiane);
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
                const makePastProp = (clé, offset) => ` + "`" + `${clé}_past_${offset}` + "`" + `;
                bdf_vars.forEach((k) => {
                    if (k in outputInPeriod) {
                        past_year_offset.forEach((offset) => {
                            const periode_offset = f.dateAddMonth(periode, 12 * offset);
                            const variable_name = makePastProp(k, offset);
                            const outputAtOffset = output_indexed.get(periode_offset);
                            if (outputAtOffset &&
                                periodes.includes(periode_offset.getTime())) {
                                outputAtOffset[variable_name] = outputInPeriod[k];
                            }
                        });
                    }
                });
            }
        }
    }
    return output_indexed;
}`,
"entr_paydex": `function entr_paydex(vPaydex, sériePériode) {
    "use strict";
    const paydexParPériode = f.makePeriodeMap();
    // initialisation (avec valeurs N/A par défaut)
    for (const période of sériePériode) {
        paydexParPériode.set(période, {
            paydex_nb_jours: null,
            paydex_nb_jours_past_3: null,
            paydex_nb_jours_past_6: null,
            paydex_nb_jours_past_12: null,
        });
    }
    // population des valeurs
    for (const entréePaydex of Object.values(vPaydex)) {
        const période = Date.UTC(entréePaydex.date_valeur.getUTCFullYear(), entréePaydex.date_valeur.getUTCMonth(), 1);
        const mois3Suivant = f.dateAddMonth(new Date(période), 3).getTime();
        const mois6Suivant = f.dateAddMonth(new Date(période), 6).getTime();
        const annéeSuivante = f.dateAddMonth(new Date(période), 12).getTime();
        const donnéesAdditionnelles = f.makePeriodeMap([
            [période, { paydex_nb_jours: entréePaydex.nb_jours }],
            [mois3Suivant, { paydex_nb_jours_past_3: entréePaydex.nb_jours }],
            [mois6Suivant, { paydex_nb_jours_past_6: entréePaydex.nb_jours }],
            [annéeSuivante, { paydex_nb_jours_past_12: entréePaydex.nb_jours }],
        ]);
        f.add(donnéesAdditionnelles, paydexParPériode);
    }
    return paydexParPériode;
}`,
"entr_sirene": `function entr_sirene(sirene_ul, sériePériode) {
    "use strict";
    const retourEntrSirene = f.makePeriodeMap();
    const sireneHashes = Object.keys(sirene_ul || {});
    sériePériode.forEach((période) => {
        if (sireneHashes.length !== 0) {
            const sirene = sirene_ul[sireneHashes[sireneHashes.length - 1]];
            const val = {};
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
            retourEntrSirene.set(période, val);
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
    const établissements = f.omit(v, "entreprise");
    const entr = Object.assign({}, v.entreprise); // on suppose que v.entreprise est défini
    const output = Object.keys(établissements).map((siret) => {
        var _a, _b, _c, _d;
        const etab = (_a = établissements[siret]) !== null && _a !== void 0 ? _a : {};
        if (etab.effectif) {
            entr.effectif_entreprise = entr.effectif_entreprise || 0 + etab.effectif;
        }
        if (etab.apart_heures_consommees) {
            entr.apart_entreprise =
                (entr.apart_entreprise || 0) + etab.apart_heures_consommees;
        }
        if (etab.montant_part_patronale || etab.montant_part_ouvriere) {
            entr.debit_entreprise =
                ((_b = entr.debit_entreprise) !== null && _b !== void 0 ? _b : 0) +
                    ((_c = etab.montant_part_patronale) !== null && _c !== void 0 ? _c : 0) +
                    ((_d = etab.montant_part_ouvriere) !== null && _d !== void 0 ? _d : 0);
        }
        return Object.assign(Object.assign(Object.assign({}, etab), entr), { nbr_etablissements_connus: Object.keys(établissements).length });
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
"lookAhead": `function lookAhead(data, attr_name, // "outcome" | "tag_default" | "tag_failure",
n_months, past) {
    "use strict";
    // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
    // demander: que va-t-il se passer) ou dans le future (past = false on
    // pourra se demander que s'est-il passé
    const chronologic = (pérA, pérB) => pérA - pérB;
    const reverse = (pérA, pérB) => pérB - pérA;
    let counter = -1;
    const output = [...data.keys()]
        .sort(past ? reverse : chronologic)
        .reduce((m, période) => {
        var _a;
        // Si on a déjà détecté quelque chose, on compte le nombre de périodes
        if (counter >= 0)
            counter = counter + 1;
        if ((_a = data.get(période)) === null || _a === void 0 ? void 0 : _a[attr_name]) {
            // si l'évènement se produit on retombe à 0
            counter = 0;
        }
        if (counter >= 0) {
            // l'évènement s'est produit
            m.set(période, {
                time_til_outcome: counter,
                outcome: counter <= n_months ? true : false,
            });
        }
        return m;
    }, f.makePeriodeMap());
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
        const periodes = serie_periode.map((date) => date.getTime());
        if (includes["apart"] || includes["all"]) {
            if (v.apconso && v.apdemande) {
                const output_apart = f.apart(v.apconso, v.apdemande);
                output_apart.forEach((current, periode) => {
                    if (!output_indexed.has(periode))
                        return; // limiter dans le scope temporel du batch.
                    const data = {
                        [this._id]: Object.assign(Object.assign({}, current), { siret: this._id }),
                    };
                    emit({
                        batch: actual_batch,
                        siren: this._id.substring(0, 9),
                        periode: new Date(periode),
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
            if (v.reporder) {
                const output_repeatable = f.repeatable(v.reporder);
                f.add(output_repeatable, output_indexed);
            }
            let output_cotisationsdettes = f.makePeriodeMap();
            if (v.cotisation && v.debit) {
                output_cotisationsdettes = f.cotisationsdettes(v.cotisation, v.debit, periodes, date_fin);
                f.add(output_cotisationsdettes, output_indexed);
            }
            if (v.delai) {
                const premièreDate = serie_periode[0];
                const dernièreDate = serie_periode[serie_periode.length - 1];
                if (premièreDate === undefined || dernièreDate === undefined) {
                    const error = (message) => {
                        throw new Error(message);
                    };
                    error("serie_periode should not contain undefined values");
                }
                else {
                    const output_delai = f.delais(v.delai, output_cotisationsdettes, {
                        premièreDate,
                        dernièreDate,
                    });
                    f.add(output_delai, output_indexed);
                }
            }
            v.procol = v.procol || {};
            f.defaillances(v.procol, output_indexed);
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
            const output_indexed = f.makePeriodeMap();
            for (const periode of serie_periode) {
                output_indexed.set(periode, {
                    siren: v.key,
                    periode,
                    exercice_bdf: 0,
                });
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
            if (v.paydex) {
                const paydexParPériode = f.entr_paydex(v.paydex, serie_periode);
                f.add(paydexParPériode, output_indexed);
            }
            v.bdf = v.bdf || {};
            v.diane = v.diane || {};
            if (v.bdf) {
                const outputBdf = f.entr_bdf(v.bdf, periodes);
                f.add(outputBdf, output_indexed);
            }
            if (v.diane) {
                /*const outputDiane =*/ f.entr_diane(v.diane, output_indexed, periodes);
                // f.add(outputDiane, output_indexed)
                // TODO: rendre f.entr_diane() pure, c.a.d. faire en sorte qu'elle ne modifie plus output_indexed directement
            }
            serie_periode.forEach((date) => {
                const entrData = output_indexed.get(date);
                if ((entrData === null || entrData === void 0 ? void 0 : entrData.arrete_bilan_bdf) !== undefined ||
                    (entrData === null || entrData === void 0 ? void 0 : entrData.arrete_bilan_diane) !== undefined) {
                    emit({
                        batch: actual_batch,
                        siren: this._id.substring(0, 9),
                        periode: entrData.periode,
                        type: "other",
                    }, {
                        entreprise: entrData,
                    });
                }
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
    const output_indexed = f.makePeriodeMap();
    for (const val of output_array) {
        output_indexed.set(val.periode, val);
    }
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
    for (const outputForKey of output_indexed.values()) {
        const code_ape = outputForKey.code_ape;
        if (code_ape) {
            const code_naf = naf.n5to1[code_ape];
            outputForKey.code_naf = code_naf;
            outputForKey.libelle_naf = code_naf ? naf.n1[code_naf] : undefined;
            const code_ape_niveau2 = code_ape.substring(0, 2);
            outputForKey.code_ape_niveau2 = code_ape_niveau2;
            const code_ape_niveau3 = code_ape.substring(0, 3);
            outputForKey.code_ape_niveau3 = code_ape_niveau3;
            const code_ape_niveau4 = code_ape.substring(0, 4);
            outputForKey.code_ape_niveau4 = code_ape_niveau4;
            outputForKey.libelle_ape2 = naf.n2[code_ape_niveau2];
            outputForKey.libelle_ape3 = naf.n3[code_ape_niveau3];
            outputForKey.libelle_ape4 = naf.n4[code_ape_niveau4];
            outputForKey.libelle_ape5 = naf.n5[code_ape];
        }
    }
}`,
"reduce": `function reduce(_key, values) {
    "use strict";
    return values.reduce((val, accu) => {
        return Object.assign(accu, val);
    }, {});
}`,
"repeatable": `function repeatable(rep) {
    "use strict";
    const output_repeatable = f.makePeriodeMap();
    for (const { periode, random_order } of Object.values(rep)) {
        output_repeatable.assign(periode, { random_order });
    }
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
            val.latitude = sirene.latitude || null;
            val.longitude = sirene.longitude || null;
            val.departement = sirene.departement || null;
            if (val.departement) {
                val.region = f.region(val.departement);
            }
            const regexp_naf = /^[0-9]{4}[A-Z]$/;
            if (sirene.ape && sirene.ape.match(regexp_naf)) {
                val.code_ape = sirene.ape;
            }
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
