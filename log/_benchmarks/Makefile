# TODO: the dependency is from global go path ...
#bench:
#	go test -bench=. -benchmem

bench:
	go test -run none -bench . -benchtime 3s -benchmem -memprofile p.out
pprof-ui:
# TODO: need to give it binary path otherwise it will throw error
	go tool pprof -http=:8080 p.out