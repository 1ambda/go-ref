# Backend 

## Requirements

| name | version | description |
|---|---|---|
| [go](https://github.com/golang/go) | 1.9.2+ | use [gvm](https://github.com/moovweb/gvm) |
| [dep](https://github.com/golang/dep) |  | `go get -u github.com/golang/dep/cmd/dep` |
| [swagger-go](https://github.com/google/protobuf) | 3.5.0+ | `go get -u github.com/go-swagger/go-swagger/cmd/swagger` |

```bash
make prepare
```

## Development

| cmd | arg | description |
| --- | --- | --- |
| make install | | install / update dependencies |
| make swagger | | generated swagger code files |
| make check | | format, style, lint |
| make build | | build applications |
| make docker | | build, tag, push image |
| make test | | run test |
| make run | APP=[*server*\|*client*] | |
| make run-cont | APP=[*server*\|*client*] | watch files and re-run |
