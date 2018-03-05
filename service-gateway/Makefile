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

SWAGGER_REST_SPEC	= ../schema/swagger/gateway-rest.yml
SWAGGER_WS_SPEC		= ../schema/swagger/gateway-websocket.yml

PROTO_DIR  = ../schema/protobuf
PROTO_FILE = $(PROTO_DIR)/gateway.proto
GENERATED_GRPC_SERVER_DIR = "pkg/grpc"

swagger-version:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - updating swagger spec version: $(VERSION)"
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_REST_SPEC) ../scripts/update-swagger-version.sh
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_WS_SPEC) ../scripts/update-swagger-version.sh

swagger-server: swagger-version
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger files"
	rm -rf pkg/generated/swagger/* || true
	mkdir -p pkg/generated/swagger || true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - validating swagger (REST) spec"
	swagger validate $(SWAGGER_REST_SPEC)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger (REST) code"
	swagger generate server --spec=$(SWAGGER_REST_SPEC) --exclude-main \
		--target=pkg/generated/swagger --model-package=rest_model \
		--server-package=rest_server --api-package=rest_api


	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - validating swagger (WS) spec"
	swagger validate $(SWAGGER_WS_SPEC)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger (WS) code"
	swagger generate server --spec=$(SWAGGER_WS_SPEC) --exclude-main \
		--target=pkg/generated/swagger --model-package=ws_model \
		--skip-operations --skip-support --exclude-spec

swagger-commit:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - commiting swagger related files"
	git add ../schema
	git add pkg/api
	git commit -m "schema: Update swagger spec"

$(APPS):
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

build: clean swagger-server $(APPS)

build-just: $(APPS)

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: clean install build
