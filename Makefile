# based on https://gist.github.com/azatoth/1030091
define GOMMON_MAKEFILE_HELP_MSG
Make commands for gommon

help           show help

Dev:
install           install binaries under ./cmd to $$GOPATH/bin
fmt               goimports
test              unit test
generate          generate code using gommon
loc               lines of code (cloc required, brew install cloc)

Dev first time:
dep-install    install dependencies based on lock file
dep-update     update dependency based on spec and code

Build:
install        install all binaries under ./cmd to $$GOPATH/bin
build          compile all binary to ./build for current platform
build-linux    compile all linux binary to ./build with -linux suffix
build-mac      compile all mac binary to ./build with -mac suffix
build-win      compile all windows binary to ./build with -win suffix
build-release  compile binary for all platforms and generate tarball to ./build

Docker:
docker-build   build runner image w/ all binaries using mulitstage build
docker-push    push runner image to docker registry

endef
export GOMMON_MAKEFILE_HELP_MSG

# TODO: might have a help verbose to and put build and docker commands in it
.PHONY: help
help:
	@echo "$$GOMMON_MAKEFILE_HELP_MSG"

# -- build vars ---
PKGS =./errors/... ./generator/... ./log/... ./noodle/... ./structure/... ./util/...
PKGST =./cmd ./errors ./generator ./log ./noodle ./structure ./util
VERSION = 0.0.8
BUILD_COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S%z)
CURRENT_USER = $(USER)
FLAGS = -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.buildTime=$(BUILD_TIME) -X main.buildUser=$(CURRENT_USER)
# -- build vars ---

.PHONY: install
install:
	go install -ldflags "$(FLAGS)" ./cmd/gommon

.PHONY: fmt
fmt:
	goimports -d -l -w $(PKGST)

# --- build ---
.PHONY: clean build build-linux build-mac build-win build-all
clean:
	rm ./build/gommon
	rm ./build/gommon-*
build:
	go build -ldflags "$(FLAGS)" -o ./build/gommon ./cmd/gommon
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(FLAGS)" -o ./build/gommon-linux ./cmd/gommon
build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(FLAGS)" -o ./build/gommon-mac ./cmd/gommon
build-win:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(FLAGS)" -o ./build/gommon-win ./cmd/gommon
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
	go test -v -cover $(PKGS)

test-cover:
# https://github.com/codecov/example-go
	go test -coverprofile=coverage.txt -covermode=atomic $(PKGS)

test-cover-html: test-cover
	go tool cover -html=coverage.txt

test-race:
	go test -race $(PKGS)

test-log:
	go test -v -cover ./log/...
test-errors:
	go test -v -cover ./errors/...
# --- test ---

# TODO: refer tools used in https://github.com/360EntSecGroup-Skylar/goreporter
.PHONY: vet
vet:
	go vet $(PKGS)

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
.PHONY: dep-install dep-update
dep-install:
	dep ensure -v

dep-update:
	dep ensure -v -update
# --- dependency management ---

# --- docker ---
.PHONY: docker-test
docker-test:
	docker-compose -f scripts/docker-compose.yml run --rm golang1.10
	docker-compose -f scripts/docker-compose.yml run --rm golanglatest

.PHONY: docker-remove-all-containers
docker-remove-all-containers:
	docker rm $(shell docker ps -a -q)
# --- docker ---
