export GO111MODULE=on
BIN_NAME="OAuth_client"

all:clean build run

build:
	go mod tidy
	bash -x build.sh

run:
	./output/bin/${BIN_NAME}

clean:
	rm -rf output

test:
	go test -v