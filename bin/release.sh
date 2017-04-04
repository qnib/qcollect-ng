#!/bin/bash

CURVERSION=$(cat .bumpversion.cfg |awk -F\= '/current_version/{print $2}' |sed -e 's/ //g')
mkdir -p bin/amd64
go build -o bin/amd64/qcollect-ng
cp bin/amd64/qcollect-ng bin/amd64/qcollect-ng_${CURVERSION}

