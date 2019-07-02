#!/usr/bin/env bash
while getopts es option; do
  case "$option" in
    e) EFFECTIF=true;;
    s) SIREN=true;;
  esac
done
shift $(($OPTIND -1))

if [ $# -eq 0 ]; then
  echo "No arguments supplied"
  exit 1
fi

cat "$@" |
    csvgrep --regex "A" --columns 41 |
    if [ -n "$EFFECTIF" ]; then csvgrep --invert-match --regex "(NN|00|01|02|03)" --columns 6; else cat; fi |
      if [ -n "$SIREN" ]; then csvcut --quoting 3 --columns 1; else cat; fi


