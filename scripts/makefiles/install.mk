GOCMD		= go
GOGET		= $(GOCMD) get -u -v
GODEP		= dep
BREW		= brew

prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing prerequisites"
	$(GOGET) github.com/alecthomas/gometalinter
	$(GOLINT) --install --update --force
	$(GOGET) github.com/ahmetb/govvv
	$(GOGET) github.com/githubnemo/CompileDaemon
	$(GOGET) github.com/onsi/ginkgo/ginkgo

	# protobuf
	$(GOGET) github.com/golang/protobuf/protoc-gen-go
	$(GOGET) github.com/golang/mock/mockgen

	# swagger
	$(GOGET) github.com/go-swagger/go-swagger/cmd/swagger
	$(BREW) install swagger-codegen

	$(GOGET) github.com/bcicen/grmon/cmd/grmon

install:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing / updating dependencies"
	@$(GODEP) ensure -update
	@$(GODEP) ensure

