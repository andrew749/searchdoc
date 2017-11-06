ARTIFACT_NAME=searchdoc

all: build

install_deps:
	go get -d ./...
build: install_deps
	go build -o $(ARTIFACT_NAME) **/*.go
