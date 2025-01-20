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

## Development

#### Keycloak

Build keycloak theme : 
```
make build-kc-theme
```

#### Config 

Create a file `config.yml` in the root project folder : 

```yml
app:
  global_domain_name: m00n.fr
  http_addr: "0.0.0.0"
  http_port: "8080"

database:
  driver: sqlite
  sqlite_driver: "./moon.db"
  #driver: postgres
  #postgres_user: moon
  #postgres_password: moon
  #postgres_dbname: moon
  #postgres_port: 5432

auth:
  realm: "moon"
  base_url: "http://localhost:8081"
  algorithm: S256
  audience: "test"
```

#### Launch

- Launch the docker compose : docker compose up
- Launch the server : ./build/moon-server
- Launch the api test : ./build/api-test
- Login with the agent : ./build/moon-agent login
- Launch the agent : ./build/moon-agent start http://localhost:5000

## TODO

- Complete this readme : `docker run ...` isn't complete (which port open ? which env variable ? etc.)
- Explain how everything works (excalidraw architecture)
- Found logo

