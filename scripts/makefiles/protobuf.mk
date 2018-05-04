.PHONY: protoc
protoc:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Cleaning grpc files"
	@rm -rf pkg/generated/grpc/* || true
	@mkdir -p pkg/generated/grpc || true
	@echo ""

	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating protobuf service files"
	protoc  --proto_path $(SCHEMA_DIR) --go_out=plugins=grpc:pkg/generated/grpc $(SCHEMA_DIR)/*.proto
	@echo ""
