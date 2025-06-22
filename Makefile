BINARY=docker-cleaner
VERSION=1.0.0

build:
	go build -o bin/$(BINARY) main.go

install:
	cp bin/$(BINARY) /usr/local/bin/$(BINARY)

release:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY)-windows-amd64.exe

clean:
	rm -rf bin/*