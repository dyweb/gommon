# TODO: the dependency is from global go path ...
bench:
	go test -run none -bench . -benchtime 3s -benchmem -memprofile p.out
bench-gommon:
	go test -run none -bench=".*/gommon" -benchtime 3s -benchmem -memprofile p.out
bench-gommon-no-fields:
	go test -run none -bench="BenchmarkWithoutFieldsText/gommon.F" -benchtime 3s -benchmem -memprofile p.out
pprof-ui:
# TODO: need to give it binary path otherwise it will throw error
	go tool pprof -http=:8080 p.out