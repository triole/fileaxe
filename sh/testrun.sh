#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir%/*}")

r "${basedir}/testdata/tmp" -r ".*" --remove -m 1m
