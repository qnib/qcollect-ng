#!/bin/bash
set -e

DIRS=$@
if [[ -z ${DIRS} ]];then
    DIRS="collectors filters handlers"
fi
for PTYPE in ${DIRS};do
    for PPATH in $(find "${PTYPE}" -mindepth 1 -maxdepth 1 -type d);do
        PLUGIN=$(basename ${PPATH})
        mkdir -p lib/${PTYPE}
        echo "> go build -ldflags \"-pluginpath=lib/${PTYPE}/\" -buildmode=plugin -o lib/${PTYPE}/${PLUGIN}.so ${PTYPE}/${PLUGIN}/plugin.go"
        go build -ldflags "-pluginpath=lib/${PTYPE}/" -buildmode=plugin -o lib/${PTYPE}/${PLUGIN}.so ${PTYPE}/${PLUGIN}/plugin.go
    done
done
