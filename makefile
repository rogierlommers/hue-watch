BINARY = hue-watch
SOURCE := *.go
BUILD_DIR=$(shell pwd)/binary

all: clean linux

linux: 
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/hue-watch-amd64 .

clean:
	rm -rf binary/*

container: linux
	docker build -t rogierlommers/hue-watch .
	docker push rogierlommers/hue-watch:latest
