docker run --rm -i golang /bin/bash << COMMANDS
git clone https://github.com/signaux-faibles/opensignauxfaibles.git --depth 1
cd opensignauxfaibles
go get -v -t -d ./...
apt-get update -y  && apt-get install -y nodejs npm dos2unix gawk
(cd dbmongo/lib/engine && go generate -x)
(cd dbmongo/lib/diane && go test -v)
COMMANDS
