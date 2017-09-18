.PHONY: all
all: test zap

.PHONY: zap
zap:
	go run echart.go main.go fixture/zap.txt zap.html

.PHONY: test
test:
	go test .

.PHONY: install
install:
	go install .

.PHONY: fmt
fmt:
	gofmt -d -l -w .