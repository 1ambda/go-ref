GOARCH  = amd64
GOOS    = $(OS)

ifeq ($(GOOS),)
  ifeq ($(shell  uname -s), Darwin)
    GOOS	= darwin
  else
    GOOS	= linux
  endif
endif

TAG			= make
MAIN		= main.go
BIN_DIR		= bin
CMD_DIR		= cmd
VENDOR_DIR	= vendor

GOCMD		= go
GOGET		= $(GOCMD) get -u -v
GODEP		= dep
GOVVV		= govvv
GOLINT		= gometalinter
GOBUILD		= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build
GOFMT		= $(GOCMD)fmt
GOTEST	= ginkgo
GOTEST_OPT	= -r -p -v

GO_FILES	= $(shell $(GOCMD) list ./... | grep -v /vendor/)

APP			= server					# default app
APPS		:= server client

# `api/$(SERVICE).proto					: protobuf definition
# `pkg/api/$(SERVICE).pb.go`			: (auto-generated) protobuf service
# `internal/mock/mock_$(SERVICE).go`	: (auto-generated) protobuf mock
SERVICES	:= hello

prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing prerequisites"
	brew update
	brew install dep
	brew install protobuf
	$(GOGET) github.com/golang/protobuf/protoc-gen-go
	$(GOGET) github.com/golang/mock/mockgen
	$(GOGET) github.com/alecthomas/gometalinter
	$(GOLINT) --install --update --force
	$(GOGET) github.com/ahmetb/govvv
	$(GOGET) github.com/githubnemo/CompileDaemon
	$(GOGET) github.com/onsi/ginkgo/ginkgo

install:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing / updating dependencies"
	@$(GODEP) ensure -update
	@$(GODEP) ensure

lint:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - linting"
	@$(GOLINT) --vendor --errors ./... --skip=internal/mock

style:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - checking code style"
	@! $(GOFMT) -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - formatting code"
	@$(GOCMD) fmt $(GO_FILES)

check: format style lint

$(SERVICES):
	@$(eval SERVICE := $@)
	# capitalize first letter
	@$(eval SERVICE_UPPERCASE := $(shell echo `echo $${SERVICE:0:1} | tr  '[a-z]' '[A-Z]'`$${SERVICE:1}))
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating protobuf service files: $(SERVICE)"
	protoc api/$(SERVICE).proto --go_out=plugins=grpc:pkg
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating protobuf mock files: $(SERVICE)"
	mockgen -source=pkg/api/$(SERVICE).pb.go -package=mock $(SERVICE_UPPERCASE)Client > internal/mock/mock_$(SERVICE).go

protoc: $(SERVICES)

$(APPS): protoc
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

build: protoc $(APPS)

test: check
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - testing"
	$(GOTEST) $(GOTEST_OPT) -skip=$(VENDOR_DIR) --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace

test-cont:
	@echo "\n[MAKEFILE] ($(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')) test-cont"
	$(GOTEST) watch $(GOTEST_OPT) -skip=$(VENDOR_DIR)

run: build
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - running app: $(APP)"
	@$(BIN_DIR)/$(APP)

run-cont:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="${BIN_DIR}" \
		-exclude-dir="internal/mock" \
		-exclude-dir="pkg/api" \
		-build="make build" -command="make run" -graceful-kill=true

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"

all: clean install format lint style build test

.PHONY: clean install build
