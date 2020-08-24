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
}
FNR==1 { # Change field names
  printf "%s", "\"Annee\""

  for (i=1; i<=NF; ++i) {

    if ($i !~ "201"){ # Field without year
      f[++nf] = i
      printf "%s%s",  OFS, $i
    } else { # Field with year
    match($i, "20..", year)
    field_name = gensub(" "year[0],"","g",$i) # Remove year from column name
    field_name = gensub("\r","","g",field_name)
    if (!visited[field_name]){
      ++nf
      ++visited[field_name]
      printf "%s%s", OFS, field_name;
    }
    f[nf , year[0]] = i
  }
}
printf "%s", ORS
}
FNR>1 && $1 !~ "Marquée" {
  first_year = 2012
  today_year = strftime("%Y")
  for (current_y = first_year; current_y <= today_year; ++current_y){
    printf "%i", current_y
    for (i=1; i<=nf; ++i) {
      if (f[i])
        if (f[i])
          printf "%s%s", OFS, $(f[i]);
        else
          printf "%s%s", OFS, "\"\"";
      else {
        # Only print fields relative to current year
        if (f[i, current_y])
          printf "%s%s", OFS, $(f[i, current_y]);
        else
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
