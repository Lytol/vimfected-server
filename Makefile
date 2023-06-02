BINARY_NAME=vimfected

all: build test run

build:
	go build -o ${BINARY_NAME} ./cmd/vimfected

test:
	go test ./...

run:
	@./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}