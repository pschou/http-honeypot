PROG_NAME := "http-honeypot"
IMAGE_NAME := "pschou/http-honeypot"
VERSION = 0.1.$(shell date -u +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"


build:
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o ${PROG_NAME} main.go
	upx --lzma ${PROG_NAME}


docker: build
	docker build -f Dockerfile --tag ${IMAGE_NAME}:${VERSION} .
