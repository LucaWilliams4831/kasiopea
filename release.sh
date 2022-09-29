#!/usr/bin/env bash

rm -rf _output/*.zip

echo "build for windows ..."
./build.sh win
cd _output/
zip -q -r kasiopea_windows_amd64_$1.zip kasiopea

cd ..
echo "build for linux ..."
./build.sh linux
cd _output/
zip -q -r kasiopea_linux_amd64_$1.zip kasiopea

cd ..
echo "build for mac ..."
./build.sh mac
cd _output/
zip -q -r kasiopea_macos_amd64_$1.zip kasiopea

echo "end ..."