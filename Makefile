.PHONY: all test build clean

all: clean test build

build: 
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  quay.io/luigizuccarelli/golang-openshift-api:1.16.6 .

push:
	podman push quay.io/luigizuccarelli/golang-openshit-api:1.16.6 
