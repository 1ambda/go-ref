lint:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - linting"
	@$(GOLINT) --vendor --errors ./... \
		--skip=internal/mock --skip=pkg/generated \
		--enable=unparam --enable=nakedret --enable=staticcheck

style:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - checking code style"
	@! $(GOFMT) -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - formatting code"
	@$(GOCMD) fmt $(GO_FILES)

check: format style lint
