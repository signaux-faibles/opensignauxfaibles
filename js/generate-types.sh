#!/bin/bash

OUT_FILE="GeneratedTypes.d.ts"

NODE_BIN=$(npm bin)

echo > "${OUT_FILE}" "\
/**
 * This file was automatically generated by generate-types.sh.
 * 
 * DO NOT MODIFY IT BY HAND.
 *
 * Instead:
 * - modify the validation/*.schema.json files;
 * - then, run generate-types.sh to regenerate this file.
 */

$("${NODE_BIN}/mongodb-json-schema-to-typescript" --input "../validation/*.schema.json" --bannerComment '')"

mkdir -p reduce.algo2/docs/
npx typescript-json-schema --ignoreErrors --refs false reduce.algo2/entr_diane.ts SortieDiane > reduce.algo2/docs/entr_diane.out.json 

"${NODE_BIN}/eslint" "${OUT_FILE}" --fix
