TARGET=target

all: test build

# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go fmt ./...
	golint ./...
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	GOOS=linux go build -o $(TARGET)/app cmd/main.go
	zip -j $(TARGET)/deployment.zip $(TARGET)/app

proto:
	protoc -I internal/http/grpc/ internal/http/grpc/*.proto --go_out=plugins=grpc:internal/http/grpc

clean:
	rm -rf $(TARGET)

rebuild:
	clean all

doc:
	godoc -http=":6060"
