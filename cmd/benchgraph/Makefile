.PHONY: all
all: test zap

.PHONY: zap
zap:
	go run echart.go echarts_tmpl.go main.go fixture/zap.txt

.PHONY: test
test:
	go test .