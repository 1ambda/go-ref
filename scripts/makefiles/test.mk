test:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - testing"
	$(GOTEST) $(GOTEST_OPT) -skip=$(VENDOR_DIR) --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace

test-cont:
	@echo "\n[MAKEFILE] ($(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')) test-cont"
	$(GOTEST) watch $(GOTEST_OPT) -skip=$(VENDOR_DIR)

