TARGET=target
VERSION=v0.0.3

all: test build-linux build-mac build-windows

# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...
	go test ./...

integration:
	go test ./... -tags=integration

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build-linux:
	GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/linux/app cmd/main.go
	zip -j $(TARGET)/deployment-linux.zip $(TARGET)/linux/app

build-mac:
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/mac/app cmd/main.go
	zip -j $(TARGET)/deployment-mac.zip $(TARGET)/mac/app

build-windows:
	GOOS=windows go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/windows/app cmd/main.go
	zip -j $(TARGET)/deployment-windows.zip $(TARGET)/windows/app

proto:
	protoc -I internal/http/grpc/ internal/http/grpc/*.proto --go_out=plugins=grpc:internal/http/grpc

clean:
	rm -rf $(TARGET)

rebuild:
	clean all

doc:
	godoc -http=":6060"
