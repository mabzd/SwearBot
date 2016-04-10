#!/bin/sh
mkdir -p bin
cp -u ./config-rename.json ./bin/config.json
cp -u ./swears.txt ./bin/swears.txt
go build -o ./bin/swbot.exe main.go