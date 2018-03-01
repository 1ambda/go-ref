run: build
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - running app: $(APP)"
	@$(BIN_DIR)/$(APP)

run-cont:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="${BIN_DIR}" \
		-exclude-dir="internal/mock" \
		-exclude-dir="pkg/api" \
		-exclude-dir="pkg/grpc" \
		-build="make build" -command="make run" -graceful-kill=true
