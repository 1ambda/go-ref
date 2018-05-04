include ../scripts/makefiles/variables.mk
include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/run.mk
include ../scripts/makefiles/test.mk
include ../scripts/makefiles/check.mk
include ../scripts/makefiles/git.mk

APP			= server
APPS		:= $(APP)

GITHUB_MODULE	= service-gateway
DOCKER_IMAGE	= $(shell echo $(DOCKER_IMAGE_PREFIX)_$(GITHUB_MODULE))

CONFIG_PACKAGE_PATH  = internal/config

SWAGGER_REST_SPEC	= ../schema/swagger/gateway-rest.yml
SWAGGER_WS_SPEC		= ../schema/swagger/gateway-websocket.yml

PROTO_SERVICES	:= hello location
SCHEMA_DIR	:= ../schema/protobuf

MOCKGEN = mockgen
MOCK_PREFIX = mock_
MOCK_PKG_DIRS := "internal/distributed" "internal/websocket"

.PHONY: dep
dep:
	@dep ensure -v

.PHONY: clean
clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning"
	@$(BAZEL) clean
	@echo ""

.PHONY: swagger-version
swagger-version:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Updating swagger spec version: $(VERSION)"
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_REST_SPEC) ../scripts/update-swagger-version.sh
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_WS_SPEC) ../scripts/update-swagger-version.sh

.PHONY: swagger-server
swagger-server: swagger-version
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning swagger files"
	rm -rf pkg/generated/swagger/* || true
	mkdir -p pkg/generated/swagger || true
	@echo ""

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Validating swagger (REST) spec"
	swagger validate $(SWAGGER_REST_SPEC)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating swagger (REST) code"
	swagger generate server --spec=$(SWAGGER_REST_SPEC) --exclude-main \
		--target=pkg/generated/swagger --model-package=rest_model \
		--server-package=rest_server --api-package=rest_api
	@echo ""

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Validating swagger (WS) spec"
	swagger validate $(SWAGGER_WS_SPEC)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating swagger (WS) code"
	swagger generate server --spec=$(SWAGGER_WS_SPEC) --exclude-main \
		--target=pkg/generated/swagger --model-package=ws_model \
		--skip-operations --skip-support --exclude-spec
	@echo ""

$(MOCK_PKG_DIRS):
	@$(eval TARGET := $@)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning mock files under $(TARGET)"
	rm -rf $(TARGET)/$(MOCK_PREFIX)*

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating mock files for $(TARGET)"
	PKG_DIR=$(TARGET) ../scripts/generate-gomock.sh
	@echo ""

.PHONY: mock-internal
mock-internal: $(MOCK_PKG_DIRS)

.PHONY: protoc
protoc:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning grpc files"
	@rm -rf pkg/generated/grpc/* || true
	@mkdir -p pkg/generated/grpc || true
	@echo ""

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating protobuf service files"
	protoc  --proto_path $(SCHEMA_DIR) --go_out=plugins=grpc:pkg/generated/grpc $(SCHEMA_DIR)/*.proto
	@echo ""

$(PROTO_SERVICES):
	@$(eval SERVICE := $@)

	# capitalize first letter
	@$(eval SERVICE_UPPERCASE := $(shell echo `echo $${SERVICE:0:1} | tr  '[a-z]' '[A-Z]'`$${SERVICE:1}))
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating protobuf mock files: $(SERVICE)"
	mockgen -source=pkg/generated/grpc/$(SERVICE).pb.go -package=pb $(SERVICE_UPPERCASE)Client > pkg/generated/grpc/mock_$(SERVICE).pb.go
	@echo ""

.PHONY: mock-proto
mock-proto: $(PROTO_SERVICES)

.PHONY: gazelle
gazelle:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating BUILD.bazel files"
	@ find vendor/ -name "BUILD" -delete
	@ find vendor/ -name "BUILD.bazel" -delete
	@$(BAZEL) run //:gazelle -- -proto disable
	@rm -f vendor/github.com/coreos/etcd/cmd/etcd
	@echo ""

.PHONY: inject-meta
inject-meta:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Injecting LDFLAGS into BUILD.bazel"
	@WORKSPACE=$(ROOT_IMPORT_PATH)/$(GITHUB_MODULE) CONFIG_PACKAGE_PATH=$(CONFIG_PACKAGE_PATH) ../scripts/get-ldflags.sh
	@echo ""

$(APPS): inject-meta
	@$(eval TARGET := $@)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Building app: $(TARGET)"
	@$(BAZEL) build //cmd/$(TARGET)
	@echo ""

.PHONY: build
build: mock-internal protoc mock-proto swagger-server gazelle $(APPS)

.PHONY: build-just
build-just: gazelle $(APPS)
