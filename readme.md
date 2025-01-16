## Getting started

### Agent

#### Go

```sh
go install moon/cmd/agent@latest
```

#### Manual

First, you need to install (Go)[https://go.dev/doc/install].
This project require at least __go 1.23.2__.

Then, you can clone and compile everything :

```sh
git clone https://moon
cd moon
make build-client
```

A binary named `moon-agent` will be located in the `build/` folder.

### Server

> This project use `keycloak` to authenticate users. You will need to install it
> on your own or use the `compose.yml` file

#### Docker

The server image can be build manually :

```sh
docker build -t moon-server .
```

or 

```sh
docker run bmehdi777/moon-server@latest
```

## TODO

- Complete this readme : `docker run ...` isn't complete (which port open ? which env variable ? etc.)
- Explain how everything works (excalidraw architecture)
- Found logo

