PROJECT_ROOT		:= github.com/Infoblox-CTO/heka-ui
BUILD_PATH		:= bin
DOCKERFILE_PATH		:= $(CURDIR)/docker

# configuration for image names
USERNAME		:= $(USER)
GIT_COMMIT		:= $(shell git describe --tags --dirty=-unsupported --always || echo pre-commit)
IMAGE_VERSION		?= $(GIT_COMMIT)
IMAGE_REGISTRY		?= infobloxcto

# configuration for server binary and image
SERVER_IMAGE		:= $(IMAGE_REGISTRY)/heka-ui
SERVER_DOCKERFILE	:= $(DOCKERFILE_PATH)/Dockerfile

# configuration for building on host machine
GO_CACHE		:= -pkgdir $(BUILD_PATH)/go-cache
GO_BUILD_FLAGS		?= $(GO_CACHE) -i -v
GO_TEST_FLAGS		?= -v -cover
GO_PACKAGES		:= $(shell go list ./... | grep -v vendor)
GO_MOD			= go.mod

export GOPRIVATE	?= github.com/Infoblox-CTO
export GOFLAGS		?= -mod=vendor


.PHONY: all
all: vendor docker

.PHONY: fmt
fmt:
	@echo "Not Implemented for now, getting error while pulling private repos"
	#@go fmt $(GO_PACKAGES)

.PHONY: lint
lint:
	@! gofmt -l . | grep -v vendor/

.PHONY: test
test: fmt
	@echo "Not Implemented for now, getting error while pulling private repos"
	#@go test $(GO_TEST_FLAGS) $(GO_PACKAGES)

.PHONY: docker
docker:
	@docker build -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE):$(IMAGE_VERSION) .
	@docker tag $(SERVER_IMAGE):$(IMAGE_VERSION) $(SERVER_IMAGE):latest
	@docker image prune -f --filter label=stage=frontend
	@docker image prune -f --filter label=stage=backend

.PHONY: push
push:
	@docker push $(SERVER_IMAGE):$(IMAGE_VERSION)

.PHONY: push-latest
push-latest:
	@docker push $(SERVER_IMAGE):latest

.PHONY: vendor
vendor:
	@export GO111MODULE=on; go mod tidy; go mod vendor; unset GO111MODULE

.PHONY: clean
clean:
	@docker rmi -f $(shell docker images -q $(SERVER_IMAGE)) || true
