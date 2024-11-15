.PHONY: build clean

build: build-client build-server

build-client:
	go build -o build/moon-client cmd/client/main.go
build-server:
	go build -o build/moon-server cmd/server/main.go

run: run-client run-server

run-client: build-client
	./build/moon-client
run-server: build-server
	./build/moon-server

clean: clean-client clean-server

clean-client:
	rm build/moon-client
clean-server:
	rm build/moon-server
