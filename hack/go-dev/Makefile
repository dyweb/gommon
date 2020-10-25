DOCKER_REPO = dyweb/go-dev
GO_VERSIONS = 1.15.3
BUILDS = $(addprefix build-, $(GO_VERSIONS))
PUSHS = $(addprefix push-, $(GO_VERSIONS))

.PHONY: build push

$(BUILDS):
	docker build -t $(DOCKER_REPO):$(subst build-,,$@) --build-arg BUILD_GO_VERSION=$(subst build-,,$@) .

$(PUSHS):
	docker push $(DOCKER_REPO):$(subst push-,,$@)

build: $(BUILDS)

push: $(PUSHS)

run:
	docker run --rm -it --entrypoint /bin/bash $(DOCKER_REPO):1.15.3