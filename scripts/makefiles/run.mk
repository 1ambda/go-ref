run-only:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - running app: $(APP)"
	@$(BIN_DIR)/$(APP)

run: build run-only

cont:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="${BIN_DIR}" \
		-exclude-dir="internal/mock" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build" -command="make run-only" -graceful-kill=true

run-just: build-just run-only

cont-just:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="${BIN_DIR}" \
		-exclude-dir="internal/mock" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build-just"  -command="make run-only" -graceful-kill=true
