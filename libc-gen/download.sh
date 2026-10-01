#!/usr/bin/env bash
set -euo pipefail

cd -P -- "$(dirname -- "$0")"

cd c/
curl -#fOL "https://gee.cs.oswego.edu/pub/misc/malloc.c"
curl -#fOL "https://github.com/nothings/stb/raw/refs/heads/master/stb_sprintf.h"
