ARTIFACT_NAME=searchdoc

build:
	go build -o $(ARTIFACT_NAME) **/*.go
