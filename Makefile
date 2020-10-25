# based on https://gist.github.com/azatoth/1030091
# TODO(at15): it is also possible to generate it automatically using awk etc.
define GOMMON_MAKEFILE_HELP_MSG
Make commands for gommon

help           show help

Dev
-----------------------------------------
install           install binaries under ./cmd to $$GOPATH/bin
fmt               goimports
test              run unit test
generate          generate code using gommon
loc               lines of code (cloc required, brew install cloc)

Build
-----------------------------------------
install        install all binaries under ./cmd to $$GOPATH/bin
build          compile all binary to ./build for current platform
build-linux    compile all linux binary to ./build with -linux suffix
build-mac      compile all mac binary to ./build with -mac suffix
build-win      compile all windows binary to ./build with -win suffix
build-release  compile binary for all platforms and generate tarball to ./build

Docker
-----------------------------------------
docker-build   build runner image w/ all binaries using mulitstage build
docker-push    push runner image to docker registry

endef
export GOMMON_MAKEFILE_HELP_MSG

# TODO: might have a help verbose to and put build and docker commands in it
.PHONY: help
help:
	@echo "$$GOMMON_MAKEFILE_HELP_MSG"

GO = GO111MODULE=on CGO_ENABLED=0 go
# -- build vars ---
PKGST =./cmd ./dcli ./errors ./generator ./httpclient ./linter ./log ./noodle ./util ./tconfig
PKGS = $(addsuffix ...,$(PKGST))
VERSION = 0.0.13
BUILD_COMMIT := $(shell git rev-parse HEAD)
BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S%z)
CURRENT_USER = $(USER)
FLAGS = -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.buildTime=$(BUILD_TIME) -X main.buildUser=$(CURRENT_USER)
DOCKER_REPO = dyweb/gommon
DCLI_PKG = github.com/dyweb/gommon/dcli.
DCLI_LDFLAGS = -X $(DCLI_PKG)buildVersion=$(VERSION) -X $(DCLI_PKG)buildCommit=$(BUILD_COMMIT) -X $(DCLI_PKG)buildBranch=$(BUILD_BRANCH) -X $(DCLI_PKG)buildTime=$(BUILD_TIME) -X $(DCLI_PKG)buildUser=$(CURRENT_USER)
# -- build vars ---

.PHONY: install
install: fmt test install-only

install-only:
	cd ./cmd/gommon && $(GO) install -ldflags "$(FLAGS)" .
	mv $(GOPATH)/bin/gommonbin $(GOPATH)/bin/gommon

.PHONY: install2
install2:
	cd ./cmd/gommon2 && $(GO) install -ldflags "$(DCLI_LDFLAGS)" .

.PHONY: fmt
fmt:
	gommon format -d -l -w $(PKGST)

# gommon format is a drop in replacement for goimports
deprecated-fmt:
	goimports -d -l -w $(PKGST)

# --- build ---
.PHONY: clean build build-linux build-mac build-win build-all
clean:
	rm -f ./build/gommon
	rm -f ./build/gommon-*
build:
	$(GO) build -ldflags "$(FLAGS)" -o ./build/gommon ./cmd/gommon
build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(FLAGS)" -o ./build/gommon-linux ./cmd/gommon
build-mac:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(FLAGS)" -o ./build/gommon-mac ./cmd/gommon
build-win:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(FLAGS)" -o ./build/gommon-win ./cmd/gommon
build-all: build build-linux build-mac build-win
build-release: clean build-all
	zip ./build/gommon-linux.zip ./build/gommon-linux
	zip ./build/gommon-mac.zip ./build/gommon-mac
	zip ./build/gommon-win.zip ./build/gommon-win
# --- build ---

.PHONY: generate
generate:
	gommon generate -v

# --- test ---
.PHONY: test test-cover test-cover-html test-race
.PHONY: test-log test-errors

test:
	$(GO) test -cover $(PKGS)

test-verbose:
	$(GO) test -v -cover $(PKGS)

test-cover:
# https://github.com/codecov/example-go
	$(GO) test -coverprofile=coverage.txt -covermode=atomic $(PKGS)

test-cover-html: test-cover
	$(GO) tool cover -html=coverage.txt

test-race:
	$(GO) test -race $(PKGS)

test-log:
	$(GO) test -v -cover ./log/...
test-errors:
	$(GO) test -v -cover ./errors/...
# --- test ---

# TODO: refer tools used in https://github.com/360EntSecGroup-Skylar/goreporter
.PHONY: vet
vet:
	$(GO) vet $(PKGS)

.PHONY: doc
doc:
	xdg-open http://localhost:6060/pkg/github.com/dyweb/gommon &
	godoc -http=":6060"

# TODO: ignore example, test, legacy etc.
# https://github.com/Aaronepower/tokei
.PHONY: loc
loc:
	tokei .

# --- dependency management ---
mod-init:
	$(GO) mod init
mod-update:
	$(GO) mod tidy
mod-graph:
	$(GO) mod graph
# --- dependency management ---

# --- docker ---
.PHONY: docker-build docker-push docker-test

docker-build:
	docker build -t $(DOCKER_REPO):$(VERSION) .

docker-push:
	docker push $(DOCKER_REPO):$(VERSION)

#.PHONY: docker-remove-all-containers
#docker-remove-all-containers:
#	docker rm $(shell docker ps -a -q)
# --- docker ---
