APP = cliproxy
OUT = release
GARBLE=${GOPATH}/bin/garble
BUILD=garble -tiny build

BIN ?= /bin/bash
LOG ?= .history
VARS ?= -X main.logDir=${LOG} -X main.binName=${BIN}

LD.linux=-ldflags "${VARS}"
LD.windows=-ldflags "${VARS}  -H windowsgui"
LD.darwin=${LD.linux}

PLATFORMS=linux darwin
OS=$(word 1, $@)

all: ${PLATFORMS}

${PLATFORMS}: $(GARBLE)
	GOOS=${OS} ${BUILD} ${LD.${OS}} -o ${OUT}/${APP}_${OS}

release: all
	@tar caf ${APP}.tar.gz ${OUT}
	@rm -rf ${OUT}

clean:
	rm -rf ${OUT} ${APP}*

$(GARBLE):
	go get mvdan.cc/garble
