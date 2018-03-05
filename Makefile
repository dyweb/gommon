PKGS=./cast/... ./config/... ./errors/... ./generator/... ./log/... ./noodle/... ./requests/... ./structure/... ./util/...
PKGST=./cast ./config ./errors ./generator ./log ./noodle ./requests ./structure ./util

.PHONY: install
install:
	go install ./cmd/gommon

.PHONY: test
test:
	go test -v -cover $(PKGS)

.PHONY: test-race
test-race:
	go test -race $(PKGS)

.PHONY: test-log
test-log:
	go test -v -cover ./log/...

.PHONY: fmt
fmt:
	gofmt -d -l -w $(PKGST)

.PHONY: vet
vet:
	go vet $(PKGST)

.PHONy: doc
doc:
	xdg-open http://localhost:6060/pkg/github.com/dyweb/gommon &
	godoc -http=":6060"

.PHONY: loc
loc:
	cloc --exclude-dir=vendor,.idea,playground,legacy .

#--- docker ---#
.PHONY: docker-test
docker-test:
	docker-compose -f scripts/docker-compose.yml run --rm golang1.9
	docker-compose -f scripts/docker-compose.yml run --rm golanglatest

.PHONY: docker-remove-all-containers
docker-remove-all-containers:
	docker rm $(shell docker ps -a -q)
#--- docker ---#