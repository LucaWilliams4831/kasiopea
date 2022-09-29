#!/usr/bin/env bash

# ./build.sh mac|win|linux

if [ ! -d "_output/" ];then
mkdir _output
fi
if [ ! -d "_output/kasiopea" ];then
mkdir _output/kasiopea
else
rm -rf _output/kasiopea/*
fi

#if [ -d ".conf/" ];then
#cp .conf/sipt.yaml _output/kasiopea.yaml
#else
cp example.yaml _output/kasiopea/example.yaml
#fi
cd assets
go generate
cd ..

if [ "$1" == "mac" ];then
# mac os
GOOS=darwin GOARCH=amd64 go build -tags release -o _output/kasiopea/kasiopea ./cmd/
GOOS=darwin GOARCH=amd64 go build -o _output/kasiopea/upgrade scripts/upgrade.go
`echo "#!/usr/bin/env bash
nohup ./kasiopea >> kasiopea.log 2>&1 &" > _output/kasiopea/start.sh`
`chmod a+x _output/kasiopea/start.sh`
elif [ "$1" == "win" ];then
# windows
GOOS=windows GOARCH=amd64 go build -tags release -o _output/kasiopea/kasiopea.exe ./cmd/
GOOS=windows GOARCH=amd64 go build -o _output/kasiopea/upgrade.exe scripts/upgrade.go
`echo "@echo off
if \"%1\" == \"h\" goto begin
mshta vbscript:createobject(\"wscript.shell\").run(\"%~nx0 h\",0)(window.close)&&exit
:begin
kasiopea >> kasiopea.log" > _output/kasiopea/startup.bat`
elif [ "$1" == "linux" ];then
# linux
GOOS=linux GOARCH=amd64 go build -tags release -o _output/kasiopea/kasiopea ./cmd/
GOOS=linux GOARCH=amd64 go build -o _output/kasiopea/upgrade scripts/upgrade.go
`echo "#!/usr/bin/env bash
nohup ./kasiopea >> kasiopea.log 2>&1 &" > _output/kasiopea/start.sh`
`chmod a+x _output/kasiopea/start.sh`
fi
