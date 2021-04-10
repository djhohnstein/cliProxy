APP = cliproxy
OUT = release
BUILD=go build -trimpath

BIN ?= /bin/bash
LOG ?= .history
VARS ?= -X main.logDir=${LOG} -X main.binName=${BIN}

LD.linux=-ldflags "-s -w ${VARS}"
LD.windows=-ldflags "-s -w ${VARS}  -H windowsgui"
LD.darwin=${LD.linux}

PLATFORMS=linux darwin
OS=$(word 1, $@)

all: ${PLATFORMS}

${PLATFORMS}:
	GOOS=${OS} ${BUILD} ${LD.${OS}} -o ${OUT}/${APP}_${OS}

release: all
	@tar caf ${APP}.tar.gz ${OUT}
	@rm -rf ${OUT}

clean:
	rm -rf ${OUT} ${APP}*
