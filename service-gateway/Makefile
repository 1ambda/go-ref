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

SWAGGER_FILE			= ../schema/swagger/gateway.yml

PROTO_DIR  = ../schema/protobuf
PROTO_FILE = $(PROTO_DIR)/gateway.proto
GENERATED_GRPC_SERVER_DIR = "pkg/grpc"

swagger-version:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - updating swagger spec version: $(VERSION)"
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_FILE) ../scripts/update-swagger-version.sh

swagger-server: swagger-version
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger files"
	rm -rf pkg/api/model/*	|| true
	rm -rf pkg/api/rest/* 	|| true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - validating swagger spec"
	swagger validate $(SWAGGER_FILE)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger code"
	swagger generate server --spec=$(SWAGGER_FILE) --target=pkg/api --exclude-main --model-package=model --server-package=rest

swagger-commit:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - commiting swagger related files"
	git add ../schema
	git add pkg/api
	git commit -m "schema: Update swagger spec"

grpc-web-server:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning grpc web client files"
	rm -rf $(GENERATED_GRPC_SERVER_DIR)	|| true
	mkdir -p $(GENERATED_GRPC_SERVER_DIR)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating grpc web client"
	protoc \
    -I $(PROTO_DIR) \
    --go_out=plugins=grpc:$(GENERATED_GRPC_SERVER_DIR) \
    $(PROTO_FILE)

$(APPS): swagger-server grpc-web-server
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

build: clean $(APPS)

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: clean install build
