include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/run.mk
include ../scripts/makefiles/test.mk
include ../scripts/makefiles/check.mk

APP			= server					# default app
APPS		:= $(APP)

GITHUB_USER		= 1ambda				# github username
GITHUB_REPO		= go-ref				# github repository
GITHUB_MODULE	= service-gateway		# module name in the repository
DOCKER_REPO		= 1ambda
DOCKER_IMAGE	= goref_service-gateway

SWAGGER_FILE			= ../schema/swagger/gateway.yml
GENERATED_CLIENT_DIR 	= ../service-front/src/app/generated/swagger

swagger-version:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - updating swagger spec version: $(VERSION)"
	@VERSION=$(VERSION) SWAGGER_FILE=$(SWAGGER_FILE) ../scripts/update-swagger-version.sh

swagger-client:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger client files"
	rm -rf $(GENERATED_CLIENT_DIR)	|| true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger client"
	swagger-codegen generate -i $(SWAGGER_FILE) -l typescript-angular -o $(GENERATED_CLIENT_DIR)

swagger: swagger-version
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger files"
	rm -rf pkg/api/model/*	|| true
	rm -rf pkg/api/rest/* 	|| true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - validating swagger spec"
	swagger validate $(SWAGGER_FILE)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating swagger code"
	swagger generate server --spec=$(SWAGGER_FILE) --target=pkg/api --exclude-main --model-package=model --server-package=rest

$(APPS): swagger
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

build: clean $(APPS)

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: clean install build
