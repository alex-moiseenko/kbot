APP=$(shell basename $(shell git remote get-url origin))
REGESTRY=gcr.io/alex-m-devops-proj
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS ?= linux
TARGETARCH ?= arm64

format: 
	gofmt -s -w ./

lint:
	golint

test:
	go test -v

get:
	go get

build: format get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o kbot -ldflags "-X="github.com/alex-moiseenko/kbot/cmd.appVersion=${VERSION}

linux:
	linux-amd64 linux-arm64

linux-amd64:
	$(MAKE) GOOS=linux GOARCH=amd64 build

linux-arm64:
	$(MAKE) GOOS=linux GOARCH=arm64 build

macos:
	darwin-amd64 darwin-arm64

darwin-amd64:
	${MAKE} build TARGETOS=darwin TARGETARCH=amd64

darwin-arm64:
	${MAKE} build TARGETOS=darwin TARGETARCH=arm64

windows:
	$(MAKE) GOOS=windows GOARCH=amd64 build

image:
	docker build . -t ${REGESTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH} --build-arg="TARGETOS=${TARGETOS}" --build-arg="TARGETARCH=${TARGETARCH}"

image-linux-arm64:
	$(MAKE) TARGETOS=linux TARGETARCH=arm64 image

image-linux-amd64:
	$(MAKE) TARGETOS=linux TARGETARCH=amd64 image

image-darwin-arm64:
	$(MAKE) TARGETOS=darwin TARGETARCH=arm64 image

image-darwin-amd64:
	$(MAKE) TARGETOS=darwin TARGETARCH=amd64 image

image-windows-amd64:
	$(MAKE) TARGETOS=windows TARGETARCH=amd64 image

image-macos:
	image-darwin-amd64 image-darwin-arm64

image-linux:
	image-linux-amd64 image-linux-arm64

push:
	docker push ${REGESTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}

clean:
	rm -rf kbot
	docker rmi ${REGESTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH} || true