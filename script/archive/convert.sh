#!/bin/bash

for i in Diane_Export*
do 	dos2unix -1252 -n $i ./aux/utf8_$i
done

head aux/utf8_Diane_Export_1.txt | grep Siren > aux/concat.txt && grep -vh Siren aux/utf8_*.txt  >> aux/concat.txt

for i in `seq 2012 1 2017`
do echo \"$i\"
awk -v tgt=$i -f aux/transform.awk aux/concat.txt | sed 's/,/./g' > "output_$i.csv" 
done


	



	
	


