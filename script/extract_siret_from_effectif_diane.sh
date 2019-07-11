DEPS=\(16\|17\|19\|23\|24\|33\|40\|47\|64\|79\|86\|87\|01\|03\|07\|15\|26\|38\|42\|43\|63\|69\|73\|74\)
cat "$@" |
   # Filter departments
   csvgrep -c 5 -r "$DEPS" -d ";" |
   # keep only siret and effectif
   csvcut -c 2,6-117 |
   # keep only companies with any effectif >= 10
   sed -n '/[0-9]\{14\},.*[0-9]\{2,\}/p' |
   # keep only sirens
   awk '{print substr($0,1,9)}' |
   sort |
   uniq |
   # Filter sirens that we already have
   comm -23 - ../BFC_PDL/siren_all_sorted_no_header.csv |
   # Format
   awk 'BEGIN{print "CF00004,SIREN"}{print "id"NR","$0}'
