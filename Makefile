PROG_NAME := "http-honeypot"
IMAGE_NAME := "pschou/http-honeypot"
VERSION := "0.1"


build:
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}'" -o ${PROG_NAME} main.go

docker: build
	docker build -f Dockerfile --tag ${IMAGE_NAME}:${VERSION} .
