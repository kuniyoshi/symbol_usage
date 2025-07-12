.PHONY: all fmt vet staticcheck build run clean

all: fmt vet staticcheck build

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

build:
	go build -o sy .

run:
	go run main.go

clean:
	rm -f sy