#!/bin/sh

if [ "$#" -ne 1 ]; then
    echo "build.sh: construit l'application datapi dans une image"
    echo "usage: build.sh branch"
    echo "exemple: ./build.sh master"
    exit 255
fi

if [ -d workspace ]; then
    echo "supprimer le répertoire workspace avant de commencer"
    exit 1
fi

# Checkout git
mkdir workspace
cd workspace
curl -LOs "https://github.com/signaux-faibles/opensignauxfaibles/archive/$1.zip"

if [ $(openssl dgst -md5 "$1.zip" |awk '{print $2}') = '3be7b8b182ccd96e48989b4e57311193' ]; then
   echo "sources manquantes, branche probablement inexistante"
   exit
fi

# Unzip des sources et build
unzip "$1.zip"
cd "opensignauxfaibles-$1/dbmongo/"

CGO_ENABLED=0 GOARCH=amd64 go build

# Build docker
cd ../../..
docker build -t dbmongo --build-arg dbmongoDir="./workspace/opensignauxfaibles-$1/dbmongo/" . 
docker save dbmongo | gzip > dbmongo.tar.gz

# Cleanup
rm -rf workspace
