build:
	go build -o gproc cmd/main.go

install:
	go mod tidy
	go build -o gproc cmd/main.go

clean:
	rm -f gproc gproc.exe
	rm -rf logs/*

test:
	go test ./...

run-example:
	./gproc start test-app ping google.com

.PHONY: build install clean test run-example