#!/bin/bash
{ echo -n `git describe --abbrev=0 --tags` ; echo -n '-' ; echo -n `git log -1 --pretty=format:%h` ; } | cat