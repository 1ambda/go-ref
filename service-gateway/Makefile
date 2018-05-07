include ../scripts/makefiles/variables.mk
include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/check.mk
include ../scripts/makefiles/git.mk
include ../scripts/makefiles/protobuf.mk
include ../scripts/makefiles/bazel.mk

APP			= server
APPS		:= $(APP)

GITHUB_MODULE	= service-gateway
DOCKER_IMAGE	= $(shell echo $(DOCKER_IMAGE_PREFIX)_$(GITHUB_MODULE))

CONFIG_PACKAGE_PATH  = internal/config

SWAGGER_REST_SPEC	= ../schema/swagger/gateway-rest.yml
SWAGGER_WS_SPEC		= ../schema/swagger/gateway-websocket.yml

PROTO_SERVICES	:= hello location
SCHEMA_DIR	:= ../schema/protobuf

MOCKGEN 	= mockgen
MOCK_PREFIX = mock
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
	@rm -rf $(TARGET)/$(MOCK_PREFIX)*

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating mock files for $(TARGET)"
	@PKG_DIR=$(TARGET) ../scripts/generate-gomock.sh
	@echo ""

.PHONY: mock-internal
mock-internal: $(MOCK_PKG_DIRS)

$(PROTO_SERVICES):
	@$(eval SERVICE := $@)

	# capitalize first letter
	@$(eval SERVICE_UPPERCASE := $(shell echo `echo $${SERVICE:0:1} | tr  '[a-z]' '[A-Z]'`$${SERVICE:1}))
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating protobuf mock files: $(SERVICE)"
	$(MOCKGEN) -source=pkg/generated/grpc/$(SERVICE).pb.go -package=pb $(SERVICE_UPPERCASE)Client > pkg/generated/grpc/$(MOCK_PREFIX)_$(SERVICE).pb.go
	@echo ""

.PHONY: mock-proto
mock-proto: $(PROTO_SERVICES)


$(APPS): inject-meta
	@$(eval TARGET := $@)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Building app: $(TARGET)"
	@$(BAZEL) build //cmd/$(TARGET)
	@echo ""

.PHONY: build
build: mock-internal protoc mock-proto swagger-server gazelle $(APPS)

.PHONY: build-just
build-just: gazelle $(APPS)

.PHONY: run
run:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Running app: $(APP)"
	@echo ""
	@bazel-bin/cmd/$(strip $(APP))/*/$(strip $(APP))

.PHONY: cont
cont:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="internal/*/mock_*.go" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build" -command="make run" -graceful-kill=true

.PHONY: cont-just
cont-just:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="internal/*/mock_*.go" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build-just"  -command="make run" -graceful-kill=true

.PHONY: test-unit
test-unit:
	@ echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Testing (Unit)"
	@ $(BAZEL) test //internal/... --features=race
	@ echo ""

.PHONY: test
test: test-unit
	@ echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning test resources (Unit)"
	# @ rm -rf internal/distributed/127.0.0.1:*
	# @ rm -rf internal/distributed/localhost:*
	@ echo ""

