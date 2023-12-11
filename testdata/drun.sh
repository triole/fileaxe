#!/bin/sh
fol="/tmp/log"

rcmd() {
  cmd=${@}
  echo -e "\033[0;93m${cmd}\033[0m"
  eval ${cmd}
}

rcmd mkdir -p "${fol}"

cd ${fol} && {
  rcmd dd if=/dev/random of=sample.log bs=2M count=1024
  rcmd logaxe "${fol}"
  ls -lah
}
