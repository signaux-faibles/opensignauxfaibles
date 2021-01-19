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

VAR_DOC_FILE="reduce.algo2/docs/variables.json"
echo "Generating ${VAR_DOC_FILE}..."
"${NODE_BIN}/ts-node" "reduce.algo2/docs/generate-docs.ts" \
  "reduce.algo2/entr_diane.ts" \
  "reduce.algo2/ccsf.ts" \
  "reduce.algo2/compte.ts" \
  "reduce.algo2/cotisation.ts" \
  "reduce.algo2/cotisationdettes.ts" \
  > "${VAR_DOC_FILE}"

"${NODE_BIN}/eslint" "${OUT_FILE}" --fix
