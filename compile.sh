#!/bin/sh
{ echo -n `git describe --abbrev=0 --tags` ; echo -n '-' ; echo -n `git log -1 --pretty=format:%h` ; } > ./bin/version.txt
mkdir -p ./bin/mods/modswears
cp -u ./swears.txt ./bin/mods/modswears/swears.txt
go build -o ./bin/swbot.exe main.go