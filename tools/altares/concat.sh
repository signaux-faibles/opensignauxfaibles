echo "stock"
go run convert_stock.go resources/SF_DATA_20230706.txt > output.csv
echo "increment 1"
go run convert_increment.go resources/S_202011095834-3_202311010315.csv >> output.csv
echo "increment 2"
go run convert_increment.go resources/S_202011095834-3_202310020319.csv >> output.csv
echo "résultat"
head output.csv
echo "."
echo "."
echo "."
tail output.csv
tar -czf altares.tar.gz output.csv
