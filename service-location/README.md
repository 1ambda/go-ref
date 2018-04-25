# service-location 

## Requirements

| name | version | description |
|---|---|---|
| [go](https://github.com/golang/go) | 1.10+ | use [gvm](https://github.com/moovweb/gvm) |
| [dep](https://github.com/golang/dep) |  | `go get -u github.com/golang/dep/cmd/dep` |
| [protobuf](https://github.com/google/protobuf) | 3.5.0+ | `brew install protobuf` |

```bash
make prepare
```

## Development

| cmd | arg | description |
| --- | --- | --- |
| make install | | install / update dependencies |
| make protoc | | generated protobuf files |
| make check | | format, style, lint |
| make build | | build applications |
| make docker | | build, tag, push image |
| make test | | run test |
| make run | APP=[*server*\|*client*] | |
| make run-cont | APP=[*server*\|*client*] | watch files and re-run |
