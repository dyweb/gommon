.PHONY: test
test:
	go test -v -cover $(shell glide novendor)

.PHONY: test-log
test-log:
	go test -v -cover ./log/...

.PHONY: fmt
fmt:
	gofmt -d -l -w ./cast ./config ./log ./requests ./runner ./structure ./util

.PHONY: docker-test
docker-test:
	docker-compose -f scripts/docker-compose.yml run --rm golang1.7
	docker-compose -f scripts/docker-compose.yml run --rm golang1.8
	docker-compose -f scripts/docker-compose.yml run --rm golang1.9
	docker-compose -f scripts/docker-compose.yml run --rm golanglatest

.PHONY: docker-remove-all-containers
docker-remove-all-containers:
	docker rm $(shell docker ps -a -q)