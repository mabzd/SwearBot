#!/bin/sh
mkdir -p bin/mods/modswears
cp -u ./modswears-config-rename.json ./bin/mods/modswears/config.json
cp -u ./swears.txt ./bin/mods/modswears/swears.txt
go build -o ./bin/swbot.exe main.go