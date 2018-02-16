.PHONY: install
install:
	go install ./cmd/gommon

.PHONY: test
test:
	go test -v -cover ./cast/... ./config/... ./generator/... ./log/... ./noodle/... ./requests/... ./runner/... ./structure/... ./util/...

.PHONY: test-log
test-log:
	go test -v -cover ./log/...

.PHONY: fmt
fmt:
	gofmt -d -l -w ./cast ./config ./generator ./log ./noodle ./requests ./runner ./structure ./util

.PHONY: docker-test
docker-test:
	docker-compose -f scripts/docker-compose.yml run --rm golang1.9
	docker-compose -f scripts/docker-compose.yml run --rm golanglatest

.PHONY: docker-remove-all-containers
docker-remove-all-containers:
	docker rm $(shell docker ps -a -q)