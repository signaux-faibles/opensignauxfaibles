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

# Generate reduce.algo2/docs/entr_diane.out.json
(cd reduce.algo2/docs/ && ${NODE_BIN}/ts-node generate-docs.ts "../entr_diane.ts")

"${NODE_BIN}/eslint" "${OUT_FILE}" --fix
