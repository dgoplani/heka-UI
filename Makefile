PROJECT_ROOT            ?= $(PWD)
DOCKER_PATH		?= $(CURDIR)/docker
GIT_COMMIT              ?= $(shell git describe --dirty=-unsupported --always --tags || echo pre-commit)
IMAGE_VERSION           ?= $(GIT_COMMIT)
IMAGE_REGISTRY          ?= infobloxcto
SERVER_IMAGE            ?= $(IMAGE_REGISTRY)/heka-ui
SERVER_DOCKERFILE       ?= $(DOCKER_PATH)/Dockerfile
DATE_TIME		?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: docker

docker:
	@docker build --build-arg GIT_TAG=$(GIT_COMMIT) --build-arg BUILD_DATE=$(DATE_TIME) -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE):$(IMAGE_VERSION) .
	@docker tag $(SERVER_IMAGE):$(IMAGE_VERSION) $(SERVER_IMAGE):latest
	@docker image prune -f --filter label=stage=frontend
	@docker image prune -f --filter label=stage=backend
