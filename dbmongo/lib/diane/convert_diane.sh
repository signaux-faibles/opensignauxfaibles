#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

while getopts b: option; do
  case "$option" in
    b) FILES=$(echo ../"${OPTARG}"/diane/*.txt);;
  esac
done
shift $(($OPTIND -1))

[ -z "$1" ] && echo "Please insert files as argument" && exit 1

AWK_SCRIPT='
BEGIN { # Semi-column separated csv as input and output
  FS = ";"
  OFS = ";"
  RE_YEAR = "[[:digit:]][[:digit:]][[:digit:]][[:digit:]]"
  RE_YEAR_SUFFIX = / ([[:digit:]][[:digit:]][[:digit:]][[:digit:]])$/
}
FNR==1 { # Heading row: Change field names
  printf "%s", "\"Annee\""

  for (field = 1; field <= NF; ++field) {

    if ($field !~ RE_YEAR_SUFFIX) { # Field without year
      to_print[++remaining] = field
      printf "%s%s",  OFS, $field
    } else { # Field with year
      match($field, RE_YEAR, year)
      field_name = gensub(" "year[0], "", "g", $field) # Remove year from column name
      field_name = gensub("\r", "", "g", field_name)
      if (!printed[field_name]) {
        ++remaining
        ++printed[field_name]
        printf "%s%s", OFS, field_name;
      }
      to_print[remaining , year[0]] = field
    }
  }
  printf "%s", ORS
}
FNR>1 && $1 !~ "Marquée" { # Data row
  first_year = 2012
  today_year = strftime("%Y")
  for (current_year = first_year; current_year <= today_year; ++current_year) {
    printf "%i", current_year
    for (field = 1; field <= remaining; ++field) {
      if (to_print[field, current_year] && $(to_print[field, current_year])) {
        printf "%s%s", OFS, $(to_print[field, current_year]);
      } else if (to_print[field] && $(to_print[field])) {
        printf "%s%s", OFS, $(to_print[field]);
      } else {
        printf "%s%s", OFS, "\"\"";
      }
    }
    printf "%s", ORS # Each year on a new line
  }
}'

# Concat all exported files /!\ FIX ME: no spaces in file_names !
cat ${FILES:-$@} |
 iconv --from-code UTF-16LE --to-code UTF-8 |
 dos2unix -ascii |
 awk "${AWK_SCRIPT}" |
 sed 's/,/./g'

# Test with: $ go test -v -update && diff -b realData/Diane_Expected_Conversion_1.csv realData/Diane_Expected_Conversion_1.txt
