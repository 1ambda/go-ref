[![Go Report Card](https://goreportcard.com/badge/github.com/1ambda/go-ref?style=flat-square)](https://goreportcard.com/report/github.com/1ambda/go-ref)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/1ambda/go-ref)
[![Release](https://img.shields.io/github/release/1ambda/go-ref.svg?style=flat-square)](https://github.com/1ambda/go-ref/releases/latest)

# Golang Reference Project

## Component
- [service-front](https://github.com/1ambda/go-ref/tree/master/service-front): webapp
  * generate rest, websocket models from [swagger spec](https://github.com/1ambda/go-ref/tree/master/schema/swagger) using [swagger-codegen](https://github.com/swagger-api/swagger-codegen)
  - written in Angular 5+
- [service-gateway](https://github.com/1ambda/go-ref/tree/master/service-gateway): serving Websocket and REST endpoints
  * generate rest, websocket model + server stub using [go-swagger](https://github.com/go-swagger/go-swagger)
  * while working as a gRPC client
  * provides scalable realtime metric servers using etcd 
- [service-location](https://github.com/1ambda/go-ref/tree/master/service-backend): serving gRPC endpoints
  * provides location service sharding (e.g location server 1 handles requests from US, ..)
  


```                                                                                                  
                                                                                                     
    +---------------+   WS   (50001)  +-----------------+                  +------------------+       
    |               |   REST (50002)  |                 |   gRPC (50003)   |                  |       
    | service-front |<--------------->| service-gateway |<---------------->| service-location |       
    |               |  (swagger-gen)  |                 |    (protobuf)    |                  |       
    +---------------+                 +-----------------+                  +------------------+       
         (webapp)                           (server)                             (server)            
                                                                                                     
                                                              < storage >                                            
                                                                                                     
                                                   +--------------+  +-------------+                 
                                                   |              |  |             |                 
                                                   | etcd cluster |  |    mysql    |                 
                                                   |              |  |             |                 
                                                   +--------------+  +-------------+                 
                                                   (distributed kv)  (transactional)                 
                                                                                       
```

## Requirements
| name | version | description |
|---|---|---|
| [go](https://github.com/golang/go) | 1.10+ | use [gvm](https://github.com/moovweb/gvm) |
| [dep](https://github.com/golang/dep) |  | `go get -u github.com/golang/dep/cmd/dep` |
| [protobuf](https://github.com/google/protobuf) | 3.5.0+ | `brew install protobuf` |
| [nodejs](https://nodejs.org/) | 9.8.0+ | use [nvm](https://github.com/creationix/nvm) |

Please check [make prepare](https://github.com/1ambda/go-ref/blob/master/scripts/makefiles/install.mk#L6) command for more description.

## Quick Start

```bash
# start dockerized storages: etcd, mysql, ...
$ docker-compose up
```

```bash
# start gateway service 
$ cd service-gateway
$ make prepare install
$ make run 
```

```bash
# start front (webapp) service 
$ cd service-front
$ npm install; npm run server:local;

# visit: https://localhost:3000 
```

## Screenshot

<img src="https://raw.githubusercontent.com/1ambda/go-ref/master/docs/screenshots/screenshot1.png" height="100%" width="100%">
<img src="https://raw.githubusercontent.com/1ambda/go-ref/master/docs/screenshots/screenshot2.png" height="100%" width="100%">

