TARGET=target

all: test build

test:
	go fmt ./...
	golint ./...
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	GOOS=linux go build -o $(TARGET)/app
	zip -j $(TARGET)/deployment.zip $(TARGET)/app

clean:
	rm -rf $(TARGET)

rebuild:
	clean all

doc:
	godoc -http=":6060"
