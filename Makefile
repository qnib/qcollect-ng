.PHONY: libs 

SRCDIR=/usr/local/src/github.com/qnib/qcollect-ng
DOCKERIMG=qnib/uplain-golang

all: gov-fetch gov-remove libs

libs: collectors filters handlers
gov-update:
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor update +l ; else govendor update +l ; fi
gov-fetch:
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor fetch +e +m ; else govendor fetch +e +m ; fi
gov-remove:
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) govendor remove +u ; else govendor remove +u ; fi
collectors: gov-update
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./build.sh collectors ; else ./build.sh collectors ; fi
filters: gov-update
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./build.sh filters ; else ./build.sh filters ; fi
handlers: gov-update
	if test -z "$$local"; then docker run -t --rm -v ${CURDIR}:$(SRCDIR) -w $(SRCDIR) $(DOCKERIMG) ./build.sh handlers ; else ./build.sh handlers ; fi
