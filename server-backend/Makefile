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
GOTEST		= ginkgo
GOTEST_OPT	= -r -p -v
GO_FILES	= $(shell $(GOCMD) list ./... | grep -v /vendor/)

APP			= server					# default app
APPS		:= $(APP) client

# `api/$(SERVICE).proto					: protobuf definition
# `pkg/api/$(SERVICE).pb.go`			: (auto-generated) protobuf service
# `pkg/api/mock_$(SERVICE).pb.go`	: (auto-generated) protobuf mock
SERVICES	:= hello
SCHEMA_DIR	:= ../schema/protobuf

GIT_COMMIT	= $(shell git rev-parse HEAD)
GIT_DIRTY	= $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
VERSION		= $(shell cat ./VERSION)

GITHUB_USER		= 1ambda				# github username
GITHUB_REPO		= go-ref				# github repository
GITHUB_MODULE	= server-backend		# module name in the repository
DOCKER_REPO		= 1ambda
DOCKER_IMAGE	= go-ref_server-backend

prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - installing prerequisites"
	$(GOGET) github.com/alecthomas/gometalinter
	$(GOLINT) --install --update --force
	$(GOGET) github.com/ahmetb/govvv
	$(GOGET) github.com/githubnemo/CompileDaemon
	$(GOGET) github.com/onsi/ginkgo/ginkgo

	$(GOGET) github.com/golang/protobuf/protoc-gen-go
	$(GOGET) github.com/golang/mock/mockgen

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
	protoc  --proto_path $(SCHEMA_DIR) $(SERVICE).proto --go_out=plugins=grpc:pkg/api
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating protobuf mock files: $(SERVICE)"
	mockgen -source=pkg/api/$(SERVICE).pb.go -package=api $(SERVICE_UPPERCASE)Client > pkg/api/mock_$(SERVICE).pb.go

protoc: $(SERVICES)

$(APPS): protoc
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

build: clean protoc $(APPS)

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

docker-prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - setting GOOS to linux for docker image build"
	@$(eval GOOS := linux)

docker-image: docker-prepare docker-clean build
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building docker image: $(APP)"
	docker build \
		--build-arg APP=${APP} \
		-t $(DOCKER_IMAGE):local .

docker-tag:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - tagging docker image: latest $(VERSION) $(GIT_COMMIT)"
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT)
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION)
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):latest

docker-push:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - pushing docker image: latest $(VERSION) $(GIT_COMMIT)"
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT)
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):${VERSION}
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):latest

docker-clean:
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT) || true
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION) || true
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):latest || true
	docker image rm $(DOCKER_IMAGE):local || true

docker: clean docker-image docker-tag docker-push

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

all: clean install format lint style build test

.PHONY: clean install build
