#!/bin/bash

scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir}" | grep -Po ".*(?=\/)")
testfol="${basedir}/testdata"

mkdir -p "${testfol}"

for i in {0..5}; do
    outf="${testfol}/logfile${i}.log"
    truncate -s 0 "${outf}"
    echo -e "\n\nWrite ${outf}"
    max=999
    for i in {000..999}; do
        printf '\r'
        printf "line ${i}/${max}"
        echo "${i} --- $(echo "${i}" | sha512sum)" >>"${outf}"
    done
done
