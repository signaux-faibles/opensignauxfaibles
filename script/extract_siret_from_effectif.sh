cat "$@" |
   # keep only siret and effectif
   csvcut -d";" -c 2,6-117 |
   # keep only companies with any effectif >= 10
   sed -n '/[0-9]\{14\},.*[0-9]\{2,\}/p' |
   # keep only sirens
   awk '{print substr($0,1,9)}' |
   sort |
   uniq
