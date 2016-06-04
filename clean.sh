#!/bin/bash
find ./bin -type f -not -name 'token.txt' -maxdepth 1 | xargs rm
find ./bin -type d -not -name 'bin' -maxdepth 1 | xargs rm -R