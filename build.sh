#!/bin/bash
# build.sh compile les éléments et dispose le nécessaire dans le répertoire dist.
# ce script ne se charge pas d'installer les éléments nécessaires à la compilation.
mkdir -p dist/static

# compilation et recopie des éléments nécessaires
go build -ldflags="-s -w"
upx sfdata
cp sfdata ../dist
cp config.toml ../dist

# création de l'archive
tar cvJf dist.tar.xz dist/

# cleanup
rm sfdata
rm -r dist/
