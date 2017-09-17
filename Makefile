.PHONY: test
test:
	go test -v -cover $(shell glide novendor)
.PHONY: fmt
fmt:
	gofmt -d -l -w ./cast ./config ./log ./requests ./runner ./structure ./util