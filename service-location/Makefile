include ../scripts/makefiles/variables.mk
include ../scripts/makefiles/basic.mk
include ../scripts/makefiles/docker.mk
include ../scripts/makefiles/install.mk
include ../scripts/makefiles/run.mk
include ../scripts/makefiles/test.mk
include ../scripts/makefiles/check.mk
include ../scripts/makefiles/git.mk

APP			= server					# default app
APPS		:= $(APP) client

GITHUB_MODULE	= service-location		# module name in the repository
DOCKER_IMAGE	= $(shell echo $(DOCKER_IMAGE_PREFIX)_$(GITHUB_MODULE))

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

$(APPS):
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	$(GOBUILD) -i -ldflags="$(LDFLAG)" -o $(BIN_DIR)/$(TARGET) -v $(CMD_DIR)/$(TARGET)/$(MAIN)

$(MOCK_PKG_DIRS):
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning mock files under $(TARGET)"
	rm -rf $(TARGET)/$(MOCK_PREFIX)*
	@echo ""

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating mock files for $(TARGET)"
	PKG_DIR=$(TARGET) ../scripts/generate-gomock.sh
	@echo ""

mock-generate: $(MOCK_PKG_DIRS)

build: clean protoc $(APPS)

build-just: $(APPS)

clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: clean install build
