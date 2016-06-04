#!/bin/sh
mkdir -p bin/mods/modswears
mkdir -p bin/mods/modchoice
cp -u ./mods-config-rename.json ./bin/mods/config.json
cp -u ./modswears-config-rename.json ./bin/mods/modswears/config.json
cp -u ./modchoice-config-rename.json ./bin/mods/modchoice/config.json
cp -u ./swears.txt ./bin/mods/modswears/swears.txt
go build -o ./bin/swbot.exe main.go