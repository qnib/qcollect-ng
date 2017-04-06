#!/bin/bash
set -e

LIB_PATH=resources/docker/lib
DIRS=$@
if [[ -z ${DIRS} ]];then
    DIRS="collectors filters handlers"
fi
for PTYPE in ${DIRS};do
    for PPATH in $(find "${PTYPE}" -mindepth 1 -maxdepth 1 -type d);do
        PLUGIN=$(basename ${PPATH})
        mkdir -p ${LIB_PATH}/${PTYPE}
        echo "> go build -ldflags \"-pluginpath=lib/${PTYPE}/${PLUGIN}\" -buildmode=plugin -o ${LIB_PATH}/${PTYPE}/${PLUGIN}.so ${PTYPE}/${PLUGIN}/plugin.go"
        go build -ldflags "-pluginpath=lib/${PTYPE}/${PLUGIN}" -buildmode=plugin -o ${LIB_PATH}/${PTYPE}/${PLUGIN}.so ${PTYPE}/${PLUGIN}/plugin.go
    done
done
