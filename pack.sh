#!/bin/bash
VER=$1
if [ "$VER" = "" ]; then
    echo 'please input pack version!'
    exit 1
fi
RELEASE="release-${VER}"
rm -rf release-*
mkdir ${RELEASE}

# windows amd64
echo 'Start pack windows amd64...'
GOOS=windows GOARCH=amd64 go build ./
tar -czvf "${RELEASE}/gis-windows-amd64.tar.gz" gis.exe config.toml log/.gitignore LICENSE README.md
rm -rf gis.exe

echo 'Start pack windows X386...'
GOOS=windows GOARCH=386 go build ./
tar -czvf "${RELEASE}/gis-windows-386.tar.gz" gis.exe config.toml log/.gitignore LICENSE README.md
rm -rf gis.exe

echo 'Start pack linux amd64'
GOOS=linux GOARCH=amd64 go build ./
tar -czvf "${RELEASE}/gis-linux-amd64.tar.gz" gis config.toml log/.gitignore LICENSE README.md
rm -rf gis

echo 'Start pack linux 386'
GOOS=linux GOARCH=386 go build ./
tar -czvf "${RELEASE}/gis-linux-386.tar.gz" gis config.toml log/.gitignore LICENSE README.md
rm -rf gis

echo 'Start pack mac amd64'
GOOS=darwin GOARCH=amd64 go build ./
tar -czvf "${RELEASE}/gis-mac-amd64.tar.gz" gis config.toml log/.gitignore LICENSE README.md
rm -rf gis

echo 'END'
