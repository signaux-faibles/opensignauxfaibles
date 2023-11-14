echo "stock"
cat resources/SF_DATA_20230706.txt > output.csv
echo "increment 1"
go run convert_increment.go resources/S_202011095834-3_202311010315.csv >> output.csv
echo "increment 2"
go run convert_increment.go resources/S_202011095834-3_202310020319.csv >> output.csv
go run remove_column_experience_paiement.go output.csv > cleaned_output.csv
#echo "ajoute les accents"
#sed -i 's/etat_organisation/état_organisation/' cleaned_output.csv
#sed -i 's/encours_etudies/encours_étudiés/' cleaned_output.csv
echo "résultat"
head cleaned_output.csv
echo "."
echo "."
echo "."
tail cleaned_output.csv
tar -czvf altares.tar.gz cleaned_output.csv
