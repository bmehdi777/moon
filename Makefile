.PHONY: build clean

build: build-client build-server build-api-test

build-client:
	go build -o build/moon-client cmd/client/main.go
build-server:
	go build -o build/moon-server cmd/server/main.go
build-api-test:
	go build -o build/api-test cmd/test-api/main.go

run: run-client run-server 

run-client: build-client
	./build/moon-client
run-server: build-server
	./build/moon-server
run-api-test:
	./build/api-test

clean: clean-client clean-server clean-api-test

clean-client:
	rm build/moon-client
clean-server:
	rm build/moon-server
clean-api-test:
	rm build/api-test
