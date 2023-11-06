BINARY_NAME=goperson

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ./cmd/go/ 

run: build
	./bin/${BINARY_NAME}-linux

clean:
	go clean
	rm ./bin/${BINARY_NAME}-linux

dep:s
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all