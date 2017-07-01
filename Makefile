all: build

deps: # sets up build dependencies, needed to do once during development
	go get -u github.com/golang/dep/...
	dep ensure

build: # builds the cli program
	go build -o sawzall ./cmd/sawzall/main.go

test: # runs the unit tests
	go test -v -race -cover ./pkg/sawzall
