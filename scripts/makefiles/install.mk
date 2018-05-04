GOCMD		= go
GOGET		= $(GOCMD) get -u -v
BREW		= brew

prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing prerequisites"
	$(GOGET) github.com/alecthomas/gometalinter
	$(GOLINT) --install --update --force
	$(GOGET) github.com/ahmetb/govvv
	$(GOGET) github.com/githubnemo/CompileDaemon

	# testing framework
	$(GOGET) github.com/onsi/ginkgo/ginkgo
	$(GOGET) github.com/onsi/gomega/...

	# protobuf
	$(GOGET) github.com/golang/protobuf/protoc-gen-go
	$(GOGET) github.com/golang/mock/mockgen

	# swagger
	$(GOGET) github.com/go-swagger/go-swagger/cmd/swagger
	$(BREW) install swagger-codegen

	# etcd cli: https://github.com/coreos/etcd/tree/master/etcdctl
	$(GOGET) github.com/coreos/etcd/etcdctl

	# mock
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

	# bazel
	$(BREW) brew install bazel
	$(GOGET) get github.com/bazelbuild/buildtools/buildifier
	$(GOGET) go get github.com/bazelbuild/buildtools/buildozer

