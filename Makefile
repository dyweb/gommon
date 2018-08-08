PKGS=./config/... ./errors/... ./generator/... ./log/... ./noodle/... ./requests/... ./structure/... ./util/...
PKGST=./cmd ./config ./errors ./generator ./log ./noodle ./requests ./structure ./util
VERSION = 0.0.1
BUILD_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%S%z)
CURRENT_USER = $(USER)
FLAGS = -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.buildTime=$(BUILD_TIME) -X main.buildUser=$(CURRENT_USER)

# TODO: define help messages

.PHONY: install
install:
	go install -ldflags "$(FLAGS)" ./cmd/gommon

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

.PHONY: fmt
fmt:
	gofmt -d -l -w $(PKGST)

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

.PHONY: update-dep
update-dep:
	dep ensure -update

#--- docker ---#
.PHONY: docker-test
docker-test:
	docker-compose -f scripts/docker-compose.yml run --rm golang1.10
	docker-compose -f scripts/docker-compose.yml run --rm golanglatest

.PHONY: docker-remove-all-containers
docker-remove-all-containers:
	docker rm $(shell docker ps -a -q)
#--- docker ---#