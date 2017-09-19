.PHONY: test
test:
	go test -v -cover $(shell glide novendor)

.PHONY: test-log
test-log:
	go test -v -cover ./log/...

.PHONY: fmt
fmt:
	gofmt -d -l -w ./cast ./config ./log ./requests ./runner ./structure ./util