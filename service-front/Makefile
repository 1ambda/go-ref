include ../scripts/makefiles/variables.mk
include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/run.mk
include ../scripts/makefiles/test.mk
include ../scripts/makefiles/check.mk
include ../scripts/makefiles/git.mk

GITHUB_MODULE	= service-front
DOCKER_IMAGE	= $(shell echo $(DOCKER_IMAGE_PREFIX)_$(GITHUB_MODULE))

SWAGGER_FILE			            = "../schema/swagger/gateway.yml"
GENERATED_SWAGGER_CLIENT_DIR  = "./src/app/generated/swagger"

PROTO_DIR  = "../schema/protobuf"
PROTO_FILE = "$(PROTO_DIR)/gateway.proto"
GENERATED_GRPC_CLIENT_DIR = "./src/app/generated/grpc"

swagger-client:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger client files"
	rm -rf $(GENERATED_SWAGGER_CLIENT_DIR)	|| true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger client"
	swagger-codegen generate -i $(SWAGGER_FILE) -l typescript-angular -o $(GENERATED_SWAGGER_CLIENT_DIR)

grpc-web-client:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning grpc web client files"
	rm -rf $(GENERATED_GRPC_CLIENT_DIR)	|| true
	mkdir -p $(GENERATED_GRPC_CLIENT_DIR)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating grpc web client"
	protoc \
    --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
    -I $(PROTO_DIR) \
    --js_out=import_style=commonjs,binary:$(GENERATED_GRPC_CLIENT_DIR) \
    --ts_out=service=true:$(GENERATED_GRPC_CLIENT_DIR) \
    $(PROTO_FILE)

