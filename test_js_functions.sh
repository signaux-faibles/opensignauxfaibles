echo "Transpiling TypeScript files, and generating the jsFunctions.go bundle..."
cd dbmongo/lib/engine
go generate -x
cd -
echo "Running tests against the JS files (including the ones transpiled from TS)..."
cd dbmongo/js/test/ && ./test_common.sh
echo "✅ Tests passed."
