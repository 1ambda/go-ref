include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/run.mk
include ../scripts/makefiles/test.mk
include ../scripts/makefiles/check.mk

APP			= server					# default app
APPS		:= $(APP) client

GITHUB_USER		= 1ambda				# github username
GITHUB_REPO		= go-ref				# github repository
GITHUB_MODULE	= service-backend		# module name in the repository
DOCKER_REPO		= 1ambda
DOCKER_IMAGE	= goref_service-backend

# `api/$(SERVICE).proto					: protobuf definition
# `pkg/api/$(SERVICE).pb.go`			: (auto-generated) protobuf service
# `pkg/api/mock_$(SERVICE).pb.go`	: (auto-generated) protobuf mock
SERVICES	:= hello
SCHEMA_DIR	:= ../schema/protobuf

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

build: clean $(APPS)

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: clean install build
