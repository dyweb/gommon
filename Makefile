# based on https://gist.github.com/azatoth/1030091
define GOMMON_MAKEFILE_HELP_MSG
Make commands for gommon

help           show help

Dev:
install           install binaries under ./cmd to $$GOPATH/bin
fmt               gofmt
test              unit test
generate          generate code using gommon
loc               lines of code (cloc required, brew install cloc)

Dev first time:
dep-install    install dependencies based on lock file
dep-update     update dependency based on spec and code

Build:
install        install all binaries under ./cmd to $$GOPATH/bin
build          compile all binary to ./build
build-linux    compile all linux binary to ./build with -linux suffix
build-mac      compile all mac binary to ./build with -mac suffix
build-release  compile all linux and mac binary and generate tarball to ./build

Docker:
docker-build   build runner image w/ all binaries using mulitstage build
docker-push    push runner image to docker registry

endef
export GOMMON_MAKEFILE_HELP_MSG

# -- build vars ---
PKGS=./errors/... ./generator/... ./log/... ./noodle/... ./requests/... ./structure/... ./util/...
PKGST=./cmd ./errors ./generator ./log ./noodle ./requests ./structure ./util
VERSION = 0.0.7
BUILD_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%S%z)
CURRENT_USER = $(USER)
FLAGS = -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.buildTime=$(BUILD_TIME) -X main.buildUser=$(CURRENT_USER)
# -- build vars ---

# TODO: might have a help verbose to and put build and docker commands in it
.PHONY: help
help:
	@echo "$$GOMMON_MAKEFILE_HELP_MSG"

# TODO: define help messages

.PHONY: install
install:
	go install -ldflags "$(FLAGS)" ./cmd/gommon

.PHONY: fmt
fmt:
	gofmt -d -l -w $(PKGST)

.PHONY: test
test:
	go test -v -cover $(PKGS)

.PHONY: test-cover
test-cover:
# https://github.com/codecov/example-go
	go test -coverprofile=coverage.txt -covermode=atomic $(PKGS)

.PHONY: test-race
test-race:
	go test -race $(PKGS)

.PHONY: test-log
test-log:
	go test -v -cover ./log/...



# TODO: refer tools used in https://github.com/360EntSecGroup-Skylar/goreporter
.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: doc
doc:
	xdg-open http://localhost:6060/pkg/github.com/dyweb/gommon &
	godoc -http=":6060"

.PHONY: loc
loc:
	cloc --exclude-dir=vendor,.idea,playground,legacy .

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