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

SERVICES	:= hello location
SCHEMA_DIR	:= ../schema/protobuf

MOCKGEN = mockgen
MOCK_PREFIX = mock_
MOCK_PKG_DIRS := "internal/hello" "internal/distributed"

$(SERVICES): protoc
	@$(eval SERVICE := $@)
	# capitalize first letter
	@$(eval SERVICE_UPPERCASE := $(shell echo `echo $${SERVICE:0:1} | tr  '[a-z]' '[A-Z]'`$${SERVICE:1}))
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating protobuf mock files: $(SERVICE)"
	mockgen -source=pkg/generated/grpc/$(SERVICE).pb.go -package=pb $(SERVICE_UPPERCASE)Client > pkg/generated/grpc/mock_$(SERVICE).pb.go

protoc:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning swagger files"
	rm -rf pkg/generated/grpc/* || true
	mkdir -p pkg/generated/grpc || true

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - generating protobuf service files: $(SERVICE)"
	protoc  --proto_path $(SCHEMA_DIR) --go_out=plugins=grpc:pkg/generated/grpc $(SCHEMA_DIR)/*.proto

$(APPS):
	@$(eval TARGET := $@)
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building app: $(TARGET)"
	@$(eval LDFLAG := $(shell govvv -flags -pkg $(shell go list ./internal/pkg/config)))
	echo $(LDFLAG)
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

build: clean protoc $(SERVICES) $(APPS)

build-just: $(APPS)

.PHONY: clean
clean:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - cleaning"
	rm -rf ./bin

.PHONY: dep
dep:
	dep ensure -v
	find vendor/ -name "BUILD" -delete
	find vendor/ -name "BUILD.bazel" -delete
	bazel run //:gazelle -- -proto disable
	rm -f vendor/github.com/coreos/etcd/cmd/etcd
