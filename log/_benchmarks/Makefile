# TODO: the dependency is from global go path ...
bench:
	go test -bench=. -benchmem
profile:
	go test -run none -bench . -benchtime 3s -benchmem -memprofile p.out
pprof-ui:
# TODO: need to give it binary path
	go tool pprof -http=:8080 p.out