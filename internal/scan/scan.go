package scan

import (
	"net"
	"os"
	"time"

	"github.com/TimofeiBoldenkov/pscan/internal/tools"
)

type PortState uint8

const (
	Open PortState = iota
	Closed
	Filtered
	Unknown
)

func (ps PortState) String() string {
	switch ps {
	case Open:
		return "open"
	case Closed:
		return "closed"
	case Filtered:
		return "filtered"
	default:
		return "unknown"
	}
}

func ScanHostPort(host string, port string, protocol string, timeout time.Duration) (PortState, error) {
	socket := net.JoinHostPort(host, port)

	_, err := net.DialTimeout(protocol, socket, timeout)

	if os.IsTimeout(err) {
		return Filtered, nil
	} else if tools.IsRefused(err) {
		return Closed, nil
	} else if err == nil {
		return Open, nil
	} else {
		return Unknown, err
	}
}
