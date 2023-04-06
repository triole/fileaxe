#!/bin/bash

scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir}" | grep -Po ".*(?=\/)")
testfol="${basedir}/testdata"

mkdir -p "${testfol}"

for i in {0..5}; do
    outf="${testfol}/logfile${i}.log"
    truncate -s 0 "${outf}"
    for i in {00000..99999}; do
        echo "${i} --- $(echo "${i}" | sha512sum)" >>"${outf}"
    done
done
