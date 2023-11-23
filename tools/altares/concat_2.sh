echo "stock"
go run main.go resources/SF_DATA_20230706.txt resources/S_202011095834-3_202311010315.csv resources/S_202011095834-3_202310020319.csv  output_2.csv
#echo "increment 1"
#go run convert_increment.go resources/S_202011095834-3_202311010315.csv >> output.csv
#echo "increment 2"
#go run convert_increment.go resources/S_202011095834-3_202310020319.csv >> output.csv
echo "résultat"
head output_2.csv
echo "."
echo "."
echo "."
tail output_2.csv
#tar -czf altares_2.tar.gz output_2.csv
