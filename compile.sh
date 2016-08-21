#!/bin/sh
./version.sh > ./bin/version.txt
mkdir -p ./bin/mods/modswears
cp -u ./swears.txt ./bin/mods/modswears/swears.txt
go build -o ./bin/swbot.exe main.go