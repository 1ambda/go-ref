.PHONY: gazelle
gazelle:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Generating BUILD.bazel files"
	@ find vendor/ -name "BUILD" -delete
	@ find vendor/ -name "BUILD.bazel" -delete
	@$(BAZEL) run //:gazelle -- -proto disable
	@rm -f vendor/github.com/coreos/etcd/cmd/etcd
	@echo ""

.PHONY: inject-meta
inject-meta:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - Injecting LDFLAGS into BUILD.bazel"
	@WORKSPACE=$(ROOT_IMPORT_PATH)/$(GITHUB_MODULE) CONFIG_PACKAGE_PATH=$(CONFIG_PACKAGE_PATH) ../scripts/get-ldflags.sh
	@echo ""

