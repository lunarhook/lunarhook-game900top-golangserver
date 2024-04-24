#!/bin/sh
export PATH=$PATH:/usr/local/go/bin
git pull
cd server/Global/main/
go build -o Russiamain main.go 
cd ..
cd ..
cd ..

ps -ef |grep Russiamain | grep -v grep |awk '{print $2}' | xargs kill -9
#cp -Rf ./server/Global/main/Russiamain ./Russiamain
nohup ./Russiamain >Russiamain.out 2>&1 &