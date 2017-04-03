.PHONY: libs 

all: gov-fetch gov-remove libs

libs: collectors filters handlers
gov-update:
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang govendor update +l
gov-fetch:
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang govendor fetch +e +m
gov-remove:
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang govendor remove +u
collectors: gov-update
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang ./build.sh collectors
filters: gov-update
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang ./build.sh filters
handlers: gov-update
	docker run -t --rm -v ${CURDIR}:/usr/local/src/github.com/qnib/qcollect-ng -w /usr/local/src/github.com/qnib/qcollect-ng qnib/uplain-golang ./build.sh handlers
