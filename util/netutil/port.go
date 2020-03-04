package netutil

import (
	"fmt"
	"net"
)

// TODO: ephemeral port range is configurable, this is what I got from wiki
const (
	ephemeralStart = 32768
	ephemeralEnd   = 61000
)

// AvailablePortBySystem use 0 as port number so system finds one for us.
// Based on https://github.com/phayes/freeport/blob/master/freeport.go
func AvailablePortBySystem() (int, error) {
	// 0 means have system find the next available port for you.
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil

}

// AvailablePortByRange scan and range of ports to find a free one.
// Based on https://github.com/sosedoff/pgweb/blob/master/pkg/connection/port.go
func AvailablePortByRange(start, end int) (int, error) {
	for p := start; p < end; p++ {
		l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", p))
		if err != nil {
			continue
		}
		l.Close()
		return p, nil
	}
	return 0, fmt.Errorf("not free port found from %d to %d", start, end)
}
