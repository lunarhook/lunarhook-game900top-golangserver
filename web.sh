#!/bin/sh
export PATH=$PATH:/usr/local/go/bin
git pull
cd src/main/
go build -o game900top main.go
cd ..
cd ..

ps -ef |grep game900top | grep -v grep |awk '{print $2}' | xargs kill -9
cp -Rf ./src/main/game900top ./game900top
nohup ./game900top >/dev/null 2>&1 &