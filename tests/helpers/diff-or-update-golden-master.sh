#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

FLAGS="$1" # may include the "--update" flag
GOLDEN_FILE="$2"
OUTPUT_FILE="$3"

COLOR_YELLOW='\033[1;33m'
COLOR_DEFAULT='\033[0m'

# Check if the --update flag was passed
if [[ "${FLAGS}" == *--update* ]]
then
    echo "🖼  Updating the golden master file from ${OUTPUT_FILE} ..."
    cp "${OUTPUT_FILE}" "${GOLDEN_FILE}"
    echo "ℹ️  Updated ${GOLDEN_FILE} => you may have to run: $ git secret hide" # to re-encrypt the golden master file, after having updated it
else
    # Diff between expected and actual output
    echo -e "${COLOR_YELLOW}"
    diff --brief "${GOLDEN_FILE}" "${OUTPUT_FILE}" # if differences are found, the script will exit with a non-zero exit code
    echo -e "${COLOR_DEFAULT}"    
    echo "✅ ${OUTPUT_FILE} matches the golden master file."
fi
