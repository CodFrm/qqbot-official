
GOOS=linux
GOARCH=amd64
VERSION=v0.1.0

PROJECT_NAME=guild-bot

NAME=$(PROJECT_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/$(PROJECT_NAME)
SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

test:
	GOOS=$(GOOS) go test -v ./...

generate:
	go generate ./... -x

build: generate
	go build -o $(PROJECT_NAME)$(SUFFIX) ./cmd/app

target:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(NAME)$(SUFFIX) ./cmd/app
