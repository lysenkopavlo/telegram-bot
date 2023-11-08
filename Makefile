BINARY_NAME=telegram_bot

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux 

run: build
	./bin/${BINARY_NAME}-linux

clean:
	go clean
	rm ./bin/${BINARY_NAME}-linux

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all