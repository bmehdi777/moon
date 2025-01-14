.PHONY: build clean theme test build-kc-theme build-agent-http clean-agent-http

build: build-agent build-server build-api-test

build-agent: build-agent-http 
	go build -o build/moon-agent cmd/agent/main.go
build-agent-http:
	npm run build --prefix assets/client/dashboard
	mv assets/client/dashboard/dist internal/pkg/agent/cmd/start/
build-server:
	go build -o build/moon-server cmd/server/main.go
build-api-test:
	go build -o build/api-test cmd/test-api/main.go
build-rpi:
	GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnu-gcc go build -o build/moon-server-rpi cmd/server/main.go

build-kc-theme: 
	npm run build-keycloak-theme --prefix keycloak/themes/moon-theme

run: run-agent run-server 

run-agent: build-client
	./build/moon-agent
run-server: build-server
	./build/moon-server
run-api-test:
	./build/api-test

certs:
	mkdir certs
	./scripts/makecert.sh

clean: clean-agent clean-server clean-api-test clean-agent-http

clean-agent:
	rm build/moon-agent
clean-agent-http:
	rm -r internal/pkg/agent/cmd/start/dist
clean-server:
	rm build/moon-server
clean-api-test:
	rm build/api-test
clean-rpi:
	rm build/moon-server-rpi

test:
	go test -v ./internal/pkg/communication

