[![Go Report Card](https://goreportcard.com/badge/github.com/1ambda/go-ref?style=flat-square)](https://goreportcard.com/report/github.com/1ambda/go-ref)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/1ambda/go-ref)
[![Release](https://img.shields.io/github/release/1ambda/go-ref.svg?style=flat-square)](https://github.com/1ambda/go-ref/releases/latest)

# Golang Reference Project

## Component
- [service-gateway](https://github.com/1ambda/go-ref/tree/master/service-gateway): App serving REST endpoints
  * generate model using [go-swagger](https://github.com/go-swagger)
  * call backend gRPC server (*TODO*)
- [service-backend](https://github.com/1ambda/go-ref/tree/master/service-backend): App serving gRPC endpoints
  * persist data into storage (*TODO*)
- [service-front](https://github.com/1ambda/go-ref/tree/master/service-front): App web dashboard

- [etcd](https://github.com/coreos/etcd): Distributed reliable key-value store

## build & start project

### Requirements
| name | version | description |
|---|---|---|
| [go](https://github.com/golang/go) | 1.9.2+ | use [gvm](https://github.com/moovweb/gvm) |
| [dep](https://github.com/golang/dep) |  | `go get -u github.com/golang/dep/cmd/dep` |
| [protobuf](https://github.com/google/protobuf) | 3.5.0+ | `brew install protobuf` |
| [nodejs](https://nodejs.org/) | 9.8.0+ | use [nvm](https://github.com/creationix/nvm) |
| [Docker](https://www.docker.com/) | | |

### build & start
- start etcd, etcd-ui, MySQL using docker
```bash
$ cd {go-ref-dir}
$ docker-compose up
```


- start service-gateway
```bash
$ make prepare install
```

- start service-front
```bash
$ npm i; npm run server:dev;
```

### Web UI
- etcd-ui
```bash
http://localhost:2381/
```
- service-front UI
```bash
http://localhost:3000/
```

  
