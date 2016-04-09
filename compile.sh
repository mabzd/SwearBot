#!/bin/sh
mkdir -p bin
cp -n ./config.json ./bin/config.json
cp -n ./swears.txt ./bin/swears.txt
go build -o ./bin/swbot.exe main.go