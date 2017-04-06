.PHONY: libs 

SRCDIR=/usr/local/src/github.com/qnib/qcollect-ng
DOCKERIMG=qnib/uplain-golang

all: gov-fetch gov-remove libs binary


libs: gov-update collectors filters handlers
gov-update:
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor update +l
gov-fetch:
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor fetch +e +m +v
gov-remove:
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor remove +u
test:
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./bin/test.sh
collectors: gov-update
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./bin/build.sh collectors
filters: gov-update
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./bin/build.sh filters
handlers: gov-update
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./bin/build.sh handlers
binary: libs
	docker run -t --rm -e SKIP_ENTRYPOINTS=1 -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) go build -o resources/docker/bin/qcollect-ng
docker:
	docker build -t qnib/qcollect-ng:local resources/docker/.
