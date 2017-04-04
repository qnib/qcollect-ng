#!/usr/bin/env bash

PROJECT=qcollect-ng
if [ ! -d resources/coverity ];then
    mkdir -p resources/coverity
fi
echo "> go test -cover -coverprofile=resources/coverity/${PROJECT}.cover ."
go test -cover -coverprofile=resources/coverity/${PROJECT}.cover . >>/dev/null
COVER_FILES="resources/coverity/${PROJECT}.cover"
for x in $(find . -type d |egrep -v "(\.$|\.git|vendor|bin|lib|resources|.idea|(collectors|filters|handlers)$)");do
    echo "> go test -cover -coverprofile=resources/coverity/${x}.cover ${x}"
    mkdir -p resources/coverity/$(dirname ${x})
    go test -cover -coverprofile=resources/coverity/${x}.cover ${x} >>/dev/null
    COVER_FILES="${COVER_FILES} resources/coverity/${x}.cover"
done
echo "> coveraggregator -o resources/coverity/coverage-all.out ${COVER_FILES}"
coveraggregator -o resources/coverity/coverage-all.out ${COVER_FILES} >>/dev/null
echo "> create HTML report"
go tool cover -html resources/coverity/coverage-all.out -o resources/coverity/report.html
