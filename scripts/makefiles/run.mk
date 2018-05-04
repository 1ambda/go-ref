run:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Running app: $(APP)"
	@echo ""
	@bazel-bin/cmd/$(strip $(APP))/*/$(strip $(APP))

cont:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="/private/var/tmp/*" \
		-exclude-dir="internal/*/mock_*.go" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build" -command="make run" -graceful-kill=true

cont-just:
	CompileDaemon \
		-exclude-dir="${VENDOR_DIR}" \
		-exclude-dir="/private/var/tmp/*" \
		-exclude-dir="internal/*/mock_*.go" \
		-exclude-dir="pkg/generated" \
		-log-prefix=false \
		-build="make build-just"  -command="make run" -graceful-kill=true
