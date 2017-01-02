#!/bin/bash

## clean
rm -rf build/

## for Mac
(
  GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/pw cmd/pw/*.go
  cd ./build/darwin-amd64 || exit
  tar cfz pw-darwin-amd64.tar.gz pw
)

## for windows
(
  GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/pw.exe cmd/pw/*.go
  cd ./build/windows-amd64 || exit
  zip -q pw-windows-amd64.tar.gz pw.exe
)

## for linux
(
  GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/pw cmd/pw/*.go
  cd ./build/linux-amd64 || exit
  tar cfz pw-linux-amd64.tar.gz pw
)
