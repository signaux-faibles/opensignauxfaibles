#!/bin/sh

set -e # will stop the script if any command fails with a non-zero exit code

# Check that swagger.yaml is valid
npx @apidevtools/swagger-cli validate swagger/swagger.yaml

# Convert swagger.yaml to swagger.json
rm swagger/swagger.json
npx jy-transform swagger/swagger.yaml swagger/swagger.json
