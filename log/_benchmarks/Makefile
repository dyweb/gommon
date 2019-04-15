GO = GO111MODULE=on go

bench:
	$(GO) test -run none -bench . -benchtime 3s -benchmem -memprofile p.out
bench-gommon:
	$(GO) test -run none -bench=".*/gommon" -benchtime 3s -benchmem -memprofile p.out
bench-gommon-no-fields:
	$(GO) test -run none -bench="BenchmarkWithoutFieldsText/gommon.F" -benchtime 3s -benchmem -memprofile p.out
bench-gommon-no-context-with-fields:
	$(GO) test -run none -bench="BenchmarkNoContextWithFieldsJSON/gommon.F" -benchtime 3s -benchmem -memprofile p.out
pprof-ui:
# TODO: need to give it binary path otherwise it will throw error
	$(GO) tool pprof -http=:8080 p.out