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

SWAGGER_CLIENT_DIR      = "./src/app/generated/swagger"
SWAGGER_REST_FILE       = "../schema/swagger/gateway-rest.yml"
SWAGGER_REST_CLIENT_DIR = "./src/app/generated/swagger/rest"
SWAGGER_WS_FILE         = "../schema/swagger/gateway-websocket.yml"
SWAGGER_WS_CLIENT_DIR   = "./src/app/generated/swagger/websocket"

swagger-client:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger client files"
	rm -rf $(SWAGGER_CLIENT_DIR)	|| true
	mkdir -p $(SWAGGER_REST_CLIENT_DIR) || true
	mkdir -p $(SWAGGER_WS_CLIENT_DIR) || true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger (REST) client"
	swagger-codegen generate -i $(SWAGGER_REST_FILE) -l typescript-angular -o $(SWAGGER_REST_CLIENT_DIR)

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger (WS) client"
	swagger-codegen generate -i $(SWAGGER_WS_FILE) -l typescript-angular -o $(SWAGGER_WS_CLIENT_DIR)
