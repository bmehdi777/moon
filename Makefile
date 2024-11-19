.PHONY: build clean

build: build-agent build-server build-api-test

build-agent:
	go build -o build/moon-agent cmd/agent/main.go
build-server:
	go build -o build/moon-server cmd/server/main.go
build-api-test:
	go build -o build/api-test cmd/test-api/main.go
build-rpi:
	GOOS=linux GOARCH=arm go build -o build/moon-server-rpi cmd/server/main.go

run: run-agent run-server 

run-agent: build-client
	./build/moon-agent
run-server: build-server
	./build/moon-server
run-api-test:
	./build/api-test

clean: clean-agent clean-server clean-api-test

clean-agent:
	rm build/moon-agent
clean-server:
	rm build/moon-server
clean-api-test:
	rm build/api-test
clean-rpi:
	rm build/moon-server-rpi
