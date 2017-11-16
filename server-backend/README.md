# Backend 

## Requirements

| name | version | description |
|---|---|---|
| [go](https://github.com/golang/go) | 1.9.2+ | use [gvm](https://github.com/moovweb/gvm) |
| [dep](https://github.com/golang/dep) |  | `go get -u github.com/golang/dep/cmd/dep` |
| [protobuf](https://github.com/google/protobuf) | 3.5.0+ | `brew install protobuf` |

```bash
make prepare
```

## Development

| cmd | arg | description |
| --- | --- | --- |
| make prepare | | install prerequisites |
| make install | | install / update dependencies |
| make protoc | | generated protobuf files |
| make lint | | [gometalinter](https://github.com/alecthomas/gometalinter) |
| make style | | check code style |
| make format | | format code style |
| make build | | build applications |
| make test | | run test |
| make test-cont | | watch files and re-run |
| make run | APP=[*server*\|*client*] | |
| make run-cont | APP=[*server*\|*client*] | watch files and re-run |
