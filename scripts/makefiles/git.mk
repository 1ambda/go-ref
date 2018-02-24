vendor-commit:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - commiting vendor related files"
	git reset HEAD
	git add ./vendor
	git add Gopkg.toml Gopkg.lock
	git status --short
	git commit -m "vendor: Update vendor packages"
