docker-prepare:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - setting GOOS to linux for docker image build"
	@$(eval GOOS := linux)

docker-image: docker-prepare docker-clean build
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - building docker image: $(APP)"
	docker build --build-arg APP=${APP} -t $(DOCKER_IMAGE):local .

docker-tag:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - tagging docker image: latest $(VERSION) $(GIT_COMMIT)"
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT)
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION)
	docker tag $(DOCKER_IMAGE):local $(DOCKER_REPO)/$(DOCKER_IMAGE):latest

docker-push:
	@echo "[$(TAG)] ($(shell TZ=UTC date -u '+%H:%M:%S')) - pushing docker image: latest $(VERSION) $(GIT_COMMIT)"
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT)
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):${VERSION}
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):latest

docker-clean:
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):$(GIT_COMMIT) || true
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION) || true
	docker image rm $(DOCKER_REPO)/$(DOCKER_IMAGE):latest || true
	docker image rm $(DOCKER_IMAGE):local || true

docker: clean docker-image docker-tag docker-push

